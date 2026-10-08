package packageinfo

import (
	"fmt"
	"strconv"
	"strings"
)

type Metadata struct {
	Name        string
	Version     string
	Description string
	Architecture string
	URL         string
	License     []string

	Depends    []string
	OptDepends []string
	Conflicts  []string
	Provides   []string
	Replaces   []string

	Packager string
	BuildDate string

	InstalledSize int64
	CompressedSize int64
}

func (m *Metadata) FullName() string {
	if m.Name == "" {
		return ""
	}

	if m.Version == "" {
		return m.Name
	}

	return m.Name + "-" + m.Version
}

func (m *Metadata) DependencyCount() int {
	return len(m.Depends)
}

func (m *Metadata) HasDependencies() bool {
	return len(m.Depends) > 0
}

func parsePKGINFO(data []byte) (*Metadata, error) {
	metadata := &Metadata{}

	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		switch key {
		case "pkgname":
			metadata.Name = value

		case "pkgver":
			metadata.Version = value

		case "pkgdesc":
			metadata.Description = value

		case "arch":
			metadata.Architecture = value

		case "url":
			metadata.URL = value

		case "license":
			metadata.License = append(
				metadata.License,
				value,
			)

		case "depend":
			metadata.Depends = append(
				metadata.Depends,
				value,
			)

		case "optdepend":
			metadata.OptDepends = append(
				metadata.OptDepends,
				value,
			)

		case "conflict":
			metadata.Conflicts = append(
				metadata.Conflicts,
				value,
			)

		case "provides":
			metadata.Provides = append(
				metadata.Provides,
				value,
			)

		case "replaces":
			metadata.Replaces = append(
				metadata.Replaces,
				value,
			)

		case "packager":
			metadata.Packager = value

		case "builddate":
			metadata.BuildDate = value

		case "size":
			size, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				metadata.InstalledSize = size
			}

		case "csize":
			size, err := strconv.ParseInt(value, 10, 64)
			if err == nil {
				metadata.CompressedSize = size
			}
		}
	}

	if metadata.Name == "" {
		return nil, fmt.Errorf("package metadata does not contain pkgname")
	}

	if metadata.Version == "" {
		return nil, fmt.Errorf("package metadata does not contain pkgver")
	}

	return metadata, nil
}
