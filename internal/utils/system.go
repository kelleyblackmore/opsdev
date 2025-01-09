package utils

import (
	"os"
	"path/filepath"
)

// GetHomeDir returns the user's home directory
func GetHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return home, nil
}

// EnsureDirectory creates a directory if it doesn't exist
func EnsureDirectory(path string) error {
	return os.MkdirAll(path, 0755)
}

// IsExecutable checks if a file exists and is executable
func IsExecutable(name string) bool {
	path, err := exec.LookPath(name)
	if err != nil {
		return false
	}
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode()&0111 != 0
}

// GetTempDir returns a temporary directory for downloads
func GetTempDir() string {
	return filepath.Join(os.TempDir(), "opsdev")
}

// CleanupTempFiles removes temporary files after installation
func CleanupTempFiles() error {
	return os.RemoveAll(GetTempDir())
}