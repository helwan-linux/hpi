package installer

import (
	"fmt"
	"os/exec"
)

func findExecutable(name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("executable name is empty")
	}

	path, err := exec.LookPath(name)
	if err != nil {
		return "", err
	}

	return path, nil
}
