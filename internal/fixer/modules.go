package fixer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// InitializeGoModule initializes a Go module in the project directory
func InitializeGoModule(projectPath string, moduleName string) error {
	// Check if go.mod already exists
	goModPath := filepath.Join(projectPath, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		// go.mod already exists, skip initialization
		return nil
	}

	// Run go mod init
	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go mod init failed: %s\n%s", err, string(output))
	}

	// Run go mod tidy to resolve dependencies
	return TidyGoModule(projectPath)
}

// TidyGoModule runs go mod tidy to resolve dependencies
func TidyGoModule(projectPath string) error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go mod tidy failed: %s\n%s", err, string(output))
	}

	return nil
}

// VerifyBuild verifies the project builds successfully
func VerifyBuild(projectPath string) error {
	cmd := exec.Command("go", "build", "-o", "/dev/null", ".")
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build failed: %s\n%s", err, string(output))
	}

	return nil
}

// GetGoModuleName extracts module name from go.mod file
func GetGoModuleName(projectPath string) (string, error) {
	goModPath := filepath.Join(projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod: %w", err)
	}

	// Parse module name from first line: "module modulename"
	lines := string(content)
	if len(lines) > 7 && lines[:7] == "module " {
		endIdx := 7
		for endIdx < len(lines) && lines[endIdx] != '\n' && lines[endIdx] != '\r' {
			endIdx++
		}
		return lines[7:endIdx], nil
	}

	return "", fmt.Errorf("could not find module name in go.mod")
}
