package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/zeromicro/mcp-zero/internal/analyzer"
	"github.com/zeromicro/mcp-zero/internal/fixer"
	"github.com/zeromicro/mcp-zero/internal/goctl"
	"github.com/zeromicro/mcp-zero/internal/responses"
	"github.com/zeromicro/mcp-zero/internal/validation"
)

// GenerateAPIFromSpecParams defines the parameters for generate_api_from_spec tool (T040-T043)
type GenerateAPIFromSpecParams struct {
	APIFile   string `json:"api_file"`
	OutputDir string `json:"output_dir,omitempty"`
	Style     string `json:"style,omitempty"`
}

// GenerateAPIFromSpec generates go-zero API code from API specification file (T044-T046)
func GenerateAPIFromSpec(ctx context.Context, req *mcp.CallToolRequest, params GenerateAPIFromSpecParams) (*mcp.CallToolResult, any, error) {
	// T040: Validate API file exists
	apiFile := params.APIFile
	if !filepath.IsAbs(apiFile) {
		cwd, _ := os.Getwd()
		apiFile = filepath.Join(cwd, apiFile)
	}

	if _, err := os.Stat(apiFile); os.IsNotExist(err) {
		return responses.FormatError(fmt.Sprintf("API file not found: %s", apiFile))
	}

	// T041: Parse API specification to extract metadata
	spec, err := analyzer.ParseAPISpecification(apiFile)
	if err != nil {
		return responses.FormatError(fmt.Sprintf("failed to parse API specification: %v", err))
	}

	// T042: Validate output directory
	outputDir := params.OutputDir
	if outputDir == "" {
		outputDir, _ = os.Getwd()
	}
	if !filepath.IsAbs(outputDir) {
		cwd, _ := os.Getwd()
		outputDir = filepath.Join(cwd, outputDir)
	}

	if err := validation.ValidateOutputDir(outputDir); err != nil {
		return responses.FormatValidationError("output_dir", outputDir, err.Error(), "Provide an absolute path to an existing writable directory")
	}

	// T043: Set code style (default: go_zero)
	style := params.Style
	if style == "" {
		style = "go_zero"
	}
	if style != "go_zero" && style != "gozero" {
		return responses.FormatValidationError("style", style, "invalid style", "Use 'go_zero' or 'gozero'")
	}

	// T044: Execute goctl api go command
	executor, err := goctl.NewExecutor()
	if err != nil {
		return responses.FormatError(fmt.Sprintf("failed to create executor: %v", err))
	}

	args := []string{
		"api",
		"go",
		"-api", apiFile,
		"-dir", outputDir,
		"-style", style,
	}

	result := executor.Execute(args...)
	if result.Error != nil {
		return responses.FormatError(fmt.Sprintf("failed to generate API code: %v\nStderr: %s", result.Error, result.Stderr))
	}

	// Get module name from service name
	moduleName := spec.ServiceName

	// T045: Fix imports and initialize modules
	if err := fixer.FixImports(outputDir, moduleName); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to fix imports: %v", err))
	}

	if err := fixer.InitializeGoModule(outputDir, moduleName); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to initialize Go module: %v", err))
	}

	if err := fixer.TidyGoModule(outputDir); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to tidy Go module: %v", err))
	}

	// T046: Verify build success
	if err := fixer.VerifyBuild(outputDir); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to verify build: %v", err))
	}

	// Format success message with endpoint list
	message := fmt.Sprintf("Successfully generated go-zero API code from specification: %s\n\nOutput directory: %s\n", spec.ServiceName, outputDir)
	message += "\nEndpoints:\n"
	for _, ep := range spec.Endpoints {
		message += fmt.Sprintf("  %s %s → %s\n", ep.Method, ep.Path, ep.Handler)
	}
	message += fmt.Sprintf("\nTotal types: %d\n", len(spec.Types))
	message += "\nNext steps:\n"
	message += fmt.Sprintf("  1. cd %s\n", outputDir)
	message += "  2. go mod tidy\n"
	message += "  3. go run .\n"

	data := map[string]any{
		"service_name":   spec.ServiceName,
		"api_file":       apiFile,
		"output_dir":     outputDir,
		"style":          style,
		"endpoint_count": len(spec.Endpoints),
		"type_count":     len(spec.Types),
	}

	return responses.FormatSuccessWithData(message, data)
}
