package tools

import (
	"errors"
	"io"
	"os/exec"
)

var (
	ErrNoCmds = errors.New("pipe: there were no commands provided")
)

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
