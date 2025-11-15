package tools

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/zeromicro/mcp-zero/internal/fixer"
	"github.com/zeromicro/mcp-zero/internal/goctl"
	"github.com/zeromicro/mcp-zero/internal/responses"
	"github.com/zeromicro/mcp-zero/internal/security"
)

type GenerateModelParams struct {
	SourceType string `json:"source_type"`
	Source     string `json:"source"`
	Table      string `json:"table"`
	OutputDir  string `json:"output_dir,omitempty"`
	Style      string `json:"style,omitempty"`
}

func GenerateModel(ctx context.Context, req *mcp.CallToolRequest, params GenerateModelParams) (*mcp.CallToolResult, any, error) {
	if params.SourceType != "mysql" && params.SourceType != "postgresql" && params.SourceType != "mongo" {
		return responses.FormatValidationError("source_type", params.SourceType, "invalid source type", "Use 'mysql', 'postgresql', or 'mongo'")
	}

	if params.Source == "" {
		return responses.FormatValidationError("source", params.Source, "source is required", "Provide database connection string")
	}

	if params.Table == "" {
		return responses.FormatValidationError("table", params.Table, "table is required", "Provide table name")
	}

	outputDir := params.OutputDir
	if outputDir == "" {
		outputDir = "./model"
	}

	style := params.Style
	if style == "" {
		style = "go_zero"
	}

	var connInfo *security.ConnectionInfo
	var err error

	connInfo, err = security.ParseConnectionString(params.Source)
	if err != nil {
		return responses.FormatError(fmt.Sprintf("failed to parse connection string: %v", err))
	}
	defer connInfo.Clear()

	executor, err := goctl.NewExecutor()
	if err != nil {
		return responses.FormatError(fmt.Sprintf("failed to create executor: %v", err))
	}

	args := []string{
		"model",
		params.SourceType,
		"datasource",
		"-url", connInfo.ToDSN(),
		"-table", params.Table,
		"-dir", outputDir,
		"-style", style,
	}

	result := executor.Execute(args...)
	if result.Error != nil {
		connInfo.Clear()
		return responses.FormatError(fmt.Sprintf("failed to generate model: %v\nStderr: %s", result.Error, result.Stderr))
	}

	connInfo.Clear()

	moduleName := "model"
	if err := fixer.FixImports(outputDir, moduleName); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to fix imports: %v", err))
	}

	if err := fixer.InitializeGoModule(outputDir, moduleName); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to initialize Go module: %v", err))
	}

	if err := fixer.TidyGoModule(outputDir); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to tidy Go module: %v", err))
	}

	if err := fixer.VerifyBuild(outputDir); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to verify build: %v", err))
	}

	message := fmt.Sprintf("Successfully generated database model for table '%s'\n\nOutput directory: %s\n", params.Table, outputDir)
	message += fmt.Sprintf("\nSource Type: %s\n", params.SourceType)
	message += fmt.Sprintf("Table: %s\n", params.Table)
	message += "\nNext steps:\n"
	message += fmt.Sprintf("  1. cd %s\n", outputDir)
	message += "  2. Review generated model code\n"
	message += "  3. Integrate with your service\n"

	absPath, _ := filepath.Abs(outputDir)
	data := map[string]any{
		"source_type": params.SourceType,
		"table":       params.Table,
		"output_dir":  absPath,
		"style":       style,
	}

	return responses.FormatSuccessWithData(message, data)
}
