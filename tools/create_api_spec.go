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

	"github.com/zeromicro/mcp-zero/internal/analyzer"
	"github.com/zeromicro/mcp-zero/internal/responses"
	"github.com/zeromicro/mcp-zero/internal/templates"
	"github.com/zeromicro/mcp-zero/internal/validation"
)

type CreateAPISpecParams struct {
	ServiceName   string `json:"service_name"`
	EndpointsJSON string `json:"endpoints_json"`
	OutputPath    string `json:"output_path,omitempty"`
}

type EndpointInput struct {
	Method   string                 `json:"method"`
	Path     string                 `json:"path"`
	Handler  string                 `json:"handler"`
	Request  map[string]interface{} `json:"request,omitempty"`
	Response map[string]interface{} `json:"response,omitempty"`
}

func CreateAPISpec(ctx context.Context, req *mcp.CallToolRequest, params CreateAPISpecParams) (*mcp.CallToolResult, any, error) {
	if err := validation.ValidateServiceName(params.ServiceName); err != nil {
		return responses.FormatValidationError("service_name", params.ServiceName, err.Error(), "Use lowercase letters, numbers, and hyphens only")
	}

	if params.EndpointsJSON == "" {
		return responses.FormatValidationError("endpoints_json", "", "endpoints_json is required", "Provide JSON array of endpoint definitions")
	}

	var endpointInputs []EndpointInput
	if err := json.Unmarshal([]byte(params.EndpointsJSON), &endpointInputs); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to parse endpoints_json: %v", err))
	}

	if len(endpointInputs) == 0 {
		return responses.FormatValidationError("endpoints_json", params.EndpointsJSON, "no endpoints provided", "Provide at least one endpoint definition")
	}

	outputPath := params.OutputPath
	if outputPath == "" {
		outputPath = fmt.Sprintf("%s.api", params.ServiceName)
	}
	if !filepath.IsAbs(outputPath) {
		cwd, _ := os.Getwd()
		outputPath = filepath.Join(cwd, outputPath)
	}

	spec := templates.APISpec{
		ServiceName: params.ServiceName,
		Types:       []templates.TypeDef{},
		Endpoints:   []templates.EndpointDef{},
	}

	typeMap := make(map[string]templates.TypeDef)

	for _, ep := range endpointInputs {
		method := strings.ToLower(ep.Method)
		handler := ep.Handler
		if handler == "" {
			handler = fmt.Sprintf("%s%sHandler", strings.Title(strings.ToLower(method)), toCamelCase(ep.Path))
		}

		endpoint := templates.EndpointDef{
			Handler: handler,
			Method:  method,
			Path:    ep.Path,
		}

		if ep.Request != nil {
			reqTypeName := handler + "Request"
			if _, exists := typeMap[reqTypeName]; !exists {
				typeMap[reqTypeName] = createTypeFromMap(reqTypeName, ep.Request)
			}
			endpoint.Request = reqTypeName
		}

		if ep.Response != nil {
			respTypeName := handler + "Response"
			if _, exists := typeMap[respTypeName]; !exists {
				typeMap[respTypeName] = createTypeFromMap(respTypeName, ep.Response)
			}
			endpoint.Response = respTypeName
		}

		spec.Endpoints = append(spec.Endpoints, endpoint)
	}

	for _, typeDef := range typeMap {
		spec.Types = append(spec.Types, typeDef)
	}

	tmpl, err := template.New("api").Parse(templates.APISpecTemplate)
	if err != nil {
		return responses.FormatError(fmt.Sprintf("failed to parse template: %v", err))
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, spec); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to generate spec: %v", err))
	}

	if err := os.WriteFile(outputPath, buf.Bytes(), 0644); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to write spec file: %v", err))
	}

	if _, err := analyzer.ParseAPISpecification(outputPath); err != nil {
		os.Remove(outputPath)
		return responses.FormatError(fmt.Sprintf("generated spec is invalid: %v", err))
	}

	message := fmt.Sprintf("Successfully created API specification: %s\n\nOutput file: %s\n", params.ServiceName, outputPath)
	message += fmt.Sprintf("\nEndpoints: %d\n", len(spec.Endpoints))
	message += fmt.Sprintf("Types: %d\n", len(spec.Types))
	message += "\nNext steps:\n"
	message += "  1. Review the generated specification\n"
	message += "  2. Use generate_api_from_spec to generate code\n"
	message += fmt.Sprintf("  3. goctl api go -api %s -dir ./output\n", outputPath)

	data := map[string]any{
		"service_name":   params.ServiceName,
		"output_path":    outputPath,
		"endpoint_count": len(spec.Endpoints),
		"type_count":     len(spec.Types),
	}

	return responses.FormatSuccessWithData(message, data)
}

func createTypeFromMap(name string, fields map[string]interface{}) templates.TypeDef {
	typeDef := templates.TypeDef{
		Name:   name,
		Fields: []templates.FieldDef{},
	}

	for fieldName, fieldValue := range fields {
		goType := "string"
		switch fieldValue.(type) {
		case float64:
			goType = "int64"
		case bool:
			goType = "bool"
		case []interface{}:
			goType = "[]string"
		case map[string]interface{}:
			goType = "map[string]string"
		}

		typeDef.Fields = append(typeDef.Fields, templates.FieldDef{
			Name:    strings.Title(fieldName),
			Type:    goType,
			JsonTag: fieldName,
		})
	}

	return typeDef
}

func toCamelCase(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	result := ""
	for _, part := range parts {
		if part != "" {
			result += strings.Title(part)
		}
	}
	return result
}
