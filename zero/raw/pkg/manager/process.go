package manager

import (
	"os/exec"
	"sync"
)

// Process is an abstraction over an exec.Cmd that provides some additional convenience
// functionalities.
type Process struct {
	waitExecutor sync.Once
	waitSync     sync.WaitGroup
	waitResult   error

	*exec.Cmd
}

// NewProcess creates a new process from an exec.Cmd.
func NewProcess(cmd *exec.Cmd) *Process {
	return &Process{Cmd: cmd}
}

// Wait is a thread-safe, idempotent implementation of a process wait. It is safe for concurrent
// usage and may be called multiple times. Concurrent calls block until the wait is complete, and
// calls after the process has exited return the result of the previously completed wait.
func (p *Process) Wait() error {
	// Wait() on a process can only be called once; this is ensured by sync.Once.
	p.waitExecutor.Do(func() {
		// The sync.WaitGroup provides effective mutual exclusion, as subsequent calls will
		// block until the wait result is available.
		p.waitSync.Add(1)
		p.waitResult = p.Cmd.Wait()
		p.waitSync.Done()
	})

	// This blocks if a process wait is in progress, and passes immediately if the wait has
	// already completed.
	p.waitSync.Wait()

	return p.waitResult
}

// WaitC is a convenience wrapper over Wait that provides a channel interface for consuming the
// process wait result.
func (p *Process) WaitC() <-chan error {
	ch := make(chan error)

	go func() {
		ch <- p.Wait()
	}()

	return ch
}

// Exited is a convenience function to indicate whether the process has exited.
func (p *Process) Exited() bool {
	return p.Cmd.ProcessState != nil
}
