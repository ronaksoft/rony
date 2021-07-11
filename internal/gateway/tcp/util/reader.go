package wsutil

import (
	"errors"
	"io"
	"io/ioutil"

	"github.com/gobwas/ws"
)

// ErrNoFrameAdvance means that Reader's Read() method was called without
// preceding NextFrame() call.
var ErrNoFrameAdvance = errors.New("no frame advance")

// FrameHandlerFunc handles parsed frame header and its body represented by
// io.Reader.
//
// Note that reader represents already unmasked body.
type FrameHandlerFunc func(ws.Header, io.Reader) error

// Reader is a wrapper around source io.Reader which represents WebSocket
// connection. It contains options for reading messages from source.
//
// Reader implements io.Reader, which Read() method reads payload of incoming
// WebSocket frames. It also takes care on fragmented frames and possibly
// intermediate control frames between them.
//
// Note that Reader's methods are not goroutine safe.
type Reader struct {
	Source io.Reader
	State  ws.State

	// SkipHeaderCheck disables checking header bits to be RFC6455 compliant.
	SkipHeaderCheck bool
	OnContinuation  FrameHandlerFunc
	OnIntermediate  FrameHandlerFunc

	opCode ws.OpCode        // Used to store message op code on fragmentation.
	frame  io.Reader        // Used to as frame reader.
	raw    io.LimitedReader // Used to discard frames without cipher.
}

// NewReader creates new frame reader that reads from r keeping given state to
// make some protocol validity checks when it needed.
func NewReader(r io.Reader, s ws.State) *Reader {
	return &Reader{
		Source: r,
		State:  s,
	}
}

// Read implements io.Reader. It reads the next message payload into p.
// It takes care on fragmented messages.
//
// The error is io.EOF only if all of message bytes were read.
// If an io.EOF happens during reading some but not all the message bytes
// Read() returns io.ErrUnexpectedEOF.
//
// The error is ErrNoFrameAdvance if no NextFrame() call was made before
// reading next message bytes.
func (r *Reader) Read(p []byte) (n int, err error) {
	if r.frame == nil {
		if !r.fragmented() {
			// Every new Read() must be preceded by NextFrame() call.
			return 0, ErrNoFrameAdvance
		}
		// Read next continuation or intermediate control frame.
		_, err := r.NextFrame()
		if err != nil {
			return 0, err
		}
		if r.frame == nil {
			// We handled intermediate control and now got nothing to read.
			return 0, nil
		}
	}

	n, err = r.frame.Read(p)
	if err != nil && err != io.EOF {
		return
	}
	if err == nil && r.raw.N != 0 {
		return
	}

	switch {
	case r.raw.N != 0:
		err = io.ErrUnexpectedEOF

	case r.fragmented():
		err = nil
		r.resetFragment()

	default:
		r.reset()
		err = io.EOF
	}

	return
}

// Discard discards current message unread bytes.
// It discards all frames of fragmented message.
func (r *Reader) Discard() (err error) {
	for {
		_, err = io.Copy(ioutil.Discard, &r.raw)
		if err != nil {
			break
		}
		if !r.fragmented() {
			break
		}
		if _, err = r.NextFrame(); err != nil {
			break
		}
	}
	r.reset()

	return err
}

// NextFrame prepares r to read next message. It returns received frame header
// and non-nil error on failure.
//
// Note that next NextFrame() call must be done after receiving or discarding
// all current message bytes.
func (r *Reader) NextFrame() (hdr ws.Header, err error) {
	hdr, err = ws.ReadHeader(r.Source)
	if err == io.EOF && r.fragmented() {
		// If we are in fragmented state EOF means that is was totally
		// unexpected.
		//
		// NOTE: This is necessary to prevent callers such that
		// ioutil.ReadAll to receive some amount of bytes without an error.
		// ReadAll() ignores an io.EOF error, thus caller may think that
		// whole message fetched, but actually only part of it.
		err = io.ErrUnexpectedEOF
	}
	if err == nil && !r.SkipHeaderCheck {
		err = ws.CheckHeader(hdr, r.State)
	}
	if err != nil {
		return hdr, err
	}

	// Save raw reader to use it on discarding frame without ciphering and
	// other streaming checks.
	r.raw = io.LimitedReader{R: r.Source, N: hdr.Length}

	frame := io.Reader(&r.raw)
	if hdr.Masked {
		frame = NewCipherReader(frame, hdr.Mask)
	}
	if r.fragmented() {
		if hdr.OpCode.IsControl() {
			if cb := r.OnIntermediate; cb != nil {
				err = cb(hdr, frame)
			}
			if err == nil {
				// Ensure that src is empty.
				_, err = io.Copy(ioutil.Discard, &r.raw)
			}

			return
		}
	} else {
		r.opCode = hdr.OpCode
	}

	// Save reader with ciphering and other streaming checks.
	r.frame = frame

	if hdr.OpCode == ws.OpContinuation {
		if cb := r.OnContinuation; cb != nil {
			err = cb(hdr, frame)
		}
	}

	if hdr.Fin {
		r.State = r.State.Clear(ws.StateFragmented)
	} else {
		r.State = r.State.Set(ws.StateFragmented)
	}

	return
}

func (r *Reader) fragmented() bool {
	return r.State.Fragmented()
}

func (r *Reader) resetFragment() {
	r.raw = io.LimitedReader{}
	r.frame = nil
}

func (r *Reader) reset() {
	r.raw = io.LimitedReader{}
	r.frame = nil
	r.opCode = 0
}
