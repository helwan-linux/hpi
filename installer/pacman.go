package installer

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
)

type Pacman struct {
	Command string
}

type CommandResult struct {
	ExitCode int
	Output   string
	Error    string
}

func NewPacman() *Pacman {
	return &Pacman{
		Command: "pacman",
	}
}

func (p *Pacman) CheckAvailable() error {
	if _, err := exec.LookPath(p.Command); err != nil {
		return fmt.Errorf(
			"pacman was not found on this system: %w",
			err,
		)
	}

	return nil
}

func (p *Pacman) Run(
	ctx context.Context,
	args ...string,
) (*CommandResult, error) {
	if err := p.CheckAvailable(); err != nil {
		return nil, err
	}

	command := exec.CommandContext(
		ctx,
		p.Command,
		args...,
	)

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

	if err != nil {
		return result, err
	}

	return result, nil
}
