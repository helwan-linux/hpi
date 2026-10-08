package system

import (
	"fmt"
	"os"
	"path/filepath"
)

const desktopFileName = "hpi.desktop"

func DesktopFilePath() string {
	return filepath.Join(
		"/usr",
		"share",
		"applications",
		desktopFileName,
	)
}

func InstallDesktopEntry(source string) error {
	if source == "" {
		return fmt.Errorf("desktop entry source is empty")
	}

	data, err := os.ReadFile(source)
	if err != nil {
		return fmt.Errorf(
			"read desktop entry: %w",
			err,
		)
	}

	destination := DesktopFilePath()

	if err := os.MkdirAll(
		filepath.Dir(destination),
		0755,
	); err != nil {
		return fmt.Errorf(
			"create applications directory: %w",
			err,
		)
	}

	if err := os.WriteFile(
		destination,
		data,
		0644,
	); err != nil {
		return fmt.Errorf(
			"install desktop entry: %w",
			err,
		)
	}

	return nil
}
