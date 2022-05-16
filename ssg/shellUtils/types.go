package shellUtils

import (
	"io"
	"os"
	"os/exec"
	"sync"
)

type ExecuteCommandConfig struct {
	Stdin          io.Reader
	Stdout         io.Writer
	Stderr         io.Writer
	AdditionalArgs []string
	AdditionalEnv  []string
	ExtraFiles     []*os.File
	FinishedChan   chan bool

	autoSetOutput bool
}

type ExecuteCommandResult struct {
	// UniqueId is a unique id assigned to this execution result.
	// setting this field is up to the user and library itself won't
	// do anything with it.
	UniqueId string
	// Stdout field is the string representation of the output from
	// command. it's set only if you use `RunCommand` or `RunCommandAsync`
	// functions.
	Stdout string
	// Stderr field is the string representation of the err-output from
	// command. it's set only if you use `RunCommand` or `RunCommandAsync`
	// functions.
	Stderr string
	// Error field is set only after execution of the command finishes and
	// it can be nil.
	Error error
	// IsKilled field is set to true only if the `Kill` method is called.
	IsKilled     bool
	IsFinished   bool
	IsReleased   bool
	FinishedChan chan bool

	autoSetOutput bool
	cmd           *exec.Cmd
	mutex         *sync.Mutex
}

type ExecuteResultEventHandler func(result *ExecuteCommandResult)
