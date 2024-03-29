package wsutil

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/ronaksoft/rony/pools"

	"github.com/gobwas/ws"
)

// Message represents a message from peer, that could be presented in one or
// more frames. That is, it contains payload of all message fragments and
// operation code of initial frame for this message.
type Message struct {
	Payload []byte
	OpCode  ws.OpCode
}

// ReadMessage is a helper function that reads next message from r. It appends
// received message(s) to the third argument and returns the result of it and
// an error if some failure happened. That is, it probably could receive more
// than one message when peer sending fragmented message in multiple frames and
// want to send some control frames between fragments. Then returned slice will
// contain those control frames at first, and then result of gluing fragments.
func ReadMessage(r io.Reader, s ws.State, m []Message) ([]Message, error) {
	im := Message{}
	rd := Reader{
		Source: r,
		State:  s,
		OnIntermediate: func(hdr ws.Header, src io.Reader) error {
			bts, err := ioutil.ReadAll(src)
			if err != nil {
				return err
			}
			m = append(m, Message{OpCode: hdr.OpCode, Payload: bts})

			return nil
		},
	}
	h, err := rd.NextFrame()
	if err != nil {
		return m, err
	}

	if h.Fin {
		im.Payload = pools.Bytes.GetLen(int(h.Length))
		// It is not possible to receive io.EOF here because Reader does not
		// return EOF if frame payload was successfully fetched.
		// Thus, we are consistent here with io.Reader behavior.
		_, err = io.ReadFull(&rd, im.Payload)
	} else {
		// Frame is fragmented, thus use ioutil.ReadAll behavior.
		var buf bytes.Buffer
		_, err = buf.ReadFrom(&rd)
		im.Payload = append(im.Payload, buf.Bytes()...)
	}
	if err != nil {
		return m, err
	}

	im.OpCode = h.OpCode

	return append(m, im), nil
}

// WriteMessage is a helper function that writes message to the w. It
// constructs single frame with given operation code and payload.
// It uses given state to prepare side-dependent things, like cipher
// payload bytes from client to server. It will not mutate p bytes if
// cipher must be made.
//
// If you want to write message in fragmented frames, use Writer instead.
func WriteMessage(w io.Writer, s ws.State, op ws.OpCode, p []byte) error {
	return writeFrame(w, s, op, true, p)
}
