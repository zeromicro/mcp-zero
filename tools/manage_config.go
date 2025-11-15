package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"gopkg.in/yaml.v3"

	"github.com/zeromicro/mcp-zero/internal/responses"
	"github.com/zeromicro/mcp-zero/internal/templates"
	"github.com/zeromicro/mcp-zero/internal/validation"
)

type ValidateConfigParams struct {
	ConfigPath  string `json:"config_path"`
	ServiceType string `json:"service_type,omitempty"` // "api" or "rpc"
}

type GenerateConfigParams struct {
	ServiceName string `json:"service_name"`
	ServiceType string `json:"service_type"` // "api" or "rpc"
	Environment string `json:"environment"`  // "development", "production", "test"
	Port        int    `json:"port,omitempty"`
	OutputPath  string `json:"output_path,omitempty"`
}

// ValidateConfig validates a go-zero configuration file
func ValidateConfig(ctx context.Context, req *mcp.CallToolRequest, params ValidateConfigParams) (*mcp.CallToolResult, any, error) {
	if params.ConfigPath == "" {
		return responses.FormatValidationError("config_path", "", "config_path is required", "Provide path to config file")
	}

	if !filepath.IsAbs(params.ConfigPath) {
		absPath, err := filepath.Abs(params.ConfigPath)
		if err != nil {
			return responses.FormatError(fmt.Sprintf("failed to resolve config path: %v", err))
		}
		params.ConfigPath = absPath
	}

	content, err := os.ReadFile(params.ConfigPath)
	if err != nil {
		return responses.FormatError(fmt.Sprintf("failed to read config file: %v", err))
	}

	// Parse config file
	var config map[string]interface{}
	ext := strings.ToLower(filepath.Ext(params.ConfigPath))

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(content, &config); err != nil {
			return responses.FormatError(fmt.Sprintf("failed to parse YAML config: %v", err))
		}
	case ".json":
		if err := json.Unmarshal(content, &config); err != nil {
			return responses.FormatError(fmt.Sprintf("failed to parse JSON config: %v", err))
		}
	default:
		return responses.FormatError(fmt.Sprintf("unsupported config file format: %s (use .yaml, .yml, or .json)", ext))
	}

	// Determine service type
	serviceType := params.ServiceType
	if serviceType == "" {
		// Auto-detect based on config fields
		if _, hasPort := config["Port"]; hasPort {
			serviceType = "api"
		} else if _, hasListenOn := config["ListenOn"]; hasListenOn {
			serviceType = "rpc"
		} else {
			serviceType = "api" // Default
		}
	}

	// Validate config
	var result *validation.ConfigValidationResult
	switch serviceType {
	case "api":
		result = validation.ValidateAPIConfig(config)
	case "rpc":
		result = validation.ValidateRPCConfig(config)
	default:
		return responses.FormatError(fmt.Sprintf("unsupported service type: %s (use 'api' or 'rpc')", serviceType))
	}

	// Format validation results
	var message strings.Builder
	message.WriteString(fmt.Sprintf("Configuration Validation: %s\n\n", params.ConfigPath))
	message.WriteString(fmt.Sprintf("Service Type: %s\n", serviceType))
	message.WriteString(fmt.Sprintf("Valid: %v\n\n", result.Valid))

	if len(result.Errors) > 0 {
		message.WriteString("=== Errors ===\n")
		for _, err := range result.Errors {
			message.WriteString(fmt.Sprintf("  ❌ %s: %s\n", err.Field, err.Message))
			if err.Value != nil {
				message.WriteString(fmt.Sprintf("     Current value: %v\n", err.Value))
			}
		}
		message.WriteString("\n")
	}

	if len(result.Warnings) > 0 {
		message.WriteString("=== Warnings ===\n")
		for _, warn := range result.Warnings {
			message.WriteString(fmt.Sprintf("  ⚠️  %s: %s\n", warn.Field, warn.Message))
			if warn.Suggestion != "" {
				message.WriteString(fmt.Sprintf("     Suggestion: %s\n", warn.Suggestion))
			}
		}
		message.WriteString("\n")
	}

	if result.Valid && len(result.Warnings) == 0 {
		message.WriteString("✅ Configuration is valid with no warnings!\n")
	} else if result.Valid {
		message.WriteString("✅ Configuration is valid but has warnings to address.\n")
	} else {
		message.WriteString("❌ Configuration has errors that must be fixed.\n")
	}

	data := map[string]any{
		"config_path":   params.ConfigPath,
		"service_type":  serviceType,
		"valid":         result.Valid,
		"error_count":   len(result.Errors),
		"warning_count": len(result.Warnings),
	}

	if !result.Valid {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{
				Text: message.String(),
			}},
			IsError: true,
		}, data, nil
	}

	return responses.FormatSuccessWithData(message.String(), data)
}

// GenerateConfigTemplate generates a configuration file template
func GenerateConfigTemplate(ctx context.Context, req *mcp.CallToolRequest, params GenerateConfigParams) (*mcp.CallToolResult, any, error) {
	if params.ServiceName == "" {
		return responses.FormatValidationError("service_name", "", "service_name is required", "Provide service name")
	}

	if params.ServiceType == "" {
		return responses.FormatValidationError("service_type", "", "service_type is required", "Use 'api' or 'rpc'")
	}

	if params.ServiceType != "api" && params.ServiceType != "rpc" {
		return responses.FormatValidationError("service_type", params.ServiceType, "invalid service type", "Use 'api' or 'rpc'")
	}

	if params.Environment == "" {
		params.Environment = "development"
	}

	// Set default port if not provided
	if params.Port == 0 {
		if params.ServiceType == "api" {
			params.Port = 8888
		} else {
			params.Port = 9090
		}
	}

	// Get template
	templateStr := templates.GetConfigTemplate(params.ServiceType, params.Environment)
	if templateStr == "" {
		return responses.FormatError(fmt.Sprintf("no template found for service type '%s' and environment '%s'", params.ServiceType, params.Environment))
	}

	// Parse and execute template
	tmpl, err := template.New("config").Parse(templateStr)
	if err != nil {
		return responses.FormatError(fmt.Sprintf("failed to parse template: %v", err))
	}

	metricsPort := params.Port + 1000
	data := map[string]interface{}{
		"ServiceName": params.ServiceName,
		"Port":        params.Port,
		"MetricsPort": metricsPort,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to execute template: %v", err))
	}

	configContent := buf.String()

	// Determine output path
	outputPath := params.OutputPath
	if outputPath == "" {
		envSuffix := ""
		if params.Environment != "development" {
			envSuffix = "-" + params.Environment
		}
		outputPath = fmt.Sprintf("etc/%s%s.yaml", params.ServiceName, envSuffix)
	}

	if !filepath.IsAbs(outputPath) {
		cwd, _ := os.Getwd()
		outputPath = filepath.Join(cwd, outputPath)
	}

	// Create directory if needed
	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to create directory: %v", err))
	}

	// Write config file
	if err := os.WriteFile(outputPath, []byte(configContent), 0644); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to write config file: %v", err))
	}

	message := fmt.Sprintf("Successfully generated %s configuration for %s environment\n\n", params.ServiceType, params.Environment)
	message += fmt.Sprintf("Output file: %s\n\n", outputPath)
	message += "Configuration includes:\n"

	switch params.Environment {
	case "production", "prod":
		message += "  ✓ Production-ready settings\n"
		message += "  ✓ File-based logging with rotation\n"
		message += "  ✓ Prometheus metrics endpoint\n"
		message += "  ✓ Service discovery configuration\n"
		message += "  ✓ Error-level logging\n"
	case "test":
		message += "  ✓ Test environment settings\n"
		message += "  ✓ Debug logging\n"
		message += "  ✓ Local host binding\n"
	default:
		message += "  ✓ Development environment settings\n"
		message += "  ✓ Console logging\n"
		message += "  ✓ Info-level logging\n"
	}

	message += "\nNext steps:\n"
	message += "  1. Review and customize the generated config\n"
	message += "  2. Use validate_config to verify your changes\n"
	message += fmt.Sprintf("  3. Start your service with: ./%s -f %s\n", params.ServiceName, outputPath)

	resultData := map[string]any{
		"service_name": params.ServiceName,
		"service_type": params.ServiceType,
		"environment":  params.Environment,
		"output_path":  outputPath,
		"port":         params.Port,
	}

	return responses.FormatSuccessWithData(message, resultData)
}
