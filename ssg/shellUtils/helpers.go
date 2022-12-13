package shellUtils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func RunCommand(command string) *ExecuteCommandResult {
	return runCommand(command, false, nil)
}

func RunPowerShell(command string) *ExecuteCommandResult {
	return runPowerShell(command, false, nil)
}

func RunCommandAsync(command string) *ExecuteCommandResult {
	return runCommand(command, true, nil)
}

func RunCommandAsyncWithChan(command string, finishedChan chan bool) *ExecuteCommandResult {
	return runCommand(command, true, finishedChan)
}

func RunPowerShellAsyncWithChan(command string, finishedChan chan bool) *ExecuteCommandResult {
	return runPowerShell(command, true, finishedChan)
}

func runCommand(command string, isAsync bool, finishedChan chan bool) *ExecuteCommandResult {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	var result *ExecuteCommandResult
	if isAsync {
		result = executeCommand(command, &ExecuteCommandConfig{
			TargetRunner:  GetCommandTargetRunner(),
			PrimaryArgs:   GetCommandPrimaryArgs(),
			Stdout:        stdout,
			Stderr:        stderr,
			autoSetOutput: true,
			FinishedChan:  finishedChan,
			IsAsync:       true,
		})
		return result
	} else {
		result = executeCommand(command, &ExecuteCommandConfig{
			TargetRunner: GetCommandTargetRunner(),
			PrimaryArgs:  GetCommandPrimaryArgs(),
			Stdout:       stdout,
			Stderr:       stderr,
		})
	}

	result.Stdout = stdout.String()
	result.Stderr = stderr.String()

	return result
}

func runPowerShell(command string, isAsync bool, finishedChan chan bool) *ExecuteCommandResult {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	var result *ExecuteCommandResult
	if isAsync {
		result = executePowerShell(command, &ExecuteCommandConfig{
			TargetRunner:           GetPowerShellRunner(),
			PrimaryArgs:            GetPowerShellPrimaryArgs(),
			Stdout:                 stdout,
			Stderr:                 stderr,
			autoSetOutput:          true,
			FinishedChan:           finishedChan,
			IsAsync:                true,
			RemovePowerShellPrompt: true,
		})
		return result
	} else {
		result = executePowerShell(command, &ExecuteCommandConfig{
			TargetRunner:           GetPowerShellRunner(),
			PrimaryArgs:            GetPowerShellPrimaryArgs(),
			Stdout:                 stdout,
			Stderr:                 stderr,
			RemovePowerShellPrompt: true,
		})
	}

	result.Stdout = stdout.String()
	result.Stderr = stderr.String()

	return result
}

func ExecuteCommand(command string, config *ExecuteCommandConfig) *ExecuteCommandResult {
	if config == nil {
		config = &ExecuteCommandConfig{
			TargetRunner: GetCommandTargetRunner(),
			PrimaryArgs:  GetCommandPrimaryArgs(),
		}
	} else if config.IsAsync {
		config.IsAsync = false
	}

	return executeCommand(command, config)
}

func ExecuteCommandAsync(command string, config *ExecuteCommandConfig) *ExecuteCommandResult {
	if config == nil {
		config = &ExecuteCommandConfig{
			TargetRunner: GetCommandTargetRunner(),
			PrimaryArgs:  GetCommandPrimaryArgs(),
			IsAsync:      true,
		}
	}

	return executeCommand(command, config)
}

func ExecutePowerShellAsync(command string, config *ExecuteCommandConfig) *ExecuteCommandResult {
	if config == nil {
		config = &ExecuteCommandConfig{
			TargetRunner: GetPowerShellRunner(),
			PrimaryArgs:  GetPowerShellPrimaryArgs(),
			IsAsync:      true,
		}
	}

	return executePowerShell(command, config)
}

func GetCommandTargetRunner() string {
	if os.PathSeparator == '/' {
		return ShellToUseUnix
	} else {
		return ShellToUseWin
	}
}

func GetPowerShellRunner() string {
	return PowerShellCmd
}

func GetCommandPrimaryArgs() []string {
	if os.PathSeparator == '/' {
		return []string{"-c"}
	}
	return []string{"/c"}
}

func GetPowerShellPrimaryArgs() []string {
	return []string{"-nologo", "-noprofile", "-NonInteractive"}
}

// executeCommand is the internal version of the execute command function.
// WARNING: the config argument MUST NOT be nil.
func executeCommand(
	command string,
	config *ExecuteCommandConfig,
) *ExecuteCommandResult {
	var cmd *exec.Cmd
	result := &ExecuteCommandResult{
		autoSetOutput: config.autoSetOutput,
		mutex:         &sync.Mutex{},
	}

	cmd = exec.Command(config.TargetRunner, append(config.PrimaryArgs, command)...)
	if len(config.ExtraFiles) != 0 {
		cmd.ExtraFiles = append(cmd.ExtraFiles, config.ExtraFiles...)
	}

	cmd.Stdout = config.Stdout
	cmd.Stdin = config.Stdin
	cmd.Stderr = config.Stderr
	cmd.Args = append(cmd.Args, config.AdditionalArgs...)
	cmd.Env = append(cmd.Env, config.AdditionalEnv...)

	result.cmd = cmd
	result.FinishedChan = config.FinishedChan

	finishUpCommand(cmd, config, result)

	return result
}

// executeCommand is the internal version of the execute powershell function.
// executes the given powershell script/command/set of commands using the
// "powershell" (it might be powershell 5.1 which ships with windows by default).
// WARNING: the config argument MUST NOT be nil.
func executePowerShell(
	command string,
	config *ExecuteCommandConfig,
) *ExecuteCommandResult {
	var cmd *exec.Cmd
	result := &ExecuteCommandResult{
		autoSetOutput: config.autoSetOutput,
		mutex:         &sync.Mutex{},
	}

	if config.RemovePowerShellPrompt && !strings.Contains(command, "function prompt") {
		// hacky way of getting rid of powershell prompt
		command = PowerShellPromptOverride + "\n" + command
	}

	cmd = exec.Command(config.TargetRunner, config.PrimaryArgs...)
	if len(config.ExtraFiles) != 0 {
		cmd.ExtraFiles = append(cmd.ExtraFiles, config.ExtraFiles...)
	}

	cmd.Stdout = config.Stdout

	pStdin, err := cmd.StdinPipe()
	if err != nil {
		result.Error = err
		return result
	}

	if config.Stdin != nil {
		wrappedStdin := &StdinWrapper{
			InnerWriter: pStdin,
		}

		result.pipedStdin = wrappedStdin
		wrappedStdin.OnWrite = append(wrappedStdin.OnWrite, func(p []byte) (n int, err error) {
			return config.Stdin.Read(p)
		})
	} else {
		result.pipedStdin = pStdin
	}

	cmd.Stderr = config.Stderr
	cmd.Args = append(cmd.Args, config.AdditionalArgs...)
	cmd.Env = append(cmd.Env, config.AdditionalEnv...)

	result.cmd = cmd
	result.FinishedChan = config.FinishedChan
	_, err = fmt.Fprint(result.pipedStdin, command)
	if err != nil {
		result.Error = err
		return result
	}

	finishUpCommand(cmd, config, result)
	return result
}

func finishUpCommand(
	cmd *exec.Cmd,
	config *ExecuteCommandConfig,
	result *ExecuteCommandResult,
) {
	if config.IsAsync {
		go func() {
			result.ClosePipes()
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
		result.ClosePipes()
		result.Error = cmd.Run()
		result.IsFinished = true
	}

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
