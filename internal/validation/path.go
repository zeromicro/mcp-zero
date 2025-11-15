package validation

import (
	"fmt"
	"os"
	"path/filepath"
)

// ValidatePath validates a file/directory path
// Checks if path is absolute and writable
func ValidatePath(path string) error {
	// Check if path is absolute
	if !filepath.IsAbs(path) {
		return fmt.Errorf("path must be absolute, got: %s", path)
	}

	// Check if parent directory exists and is writable
	parentDir := filepath.Dir(path)
	if err := checkWritable(parentDir); err != nil {
		return fmt.Errorf("parent directory %s is not writable: %w", parentDir, err)
	}

	return nil
}

// ValidateOutputDir validates an output directory
// Creates it if it doesn't exist
func ValidateOutputDir(dir string) error {
	// Check if absolute
	if !filepath.IsAbs(dir) {
		return fmt.Errorf("output directory must be absolute path, got: %s", dir)
	}

	// Check if directory exists
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			// Try to create directory
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %w", err)
			}
			return nil
		}
		return fmt.Errorf("failed to check output directory: %w", err)
	}

	// Check if it's a directory
	if !info.IsDir() {
		return fmt.Errorf("output path exists but is not a directory: %s", dir)
	}

	// Check if writable
	return checkWritable(dir)
}

// checkWritable checks if a directory is writable by trying to create a temp file
func checkWritable(dir string) error {
	// Try to create a temp file
	tempFile, err := os.CreateTemp(dir, ".mcp-zero-test-*")
	if err != nil {
		return fmt.Errorf("directory is not writable")
	}

	// Clean up
	tempFile.Close()
	os.Remove(tempFile.Name())

	return nil
}

// EnsureDirectoryExists ensures a directory exists, creating it if necessary
func EnsureDirectoryExists(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(dir, 0755)
		}
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("path exists but is not a directory: %s", dir)
	}

	return nil
}
