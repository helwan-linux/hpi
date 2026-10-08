package installer

import (
	"context"
	"fmt"
	"strings"

	packageinfo "github.com/helwan-linux/hpi/package"
)

type Transaction struct {
	Pacman     *Pacman
	Resolver   *Resolver
	Permission *PermissionManager
}

func NewTransaction() *Transaction {
	pacman := NewPacman()

	return &Transaction{
		Pacman:     pacman,
		Resolver:   NewResolver(),
		Permission: NewPermissionManager(),
	}
}

func (t *Transaction) Check(
	ctx context.Context,
	pkg *packageinfo.Package,
) (*Resolution, error) {
	if pkg == nil {
		return nil, fmt.Errorf("package is nil")
	}

	if t.Resolver == nil {
		t.Resolver = NewResolver()
	}

	return t.Resolver.Resolve(ctx, pkg)
}

func (t *Transaction) Install(
	ctx context.Context,
	packagePath string,
) (*CommandResult, error) {
	if packagePath == "" {
		return nil, fmt.Errorf("package path is empty")
	}

	if t.Pacman == nil {
		t.Pacman = NewPacman()
	}

	if t.Permission == nil {
		t.Permission = NewPermissionManager()
	}

	if err := t.Pacman.CheckAvailable(); err != nil {
		return nil, err
	}

	/*
		pacman -U is intentionally used as one transaction.

		This lets pacman handle the package installation
		as a single operation instead of HPI installing
		packages one by one.

		--needed avoids reinstalling packages that are
		already present.

		--noconfirm is used because HPI is a graphical
		installer and cannot depend on an interactive
		terminal prompt.
	*/

	result, err := t.Permission.Run(
		ctx,
		t.Pacman.Command,
		"-U",
		"--needed",
		"--noconfirm",
		packagePath,
	)

	if err != nil {
		message := ""

		if result != nil {
			message = strings.TrimSpace(result.Error)

			if message == "" {
				message = strings.TrimSpace(result.Output)
			}
		}

		if message == "" {
			message = err.Error()
		}

		return result, fmt.Errorf(
			"package installation failed: %s",
			message,
		)
	}

	return result, nil
}
