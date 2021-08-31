package manager

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"zero/internal/config"
)

// Manager is a process manager that mediates network sockets and manages graceful process reloads.
type Manager struct {
	cfg config.Service
	cmd func() *Process

	active *Process
	mutex  sync.Mutex
}

// NewManager creates a new process manager.
func NewManager(cfg config.Service) (*Manager, error) {
	var env []string
	var files []*os.File
	var maxFD int

	for _, listener := range cfg.Listeners {
		if listener.SocketFD > maxFD {
			maxFD = listener.SocketFD
		}
	}

	for _, envvar := range cfg.Runtime.Environment {
		if envvar.Value != "" {
			env = append(env, fmt.Sprintf("%s=%s", envvar.Key, envvar.Value))
		} else {
			env = append(env, fmt.Sprintf("%s=%s", envvar.Key, os.Getenv(envvar.Key)))
		}
	}

	files = make([]*os.File, maxFD-2) // FDs 0-2 are reserved for stdin, stdout, and stderr

	for _, listener := range cfg.Listeners {
		var file *os.File

		ln, err := listener.Address.Listen()
		if err != nil {
			return nil, err
		}

		switch tln := ln.(type) {
		case *net.TCPListener:
			file, err = tln.File()
			if err != nil {
				return nil, err
			}
		case *net.UnixListener:
			file, err = tln.File()
			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf(
				"manager: unsupported listener type: name=%s address=%v type=%T",
				listener.Name,
				listener.Address,
				ln,
			)
		}

		// The ith index in ExtraFiles corresponds to socket FD i + 3.
		files[listener.SocketFD-3] = file
	}

	cmd := func() *Process {
		return NewProcess(&exec.Cmd{
			Dir:        cfg.Runtime.Directory,
			Path:       cfg.Runtime.Path,
			Args:       append([]string{cfg.Runtime.Path}, cfg.Runtime.Args...),
			Env:        env,
			ExtraFiles: files,
			Stdin:      os.Stdin,
			Stdout:     os.Stdout,
			Stderr:     os.Stderr,
			// Ensure that the service (child) process is always killed if Zero itself
			// (parent) unexpectedly exits.
			SysProcAttr: &syscall.SysProcAttr{Pdeathsig: syscall.SIGKILL},
		})
	}

	return &Manager{
		cfg: cfg,
		cmd: cmd,
	}, nil
}

// Start starts the process asynchronously.
func (m *Manager) Start() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.active != nil {
		return fmt.Errorf("manager: cannot start when an active process exists")
	}

	m.active = m.cmd()

	if err := m.active.Start(); err != nil {
		return fmt.Errorf("manager: error starting service process: err=%v", err)
	}

	select {
	case err := <-m.active.WaitC():
		if err != nil {
			return fmt.Errorf("manager: unexpected service process exit: err=%v", err)
		}

		return fmt.Errorf(
			"manager: process exited before required threshold uptime: uptime=%v",
			m.cfg.Options.ReloadUptime,
		)
	case <-time.After(m.cfg.Options.ReloadUptime):
	}

	return nil
}

// Reload gracefully reloads the process by launching a new process, followed by shutting down the
// active process. It returns an error if there is an error starting the new process or shutting
// down the old process, and an *os.ProcessState describing the exited process if the active process
// was successfully replaced.
func (m *Manager) Reload() (*os.ProcessState, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.active == nil {
		return nil, fmt.Errorf("manager: cannot reload with no current active process")
	}

	pending := m.cmd()
	if err := pending.Start(); err != nil {
		return nil, fmt.Errorf("manager: error starting new process instance: err=%v", err)
	}

	select {
	case err := <-pending.WaitC():
		if err != nil {
			return nil, fmt.Errorf(
				"manager: unexpected process instance exit: err=%v",
				err,
			)
		}

		return nil, fmt.Errorf(
			"manager: process exited before required threshold uptime: uptime=%v",
			m.cfg.Options.ReloadUptime,
		)
	case <-time.After(m.cfg.Options.ReloadUptime):
	}

	defer func() {
		// Regardless of whether the shutdown was successful, the active process should be
		// swapped since its successor is already running.
		m.active = pending
	}()

	if err := m.shutdown(); err != nil {
		return m.active.ProcessState, fmt.Errorf(
			"manager: error shutting down active process instance: err=%v",
			err,
		)
	}

	return m.active.ProcessState, nil
}

// Close shuts down the manager by shutting down the currently active process. It obtains an
// exclusive lock to avoid contention with an in-progress reload.
func (m *Manager) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if err := m.shutdown(); err != nil {
		return fmt.Errorf(
			"manager: error shutting down active process instance: err=%v",
			err,
		)
	}

	m.active = nil

	return nil
}

// Process returns the currently active process.
func (m *Manager) Process() (*os.Process, error) {
	if m.active == nil {
		return nil, fmt.Errorf("manager: no active process")
	} else if m.active.Process == nil {
		return nil, fmt.Errorf("manager: process is not yet started")
	}

	return m.active.Process, nil
}

// shutdown shuts down the currently active process. It attempts to gracefully shut down the process
// using its graceful shutdown signal, but will forcefully kill the process if it is still alive
// past its shutdown timeout.
func (m *Manager) shutdown() error {
	if err := m.active.Process.Signal(m.cfg.Options.ShutdownSignal.Signal); err != nil {
		return err
	}

	// Synchronously wait for the process to exit if there is no configured shutdown timeout.
	if m.cfg.Options.ShutdownTimeout == 0 {
		return m.active.Wait()
	}

	select {
	case err := <-m.active.WaitC():
		if err != nil {
			return err
		}
	case <-time.After(m.cfg.Options.ShutdownTimeout):
		// Process has not exited after the configured timeout grace period; kill the
		// process immediately
		return m.active.Process.Kill()
	}

	return nil
}
