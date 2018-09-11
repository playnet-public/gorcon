package watcher

import (
	"io"
	"os/exec"
)

// Process provides an interface for managing it
//go:generate counterfeiter -o ../mocks/watcher_process.go --fake-name Process . Process
type Process interface {
	SetOut(io.Writer, io.Writer)
	Run() error
	Stop() error
}

// OSProcess implements process using default os processes
type OSProcess struct {
	Cmd *exec.Cmd
}

// NewOSProcess with command and args
func NewOSProcess(cmd string, args ...string) *OSProcess {
	return &OSProcess{
		Cmd: exec.Command(cmd, args...),
	}
}

// SetOut sets the process stdout and stderr values to the passed in writers
func (p *OSProcess) SetOut(stderr, stdout io.Writer) {
	p.Cmd.Stderr = stderr
	p.Cmd.Stdout = stdout
}

// Run a new process and return error once it ends
func (p *OSProcess) Run() error {
	return p.Cmd.Run()
}

// Stop the current process by sending a termination signal
func (p *OSProcess) Stop() error {
	return nil
}
