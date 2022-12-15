package shellUtils

import (
	"strings"
	"time"
	"unicode"
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

func (r *ExecuteCommandResult) ClosePipes() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.pipedStdin != nil {
		err := r.pipedStdin.Close()
		r.pipedStdin = nil
		if err != nil {
			return err
		}
	}

	return nil
}

// PurifyOutput will purify the stdout of the result and will
// set the purified version of the output. If you just want to
// get the purified version without setting the new value, use
func (r *ExecuteCommandResult) PurifyPowerShellOutput() string {
	r.Stdout = r.GetPurifiedPowerShellOutput()
	return r.Stdout
}

func (r *ExecuteCommandResult) GetPurifiedPowerShellOutput() string {
	myStrs := strings.Split(r.Stdout, "\n")
	if len(myStrs) == 0 {
		return ""
	}

	// if the override doesn't exist, don't do any operation, just
	// return the original value.
	if !strings.Contains(myStrs[0], PowerShellPromptOverride) {
		return r.Stdout
	}

	purifiedOutput := &strings.Builder{}
	for i, currentLine := range myStrs {
		if i == 0 {
			// it's guaranteed that the first line is our prompt override.
			continue
		}

		if strings.HasPrefix(currentLine, PromptIgnoreSSG) {
			continue
		}

		if strings.HasPrefix(currentLine, ">>") {
			// looks like something we should be removing
			currentLine = strings.TrimLeft(currentLine, ">")
			if unicode.IsSpace(rune(currentLine[0])) {
				continue
			}
		}

		purifiedOutput.WriteString(currentLine)
		purifiedOutput.WriteRune('\n')
	}

	return strings.TrimSpace(purifiedOutput.String())
}

// --------------------------------------------------------

func (r *StdinWrapper) Write(p []byte) (n int, err error) {
	if len(r.OnWrite) > 0 {
		for _, currentHandler := range r.OnWrite {
			n, err = currentHandler(p)
			if err != nil {
				return n, err
			}
		}
	}

	return r.InnerWriter.Write(p)
}

func (r *StdinWrapper) Close() error {
	if len(r.OnClose) > 0 {
		var err error
		for _, currentHandler := range r.OnClose {
			err = currentHandler()
			if err != nil {
				return err
			}
		}
	}

	return r.InnerWriter.Close()
}
