package shellUtils

import (
	"io"
	"os"
	"os/exec"
	"sync"
)

type ExecuteCommandConfig struct {
	TargetRunner   string
	PrimaryArgs    []string
	Stdin          io.Reader
	Stdout         io.Writer
	Stderr         io.Writer
	AdditionalArgs []string
	AdditionalEnv  []string

	// ExtraFiles specifies additional open files to be inherited by the new process.
	// It does not include standard input, standard output, or standard error.
	// If non-nil, entry i becomes file descriptor 3+i.
	// ExtraFiles is not supported on Windows.
	ExtraFiles []*os.File

	FinishedChan chan bool

	// IsAsync field determines whether the command should be run
	// async or not. if this field is set to false, the execute
	// function will block the current goroutine until the process
	// ends completely.
	IsAsync                bool
	RemovePowerShellPrompt bool

	// autoSetOutput determines whether the output reader should
	// set automatically or not.
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
	pipedStdin    io.WriteCloser
	mutex         *sync.Mutex
}

type StdinWrapper struct {
	InnerWriter io.WriteCloser

	OnWrite []func(p []byte) (n int, err error)
	OnClose []func() error
}

type ExecuteResultEventHandler func(result *ExecuteCommandResult)
