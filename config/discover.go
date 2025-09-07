package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func findFirstMatchingFile(filenames []string) (string, error) {
	startDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := startDir

	for {
		for _, filename := range filenames {
			path := filepath.Join(dir, filename)
			if fileExists(path) {
				return path, nil
			}
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	if home, err := os.UserHomeDir(); err == nil {
		for _, name := range filenames {
			candidate := filepath.Join(home, name)
			if fileExists(candidate) {
				return candidate, nil
			}
		}

		configDir := filepath.Join(home, ".config")
		for _, name := range filenames {
			candidate := filepath.Join(configDir, name)
			if fileExists(candidate) {
				return candidate, nil
			}
		}
	}

	return "", fmt.Errorf("no matching config file found in any parent directory")
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
