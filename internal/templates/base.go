package templates

import (
	"bytes"
	"fmt"
	"text/template"
)

// Template represents a code template
type Template struct {
	Name        string
	Type        string // "middleware", "error_handler", "deployment"
	Description string
	Content     string
	Parameters  []TemplateParameter
}

// TemplateParameter defines a parameter for template execution
type TemplateParameter struct {
	Name        string
	Type        string // "string", "int", "bool"
	Description string
	Required    bool
	Default     interface{}
}

// ExecuteTemplate executes a template with the given parameters
func ExecuteTemplate(tmpl *Template, params map[string]interface{}) (string, error) {
	// Validate required parameters
	for _, param := range tmpl.Parameters {
		if param.Required {
			if _, exists := params[param.Name]; !exists {
				if param.Default != nil {
					params[param.Name] = param.Default
				} else {
					return "", fmt.Errorf("required parameter '%s' is missing", param.Name)
				}
			}
		}
	}

	// Parse and execute template
	t, err := template.New(tmpl.Name).Parse(tmpl.Content)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, params); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// GetTemplate returns a template by name and type
func GetTemplate(templateType, templateName string) (*Template, error) {
	switch templateType {
	case "middleware":
		return GetMiddlewareTemplate(templateName)
	case "error_handler", "error-handler":
		return GetErrorHandlerTemplate(templateName)
	case "deployment":
		return GetDeploymentTemplate(templateName)
	default:
		return nil, fmt.Errorf("unknown template type: %s", templateType)
	}
}

// ListTemplates returns available templates by type
func ListTemplates(templateType string) []string {
	switch templateType {
	case "middleware":
		return []string{"auth", "logging", "rate-limiting"}
	case "error_handler", "error-handler":
		return []string{"basic", "detailed"}
	case "deployment":
		return []string{"docker", "kubernetes", "systemd"}
	default:
		return []string{}
	}
}
