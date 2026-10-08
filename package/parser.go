package packageinfo

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/klauspost/compress/zstd"
)

const (
	PackageSuffix = ".pkg.tar.zst"
	PKGINFOPath   = ".PKGINFO"
)

type Package struct {
	Path     string
	Metadata *Metadata
}

func Open(path string) (*Package, error) {
	if strings.TrimSpace(path) == "" {
		return nil, fmt.Errorf("package path is empty")
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("cannot access package: %w", err)
	}

	if info.IsDir() {
		return nil, fmt.Errorf("path is a directory, not a package")
	}

	if !strings.HasSuffix(path, PackageSuffix) {
		return nil, fmt.Errorf(
			"unsupported package format: expected %s",
			PackageSuffix,
		)
	}

	metadataData, err := readPKGINFO(path)
	if err != nil {
		return nil, err
	}

	metadata, err := parsePKGINFO(metadataData)
	if err != nil {
		return nil, fmt.Errorf("invalid package metadata: %w", err)
	}

	return &Package{
		Path:     path,
		Metadata: metadata,
	}, nil
}

func readPKGINFO(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open package: %w", err)
	}

	defer file.Close()

	decoder, err := zstd.NewReader(file)
	if err != nil {
		return nil, fmt.Errorf("open zstd stream: %w", err)
	}

	defer decoder.Close()

	reader := tar.NewReader(decoder)

	for {
		header, err := reader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf(
				"read package archive: %w",
				err,
			)
		}

		if header.Name != PKGINFOPath {
			continue
		}

		if header.Typeflag != tar.TypeReg {
			return nil, fmt.Errorf(
				"%s is not a regular file",
				PKGINFOPath,
			)
		}

		data, err := io.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf(
				"read %s: %w",
				PKGINFOPath,
				err,
			)
		}

		return data, nil
	}

	return nil, fmt.Errorf(
		"package does not contain %s",
		PKGINFOPath,
	)
}
