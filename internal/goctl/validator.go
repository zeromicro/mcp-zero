package goctl

import (
	"fmt"
	"os"
	"path/filepath"
)

// Validator validates generated code and project structure
type Validator struct{}

// NewValidator creates a new validator instance
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateServiceProject validates a generated service project
// Checks for required files and structure
func (v *Validator) ValidateServiceProject(projectPath string, serviceType string) error {
	// Check project directory exists
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		return fmt.Errorf("project directory does not exist: %s", projectPath)
	}

	// Check for go.mod
	goModPath := filepath.Join(projectPath, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return fmt.Errorf("go.mod not found - project not initialized")
	}

	// Type-specific validation
	switch serviceType {
	case "api":
		return v.validateAPIService(projectPath)
	case "rpc":
		return v.validateRPCService(projectPath)
	default:
		return fmt.Errorf("unknown service type: %s", serviceType)
	}
}

// validateAPIService validates API service structure
func (v *Validator) validateAPIService(projectPath string) error {
	// Check for .api file
	apiFiles, err := filepath.Glob(filepath.Join(projectPath, "*.api"))
	if err != nil {
		return fmt.Errorf("failed to search for .api files: %w", err)
	}
	if len(apiFiles) == 0 {
		return fmt.Errorf("no .api specification file found")
	}

	// Check for main entry point (should have a .go file)
	goFiles, err := filepath.Glob(filepath.Join(projectPath, "*.go"))
	if err != nil {
		return fmt.Errorf("failed to search for .go files: %w", err)
	}
	if len(goFiles) == 0 {
		return fmt.Errorf("no main .go file found")
	}

	// Check for etc directory (config)
	etcDir := filepath.Join(projectPath, "etc")
	if _, err := os.Stat(etcDir); os.IsNotExist(err) {
		return fmt.Errorf("etc directory not found - incomplete generation")
	}

	// Check for internal directory
	internalDir := filepath.Join(projectPath, "internal")
	if _, err := os.Stat(internalDir); os.IsNotExist(err) {
		return fmt.Errorf("internal directory not found - incomplete generation")
	}

	return nil
}

// validateRPCService validates RPC service structure
func (v *Validator) validateRPCService(projectPath string) error {
	// Check for .proto file
	protoFiles, err := filepath.Glob(filepath.Join(projectPath, "*.proto"))
	if err != nil {
		return fmt.Errorf("failed to search for .proto files: %w", err)
	}
	if len(protoFiles) == 0 {
		return fmt.Errorf("no .proto specification file found")
	}

	// Check for main entry point
	goFiles, err := filepath.Glob(filepath.Join(projectPath, "*.go"))
	if err != nil {
		return fmt.Errorf("failed to search for .go files: %w", err)
	}
	if len(goFiles) == 0 {
		return fmt.Errorf("no main .go file found")
	}

	// Check for etc directory
	etcDir := filepath.Join(projectPath, "etc")
	if _, err := os.Stat(etcDir); os.IsNotExist(err) {
		return fmt.Errorf("etc directory not found - incomplete generation")
	}

	// Check for internal directory
	internalDir := filepath.Join(projectPath, "internal")
	if _, err := os.Stat(internalDir); os.IsNotExist(err) {
		return fmt.Errorf("internal directory not found - incomplete generation")
	}

	return nil
}

// FileExists checks if a file exists
func (v *Validator) FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DirectoryExists checks if a directory exists
func (v *Validator) DirectoryExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
