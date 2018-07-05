package watcher

import (
	"io"
)

// Process provides an interface for managing it
//go:generate counterfeiter -o ../mocks/watcher_process.go --fake-name Process . Process
type Process interface {
	SetOut(io.Writer, io.Writer)
	Start() error
	Wait() error
	Stop() error
}

// OSProcess implements process using default os processes
type OSProcess struct {
}
