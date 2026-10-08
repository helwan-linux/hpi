package installer

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

type PermissionManager struct {
	Helper string
}

func NewPermissionManager() *PermissionManager {
	return &PermissionManager{
		Helper: "pkexec",
	}
}

func (p *PermissionManager) IsRoot() bool {
	return os.Geteuid() == 0
}

func (p *PermissionManager) RequiresElevation() bool {
	return !p.IsRoot()
}

func (p *PermissionManager) CheckHelper() error {
	if p.IsRoot() {
		return nil
	}

	if p.Helper == "" {
		return fmt.Errorf("no privilege helper configured")
	}

	if _, err := findExecutable(p.Helper); err != nil {
		return fmt.Errorf(
			"required privilege helper %q was not found: %w",
			p.Helper,
			err,
		)
	}

	return nil
}

func (p *PermissionManager) Run(
	ctx context.Context,
	command string,
	args ...string,
) (*CommandResult, error) {
	if command == "" {
		return nil, fmt.Errorf("command is empty")
	}

	if err := p.CheckHelper(); err != nil {
		return nil, err
	}

	var cmd *exec.Cmd

	if p.IsRoot() {
		cmd = exec.CommandContext(
			ctx,
			command,
			args...,
		)
	} else {
		privilegedArgs := make([]string, 0, len(args)+1)
		privilegedArgs = append(privilegedArgs, command)
		privilegedArgs = append(privilegedArgs, args...)

		cmd = exec.CommandContext(
			ctx,
			p.Helper,
			privilegedArgs...,
		)
	}

	output, err := runCommand(cmd)

	if err != nil {
		return output, err
	}

	return output, nil
}
