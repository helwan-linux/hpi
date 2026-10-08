package system

import (
	"fmt"
	"os"
	"path/filepath"
)

const nemoScriptName = "hpi-install"

func NemoScriptPath() string {
	return filepath.Join(
		"/usr",
		"share",
		"nemo",
		"scripts",
		nemoScriptName,
	)
}

func NemoScript() string {
	return `#!/bin/sh

if [ -z "$NEMO_SCRIPT_SELECTED_URIS" ]; then
    exit 0
fi

for uri in $NEMO_SCRIPT_SELECTED_URIS; do
    file="$(printf '%s' "$uri" | sed 's#^file://##')"

    case "$file" in
        *.pkg.tar.zst)
            hpi "$file"
            ;;
    esac
done
`
}

func InstallNemoIntegration() error {
	destination := NemoScriptPath()

	if err := os.MkdirAll(
		filepath.Dir(destination),
		0755,
	); err != nil {
		return fmt.Errorf(
			"create Nemo scripts directory: %w",
			err,
		)
	}

	if err := os.WriteFile(
		destination,
		[]byte(NemoScript()),
		0755,
	); err != nil {
		return fmt.Errorf(
			"install Nemo integration: %w",
			err,
		)
	}

	return nil
}
