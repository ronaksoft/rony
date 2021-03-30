package tools

import (
	"errors"
	"io"
	"os/exec"
)

var (
	ErrNoCommands = errors.New("pipe: there were no commands provided")
)

func NewPipeCommands(cmds ...*exec.Cmd) (*cmdPipe, error) {
	if len(cmds) == 0 {
		return nil, ErrNoCommands
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
		commands:     cmds,
	}, nil
}

type cmdPipe struct {
	// Change these *BEFORE* calling start
	Stdin        io.Reader
	Stdout       io.Writer
	writeClosers []io.WriteCloser
	commands     []*exec.Cmd
}

func (c *cmdPipe) Start() error {
	c.commands[0].Stdin = c.Stdin
	c.commands[len(c.commands)-1].Stdout = c.Stdout
	for i := range c.commands {
		if err := c.commands[i].Start(); err != nil {
			return err
		}
	}
	return nil
}

func (c *cmdPipe) Wait() error {
	for i := range c.commands {
		if err := c.commands[i].Wait(); err != nil {
			return err
		}
		if i != len(c.commands)-1 {
			if err := c.writeClosers[i].Close(); err != nil {
				return err
			}
		}
	}
	return nil
}
