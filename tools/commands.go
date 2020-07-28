package tools

import (
	"errors"
	"io"
	"os/exec"
)

var (
	ErrNoCmds = errors.New("pipe: there were no commands provided")
)

// Commands takes in a series of *exec.Cmd commands and
// runs them in the order provided, piping the output
// from each command to the input of the next.
//
// This is pretty experimental, and I don't find it to
// be very useful unless you intend to ignore errors
// at which point it is easier than using the `Cmds`
// type.
//
// Commands returns an io.ReadCloser and an
// io.WriteCloser which can be used much like you would
// use [io.Pipe](https://golang.org/pkg/io/#Pipe).
// It also returns a third argument which is an error
// channel that will eventually pipe any errors
// encountered or close if things run successfully.
// If the reader or writer are nil, then there should
// already be an error in the error channel.
//
// NOTE: YOU MUST call Close() on the io.WriteCloser()
// to kick off the chaining of commands.
func Commands(cmds ...*exec.Cmd) (io.ReadCloser, io.WriteCloser, chan error) {
	errCh := make(chan error, 1)
	cmd, err := NewPipeCommands(cmds...)
	if err != nil {
		errCh <- err
		return nil, nil, errCh
	}

	// Get the io.WriteCloser and io.ReadCloser that we
	// will be returning
	wc, err := cmd.cmds[0].StdinPipe()
	if err != nil {
		errCh <- err
		return nil, nil, errCh
	}
	cmd.Stdin = cmd.cmds[0].Stdin

	rc, err := cmd.cmds[len(cmds)-1].StdoutPipe()
	if err != nil {
		errCh <- err
		return nil, nil, errCh
	}
	cmd.Stdout = cmd.cmds[len(cmds)-1].Stdout

	// Start our command chain and then wait in a goroutine
	// where we just monitor for errors
	if err := cmd.Start(); err != nil {
		errCh <- err
		return nil, nil, errCh
	}
	go func() {
		if err := cmd.Wait(); err != nil {
			errCh <- err
			return
		}
		// Close this if we finish waiting and have no errors
		close(errCh)
	}()

	// Return away
	return rc, wc, errCh
}

func NewPipeCommands(cmds ...*exec.Cmd) (*cmdPipe, error) {
	if len(cmds) == 0 {
		return nil, ErrNoCmds
	}

	// Link each command's Stdout w/ the next command's
	// Stdin
	writeClosers := make([]io.WriteCloser, len(cmds))
	for i := range cmds {
		if i == 0 {
			continue
		}
		r, w := io.Pipe()
		writeClosers[i-1] = w
		cmds[i-1].Stdout = w
		cmds[i].Stdin = r
	}

	return &cmdPipe{
		writeClosers: writeClosers,
		cmds:         cmds,
	}, nil
}

type cmdPipe struct {
	// Change these *BEFORE* calling start
	Stdin        io.Reader
	Stdout       io.Writer
	writeClosers []io.WriteCloser
	cmds         []*exec.Cmd
}

func (c *cmdPipe) Start() error {
	c.cmds[0].Stdin = c.Stdin
	c.cmds[len(c.cmds)-1].Stdout = c.Stdout
	for i := range c.cmds {
		if err := c.cmds[i].Start(); err != nil {
			return err
		}
	}
	return nil
}

func (c *cmdPipe) Wait() error {
	for i := range c.cmds {
		if err := c.cmds[i].Wait(); err != nil {
			return err
		}
		if i != len(c.cmds)-1 {
			if err := c.writeClosers[i].Close(); err != nil {
				return err
			}
		}
	}
	return nil
}