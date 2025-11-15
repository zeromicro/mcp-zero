package tools

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/zeromicro/mcp-zero/internal/fixer"
	"github.com/zeromicro/mcp-zero/internal/goctl"
	"github.com/zeromicro/mcp-zero/internal/responses"
	"github.com/zeromicro/mcp-zero/internal/validation"
)

// CreateAPIServiceParams defines the parameters for creating an API service
type CreateAPIServiceParams struct {
	ServiceName string `json:"service_name"`
	Port        int    `json:"port,omitempty"`
	OutputDir   string `json:"output_dir,omitempty"`
	Style       string `json:"style,omitempty"`
}

// CreateAPIService creates a new go-zero API service
func CreateAPIService(ctx context.Context, req *mcp.CallToolRequest, params CreateAPIServiceParams) (*mcp.CallToolResult, any, error) {
	// Validate service name
	if err := validation.ValidateServiceName(params.ServiceName); err != nil {
		return responses.FormatValidationError("service_name", params.ServiceName, err.Error(), "Use lowercase letters, numbers, and hyphens only")
	}

	// Validate and set default port
	port := params.Port
	if port == 0 {
		port = 8888
	}
	if err := validation.ValidatePort(port); err != nil {
		return responses.FormatValidationError("port", fmt.Sprintf("%d", port), err.Error(), "Use a port number between 1024 and 65535")
	}

	// Set default style
	style := params.Style
	if style == "" {
		style = "go_zero"
	}

	// Validate output directory
	outputDir := params.OutputDir
	if outputDir == "" {
		outputDir = "."
	}
	if err := validation.ValidateOutputDir(outputDir); err != nil {
		return responses.FormatValidationError("output_dir", outputDir, err.Error(), "Provide an absolute path to an existing writable directory")
	}

	// Prepare service directory
	serviceDir := filepath.Join(outputDir, params.ServiceName)

	// Execute goctl api new command
	executor, err := goctl.NewExecutor()
	if err != nil {
		return responses.FormatError(fmt.Sprintf("failed to create executor: %v", err))
	}

	// goctl api new creates service in current directory, so we execute in outputDir
	args := []string{
		"api",
		"new",
		params.ServiceName,
		"--style", style,
	}

	result := executor.ExecuteInDir(outputDir, args...)
	if result.Error != nil {
		return responses.FormatError(fmt.Sprintf("failed to create API service: %v\nStderr: %s", result.Error, result.Stderr))
	}

	// Use a proper module path format (avoid module names starting with numbers)
	moduleName := "github.com/example/" + params.ServiceName

	// Fix imports
	if err := fixer.FixImports(serviceDir, moduleName); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to fix imports: %v", err))
	}

	// Initialize Go module
	if err := fixer.InitializeGoModule(serviceDir, moduleName); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to initialize Go module: %v", err))
	}

	// Tidy module
	if err := fixer.TidyGoModule(serviceDir); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to tidy Go module: %v", err))
	}

	// Update config file with port
	if err := fixer.UpdateConfigFile(serviceDir, params.ServiceName, port); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to update config file: %v", err))
	}

	// Verify build
	if err := fixer.VerifyBuild(serviceDir); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to verify build: %v", err))
	}

	// Validate project structure
	validator := goctl.NewValidator()
	if err := validator.ValidateServiceProject(serviceDir, "api"); err != nil {
		return responses.FormatError(fmt.Sprintf("project structure validation failed: %v", err))
	}

	// Return success response
	additionalInfo := map[string]string{
		"port":  fmt.Sprintf("%d", port),
		"style": style,
	}
	return responses.FormatServiceCreated("api", params.ServiceName, serviceDir, additionalInfo)
}
