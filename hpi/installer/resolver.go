package installer

import (
	"context"
	"fmt"
	"strings"

	packageinfo "github.com/helwan-linux/hpi/package"
)

type DependencyStatus int

const (
	DependencyInstalled DependencyStatus = iota
	DependencyAvailable
	DependencyMissing
	DependencyUnknown
)

type DependencyResult struct {
	Dependency packageinfo.Dependency
	Status     DependencyStatus
	Reason     string
}

type Resolution struct {
	Dependencies []DependencyResult

	Installed []DependencyResult
	Available []DependencyResult
	Missing   []DependencyResult
	Unknown   []DependencyResult

	Online bool
}

func (r *Resolution) CanInstallOffline() bool {
	return len(r.Missing) == 0 && len(r.Available) == 0
}

func (r *Resolution) HasMissing() bool {
	return len(r.Missing) > 0
}

func (r *Resolution) HasRepositoryPackages() bool {
	return len(r.Available) > 0
}

func (r *Resolution) RequiresInternet() bool {
	return len(r.Available) > 0
}

func (r *Resolution) IsComplete() bool {
	return len(r.Missing) == 0 && len(r.Unknown) == 0
}

type Resolver struct {
	Pacman *Pacman
}

func NewResolver() *Resolver {
	return &Resolver{
		Pacman: NewPacman(),
	}
}

func (r *Resolver) Resolve(
	ctx context.Context,
	pkg *packageinfo.Package,
) (*Resolution, error) {
	if pkg == nil {
		return nil, fmt.Errorf("package is nil")
	}

	if r.Pacman == nil {
		r.Pacman = NewPacman()
	}

	if err := r.Pacman.CheckAvailable(); err != nil {
		return nil, err
	}

	resolution := &Resolution{
		Dependencies: make([]DependencyResult, 0),
		Installed:    make([]DependencyResult, 0),
		Available:    make([]DependencyResult, 0),
		Missing:      make([]DependencyResult, 0),
		Unknown:      make([]DependencyResult, 0),
	}

	online, err := r.detectOnline(ctx)
	if err != nil {
		return nil, err
	}

	resolution.Online = online

	for _, dependency := range pkg.Dependencies() {
		result, err := r.checkDependency(
			ctx,
			dependency,
			online,
		)
		if err != nil {
			return nil, err
		}

		resolution.Dependencies = append(
			resolution.Dependencies,
			*result,
		)

		switch result.Status {
		case DependencyInstalled:
			resolution.Installed = append(
				resolution.Installed,
				*result,
			)

		case DependencyAvailable:
			resolution.Available = append(
				resolution.Available,
				*result,
			)

		case DependencyMissing:
			resolution.Missing = append(
				resolution.Missing,
				*result,
			)

		case DependencyUnknown:
			resolution.Unknown = append(
				resolution.Unknown,
				*result,
			)
		}
	}

	return resolution, nil
}

func (r *Resolver) checkDependency(
	ctx context.Context,
	dependency packageinfo.Dependency,
	online bool,
) (*DependencyResult, error) {
	result := &DependencyResult{
		Dependency: dependency,
	}

	installed, err := r.isInstalled(ctx, dependency)
	if err != nil {
		return nil, err
	}

	if installed {
		result.Status = DependencyInstalled
		result.Reason = "dependency is already installed"
		return result, nil
	}

	if !online {
		result.Status = DependencyMissing
		result.Reason = "dependency is not installed and the system is offline"
		return result, nil
	}

	available, err := r.repositoryHasPackage(
		ctx,
		dependency,
	)
	if err != nil {
		return nil, err
	}

	if available {
		result.Status = DependencyAvailable
		result.Reason = "dependency is available from configured repositories"
		return result, nil
	}

	result.Status = DependencyMissing
	result.Reason = "dependency is not installed and is not available from configured repositories"

	return result, nil
}

func (r *Resolver) isInstalled(
	ctx context.Context,
	dependency packageinfo.Dependency,
) (bool, error) {
	result, err := r.Pacman.Run(
		ctx,
		"-T",
		dependency.RawName,
	)

	if result == nil && err != nil {
		return false, err
	}

	if result != nil && result.ExitCode == 0 {
		return true, nil
	}

	return false, nil
}

func (r *Resolver) repositoryHasPackage(
	ctx context.Context,
	dependency packageinfo.Dependency,
) (bool, error) {
	result, err := r.Pacman.Run(
		ctx,
		"-Sp",
		"--print-format",
		"%n",
		dependency.RawName,
	)

	if result == nil && err != nil {
		return false, err
	}

	if result == nil {
		return false, nil
	}

	if result.ExitCode != 0 {
		return false, nil
	}

	return strings.TrimSpace(result.Output) != "", nil
}

func (r *Resolver) detectOnline(
	ctx context.Context,
) (bool, error) {
	/*
		pacman -Sp performs a repository lookup without
		downloading or installing anything.

		If pacman can successfully resolve a package from
		the configured repositories, repository/network
		access is available.

		"pacman" itself is used only as a connectivity
		probe because it is expected to exist on an
		Arch-based system.
	*/

	result, err := r.Pacman.Run(
		ctx,
		"-Sp",
		"--print-format",
		"%n",
		"pacman",
	)

	if result == nil && err != nil {
		return false, err
	}

	if result == nil {
		return false, nil
	}

	return result.ExitCode == 0 &&
		strings.TrimSpace(result.Output) != "", nil
}
