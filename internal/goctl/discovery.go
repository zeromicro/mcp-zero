package goctl

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// DiscoverGoctl finds the goctl executable using multiple fallback strategies
// Returns the absolute path to goctl or an error if not found
func DiscoverGoctl() (string, error) {
	// Strategy 1: Check GOCTL_PATH environment variable (highest priority)
	if goctlPath := os.Getenv("GOCTL_PATH"); goctlPath != "" {
		if isExecutable(goctlPath) {
			return goctlPath, nil
		}
	}

	// Strategy 2: Search common installation locations
	commonPaths := []string{
		"/usr/local/bin/goctl",
		filepath.Join(os.Getenv("HOME"), "go", "bin", "goctl"),
		filepath.Join(os.Getenv("HOME"), "Develop", "go", "bin", "goctl"),
	}

	// Add GOPATH/bin if GOPATH is set
	if goPath := os.Getenv("GOPATH"); goPath != "" {
		commonPaths = append(commonPaths, filepath.Join(goPath, "bin", "goctl"))
	}

	for _, path := range commonPaths {
		if isExecutable(path) {
			return path, nil
		}
	}

	// Strategy 3: Search in PATH
	if path, err := exec.LookPath("goctl"); err == nil {
		return path, nil
	}

	// Not found - return actionable error
	return "", fmt.Errorf("goctl not found. Install with:\n  go install github.com/zeromicro/go-zero/tools/goctl@latest\nOr set GOCTL_PATH environment variable")
}

// isExecutable checks if a file exists and is executable
func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	// Check if it's a regular file and executable
	mode := info.Mode()
	return !info.IsDir() && (mode&0111 != 0)
}
