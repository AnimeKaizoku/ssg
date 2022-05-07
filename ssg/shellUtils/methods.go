package shellUtils

import (
	"time"
)

// WaitAndRun method waits for the execution to either finish or gets cancelled.
func (r *ExecuteCommandResult) WaitAndRun(
	interval, timeout time.Duration,
	handler ExecuteResultEventHandler,
) {
	if handler == nil {
		// prevent from panic
		return
	}

	started := time.Now()
	for !r.IsDone() {
		if time.Since(started) >= timeout {
			// execution exceeded the timeout
			go handler(r)
			return
		}
		time.Sleep(interval)
	}

	if r.Exited() && !r.IsFinished {
		r.IsFinished = true
	}

	go handler(r)
}

// Exited reports whether the program has exited.
// On Unix systems this reports true if the program exited due to calling exit,
// but false if the program terminated due to a signal.
func (r *ExecuteCommandResult) Exited() bool {
	if r.cmd == nil {
		return true
	}

	if r.cmd.ProcessState == nil {
		return false
	}

	return r.cmd.ProcessState.Exited()
}

func (r *ExecuteCommandResult) IsDone() bool {
	return r.IsFinished || r.IsKilled || r.IsReleased || r.Exited()
}

func (r *ExecuteCommandResult) UserTime() time.Duration {
	if r.cmd == nil || r.cmd.ProcessState == nil {
		return 0
	}

	return r.cmd.ProcessState.UserTime()
}

func (r *ExecuteCommandResult) Kill() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.cmd == nil || r.cmd.Process == nil {
		return nil
	}

	r.IsKilled = true
	r.IsFinished = true

	return r.cmd.Process.Kill()
}

func (r *ExecuteCommandResult) Release() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.cmd == nil || r.cmd.Process == nil {
		return nil
	}

	r.IsReleased = true
	r.IsFinished = true

	return r.cmd.Process.Release()
}
