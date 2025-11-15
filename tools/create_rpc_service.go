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

type CreateRPCServiceParams struct {
	ServiceName  string `json:"service_name"`
	ProtoContent string `json:"proto_content"`
	OutputDir    string `json:"output_dir,omitempty"`
	Style        string `json:"style,omitempty"`
}

func CreateRPCService(ctx context.Context, req *mcp.CallToolRequest, params CreateRPCServiceParams) (*mcp.CallToolResult, any, error) {
	if err := validation.ValidateServiceName(params.ServiceName); err != nil {
		return responses.FormatValidationError("service_name", params.ServiceName, err.Error(), "Use lowercase letters, numbers, and hyphens only")
	}

	style := params.Style
	if style == "" {
		style = "go_zero"
	}

	outputDir := params.OutputDir
	if outputDir == "" {
		outputDir = "."
	}
	if err := validation.ValidateOutputDir(outputDir); err != nil {
		return responses.FormatValidationError("output_dir", outputDir, err.Error(), "Provide an absolute path to an existing writable directory")
	}

	serviceDir := filepath.Join(outputDir, params.ServiceName)

	// Write proto file to output directory (not temp) so protoc can find it
	protoFile := filepath.Join(outputDir, params.ServiceName+".proto")
	if err := os.WriteFile(protoFile, []byte(params.ProtoContent), 0644); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to write proto file: %v", err))
	}
	defer os.Remove(protoFile)

	spec, err := analyzer.ParseProtoSpecification(protoFile)
	if err != nil {
		return responses.FormatError(fmt.Sprintf("failed to parse proto specification: %v", err))
	}

	executor, err := goctl.NewExecutor()
	if err != nil {
		return responses.FormatError(fmt.Sprintf("failed to create executor: %v", err))
	}

	// Use relative path for proto file and execute in outputDir
	args := []string{
		"rpc",
		"protoc",
		params.ServiceName + ".proto",
		"--go_out=.",
		"--go-grpc_out=.",
		"--zrpc_out=.",
		"--style", style,
	}

	result := executor.ExecuteInDir(outputDir, args...)
	if result.Error != nil {
		return responses.FormatError(fmt.Sprintf("failed to create RPC service: %v\nStderr: %s", result.Error, result.Stderr))
	}

	// Use a proper module path format (avoid module names starting with numbers)
	moduleName := "github.com/example/" + params.ServiceName

	if err := fixer.FixImports(serviceDir, moduleName); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to fix imports: %v", err))
	}

	if err := fixer.InitializeGoModule(serviceDir, moduleName); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to initialize Go module: %v", err))
	}

	if err := fixer.TidyGoModule(serviceDir); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to tidy Go module: %v", err))
	}

	if err := fixer.VerifyBuild(serviceDir); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to verify build: %v", err))
	}

	validator := goctl.NewValidator()
	if err := validator.ValidateServiceProject(serviceDir, "rpc"); err != nil {
		return responses.FormatError(fmt.Sprintf("project structure validation failed: %v", err))
	}

	message := fmt.Sprintf("Successfully created RPC service '%s'\n\nOutput directory: %s\n", params.ServiceName, serviceDir)
	message += fmt.Sprintf("\nService: %s\n", spec.ServiceName)
	message += "\nMethods:\n"
	for _, method := range spec.Methods {
		streamInfo := ""
		if method.Stream != "" {
			streamInfo = fmt.Sprintf(" [%s stream]", method.Stream)
		}
		message += fmt.Sprintf("  %s(%s) returns (%s)%s\n", method.Name, method.Request, method.Response, streamInfo)
	}
	message += fmt.Sprintf("\nMessages: %d\n", len(spec.Messages))
	message += "\nNext steps:\n"
	message += fmt.Sprintf("  1. cd %s\n", serviceDir)
	message += "  2. go mod tidy\n"
	message += "  3. go run .\n"

	data := map[string]any{
		"service_type":  "rpc",
		"service_name":  params.ServiceName,
		"output_dir":    serviceDir,
		"style":         style,
		"method_count":  len(spec.Methods),
		"message_count": len(spec.Messages),
	}

	return responses.FormatSuccessWithData(message, data)
}
