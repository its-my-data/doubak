package util

import (
	"flag"
	"github.com/its-my-data/doubak/proto"
	"os"
	"path/filepath"
)

const CollectorPathPrefix = "collector/"

// GetPathWithCreation returns the concatenated path with output path and have them created in advance.
func GetPathWithCreation(subdirs string) (string, error) {
	baseDir := flag.Lookup(proto.Flag_output_dir.String()).Value.String()
	return GetPathWithCreationWithBase(baseDir, subdirs)
}

// GetPathWithCreationWithBase returns the concatenated path with base path and have them created in advance.
func GetPathWithCreationWithBase(base, subdirs string) (string, error) {
	// Try to create the base directory first if not exists.
	if _, err := os.Stat(base); os.IsNotExist(err) {
		// Only create the top-most base directory.
		err := os.Mkdir(base, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	// Recursively create sub-directories.
	newPath := filepath.Join(base, subdirs)
	return newPath, os.MkdirAll(newPath, os.ModePerm)
}
