package shellUtils

const (
	ShellToUseUnix = "bash"
	ShellToUseWin  = "cmd"
	gitCmd         = "git rev-parse --short HEAD &&" +
		" git rev-parse --verify HEAD &&" +
		" git fetch && " +
		"git rev-list --left-right --count origin/master...master"
	PowerShellCmd            = "powershell"
	PromptIgnoreSSG          = "PROMPT_IGNORE_SSG>"
	PowerShellPromptOverride = "function prompt { \"" + PromptIgnoreSSG + "\" }"
)
