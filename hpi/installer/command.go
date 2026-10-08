package installer

import (
	"bytes"
	"os/exec"
)

func runCommand(
	command *exec.Cmd,
) (*CommandResult, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()

	result := &CommandResult{
		Output: stdout.String(),
		Error:  stderr.String(),
	}

	if command.ProcessState != nil {
		result.ExitCode = command.ProcessState.ExitCode()
	}

	return result, err
}
