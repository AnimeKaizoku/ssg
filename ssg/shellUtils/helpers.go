package shellUtils

import (
	"bytes"
	"os"
	"os/exec"
	"sync"
)

func RunCommand(command string) *ExecuteCommandResult {
	return runCommand(command, false, nil)
}

func RunCommandAsync(command string) *ExecuteCommandResult {
	return runCommand(command, true, nil)
}

func RunCommandAsyncWithChan(command string, finishedChan chan bool) *ExecuteCommandResult {
	return runCommand(command, true, finishedChan)
}

func runCommand(command string, isAsync bool, finishedChan chan bool) *ExecuteCommandResult {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	var result *ExecuteCommandResult
	if isAsync {
		result = executeCommand(command, &ExecuteCommandConfig{
			Stdout:        stdout,
			Stderr:        stderr,
			autoSetOutput: true,
			FinishedChan:  finishedChan,
		}, true)
		return result
	} else {
		result = ExecuteCommand(command, &ExecuteCommandConfig{
			Stdout: stdout,
			Stderr: stderr,
		})
	}

	result.Stdout = stdout.String()
	result.Stderr = stderr.String()

	return result
}

func ExecuteCommand(command string, config *ExecuteCommandConfig) *ExecuteCommandResult {
	return executeCommand(command, config, false)
}

func ExecuteCommandAsync(command string, config *ExecuteCommandConfig) *ExecuteCommandResult {
	return executeCommand(command, config, true)
}

func executeCommand(
	command string,
	config *ExecuteCommandConfig,
	isAsync bool,
) *ExecuteCommandResult {
	if config == nil {
		config = &ExecuteCommandConfig{}
	}

	var cmd *exec.Cmd
	result := &ExecuteCommandResult{
		autoSetOutput: config.autoSetOutput,
		mutex:         &sync.Mutex{},
	}
	if os.PathSeparator == '/' {
		cmd = exec.Command(ShellToUseUnix, "-c", command)
		cmd.ExtraFiles = append(cmd.ExtraFiles, config.ExtraFiles...)
	} else {
		cmd = exec.Command(ShellToUseWin, "/C", command)
	}

	cmd.Stdout = config.Stdout
	cmd.Stdin = config.Stdin
	cmd.Stderr = config.Stderr
	cmd.Args = append(cmd.Args, config.AdditionalArgs...)
	cmd.Env = append(cmd.Env, config.AdditionalEnv...)

	result.cmd = cmd

	if isAsync {
		go func() {
			result.Error = cmd.Run()
			result.IsFinished = true
			if result.autoSetOutput {
				stdout, ok := result.cmd.Stdout.(*bytes.Buffer)
				if ok && stdout != nil {
					result.Stdout = stdout.String()
				}

				stderr, ok := result.cmd.Stderr.(*bytes.Buffer)
				if ok && stderr != nil {
					result.Stderr = stderr.String()
				}
			}

			if result.FinishedChan != nil {
				result.FinishedChan <- true
			}
		}()
	} else {
		result.Error = cmd.Run()
		result.IsFinished = true
	}

	return result
}

// GetGitStats function will return the git stats in the following format:
// "d8e6e45 \n d8e6e45d52f7bf164a995e22abb81ffc6e3eeae1 \n 3 0"
func GetGitStatsString() string {
	result := RunCommand(gitCmd)
	stdout, err := result.Stdout, result.Error
	if err != nil || len(stdout) == 0 {
		return ""
	}

	return stdout
}

func StartProcess(args ...string) (p *os.Process, err error) {
	if args[0], err = exec.LookPath(args[0]); err == nil {
		var procAttr os.ProcAttr
		procAttr.Files = []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		}
		p, err := os.StartProcess(args[0], args, &procAttr)
		if err == nil {
			return p, nil
		}
	}
	return nil, err
}
