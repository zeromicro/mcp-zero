package responses

import (
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func FormatSuccess(message string) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: message}},
	}, nil, nil
}

func FormatSuccessWithData(message string, data any) (*mcp.CallToolResult, any, error) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return FormatError(fmt.Sprintf("failed to serialize data: %v", err))
	}
	fullMessage := message
	if len(jsonData) > 0 {
		fullMessage = fmt.Sprintf("%s\n\n%s", message, string(jsonData))
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fullMessage}},
	}, data, nil
}

func FormatError(message string) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error: %s", message)}},
		IsError: true,
	}, nil, fmt.Errorf("%s", message)
}

func FormatValidationError(field, value, reason, suggestion string) (*mcp.CallToolResult, any, error) {
	message := fmt.Sprintf("Validation Error\n\nField: %s\nValue: %s\nReason: %s", field, value, reason)
	if suggestion != "" {
		message += fmt.Sprintf("\nSuggestion: %s", suggestion)
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: message}},
		IsError: true,
	}, nil, fmt.Errorf("validation failed for %s: %s", field, reason)
}

func FormatServiceCreated(serviceType, serviceName, outputDir string, additionalInfo map[string]string) (*mcp.CallToolResult, any, error) {
	message := fmt.Sprintf("Successfully created %s service '%s'\n\nOutput directory: %s\n", serviceType, serviceName, outputDir)
	if len(additionalInfo) > 0 {
		message += "\nAdditional Information:\n"
		for key, value := range additionalInfo {
			message += fmt.Sprintf("  %s: %s\n", key, value)
		}
	}
	message += "\nNext steps:\n"
	message += fmt.Sprintf("  1. cd %s\n", outputDir)
	message += "  2. go mod tidy\n"
	message += "  3. go run .\n"
	data := map[string]any{
		"service_type":    serviceType,
		"service_name":    serviceName,
		"output_dir":      outputDir,
		"additional_info": additionalInfo,
	}
	return FormatSuccessWithData(message, data)
}
