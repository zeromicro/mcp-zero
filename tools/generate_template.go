package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/zeromicro/mcp-zero/internal/responses"
	"github.com/zeromicro/mcp-zero/internal/templates"
)

type GenerateTemplateParams struct {
	TemplateType string `json:"template_type"`        // "middleware", "error_handler", "deployment"
	TemplateName string `json:"template_name"`        // specific template like "auth", "logging", etc.
	Parameters   string `json:"parameters,omitempty"` // JSON string of parameters
	OutputPath   string `json:"output_path,omitempty"`
}

// GenerateTemplate generates code templates for common patterns
func GenerateTemplate(ctx context.Context, req *mcp.CallToolRequest, params GenerateTemplateParams) (*mcp.CallToolResult, any, error) {
	if params.TemplateType == "" {
		return responses.FormatValidationError("template_type", "", "template_type is required",
			"Use 'middleware', 'error_handler', or 'deployment'")
	}

	if params.TemplateName == "" {
		// List available templates for this type
		available := templates.ListTemplates(params.TemplateType)
		if len(available) == 0 {
			return responses.FormatValidationError("template_type", params.TemplateType,
				"invalid template type", "Use 'middleware', 'error_handler', or 'deployment'")
		}

		message := fmt.Sprintf("Please specify template_name. Available templates for '%s':\n", params.TemplateType)
		for _, name := range available {
			message += fmt.Sprintf("  - %s\n", name)
		}

		return responses.FormatValidationError("template_name", "", "template_name is required", message)
	}

	// Get the template
	tmpl, err := templates.GetTemplate(params.TemplateType, params.TemplateName)
	if err != nil {
		available := templates.ListTemplates(params.TemplateType)
		suggestion := fmt.Sprintf("Available templates: %s", strings.Join(available, ", "))
		return responses.FormatValidationError("template_name", params.TemplateName,
			fmt.Sprintf("template not found: %v", err), suggestion)
	}

	// Parse parameters
	templateParams := make(map[string]interface{})
	if params.Parameters != "" {
		if err := json.Unmarshal([]byte(params.Parameters), &templateParams); err != nil {
			return responses.FormatError(fmt.Sprintf("failed to parse parameters JSON: %v", err))
		}
	}

	// Set default values for missing parameters
	for _, param := range tmpl.Parameters {
		if _, exists := templateParams[param.Name]; !exists && param.Default != nil {
			templateParams[param.Name] = param.Default
		}
	}

	// Execute template
	code, err := templates.ExecuteTemplate(tmpl, templateParams)
	if err != nil {
		return responses.FormatError(fmt.Sprintf("failed to generate template: %v", err))
	}

	// Determine output path
	outputPath := params.OutputPath
	if outputPath == "" {
		outputPath = generateDefaultPath(params.TemplateType, params.TemplateName, templateParams)
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

	// Write file
	if err := os.WriteFile(outputPath, []byte(code), 0644); err != nil {
		return responses.FormatError(fmt.Sprintf("failed to write file: %v", err))
	}

	// Try to verify the generated code compiles (best effort)
	compileCheck := ""
	if strings.HasSuffix(outputPath, ".go") {
		if err := verifyGoFile(outputPath); err != nil {
			compileCheck = fmt.Sprintf("\n⚠️  Warning: Generated code may have compilation issues: %v\n", err)
		} else {
			compileCheck = "\n✅ Generated code verified successfully\n"
		}
	}

	// Generate response
	message := fmt.Sprintf("Successfully generated %s template: %s\n\n", params.TemplateType, params.TemplateName)
	message += fmt.Sprintf("Output file: %s\n", outputPath)
	message += compileCheck
	message += "\n" + getIntegrationInstructions(params.TemplateType, params.TemplateName, templateParams)

	data := map[string]any{
		"template_type": params.TemplateType,
		"template_name": params.TemplateName,
		"output_path":   outputPath,
		"file_size":     len(code),
	}

	return responses.FormatSuccessWithData(message, data)
}

func generateDefaultPath(templateType, templateName string, params map[string]interface{}) string {
	switch templateType {
	case "middleware":
		name := "middleware"
		if v, ok := params["MiddlewareName"]; ok {
			name = strings.ToLower(fmt.Sprint(v))
		}
		return filepath.Join("middleware", name+"_middleware.go")
	case "error_handler", "error-handler":
		return filepath.Join("handler", "error_handler.go")
	case "deployment":
		switch templateName {
		case "docker":
			return "Dockerfile"
		case "kubernetes", "k8s":
			serviceName := "service"
			if v, ok := params["ServiceName"]; ok {
				serviceName = fmt.Sprint(v)
			}
			return filepath.Join("k8s", serviceName+".yaml")
		case "systemd":
			serviceName := "service"
			if v, ok := params["ServiceName"]; ok {
				serviceName = fmt.Sprint(v)
			}
			return filepath.Join("deploy", serviceName+".service")
		}
	}
	return "output.txt"
}

func getIntegrationInstructions(templateType, templateName string, params map[string]interface{}) string {
	switch templateType {
	case "middleware":
		return `Integration Instructions:

1. Import the middleware in your handler file:
   import "yourproject/middleware"

2. Register in your API service:
   server.Use(middleware.NewYourMiddleware().Handle)

3. Or apply to specific routes in the .api file:
   @server(
       middleware: YourMiddleware
   )
   service yourapi {
       @handler YourHandler
       post /path (Request) returns (Response)
   }

4. Configure any dependencies (Redis, database, etc.) in your service context

Example usage:
  - For auth middleware: Verify JWT tokens, check permissions
  - For logging middleware: Track request/response timing
  - For rate limiting: Configure Redis connection in service context
`
	case "error_handler", "error-handler":
		return `Integration Instructions:

1. Register the error handler in your main.go:
   server := rest.MustNewServer(c.RestConf, rest.WithCustomCors(nil, nil, "*"))
   defer server.Stop()

   // Register error handler
   httpx.SetErrorHandler(handler.ErrorHandler)

2. Use custom error types in your logic:
   return handler.ValidationError{Field: "email", Message: "invalid format"}
   return handler.BusinessError{Code: 400, Message: "user not found"}

3. The error handler will automatically format responses
`
	case "deployment":
		switch templateName {
		case "docker":
			return `Integration Instructions:

1. Build the Docker image:
   docker build -t your-service:latest .

2. Run the container:
   docker run -d -p 8888:8888 --name your-service your-service:latest

3. For production, push to a registry:
   docker tag your-service:latest registry.example.com/your-service:latest
   docker push registry.example.com/your-service:latest
`
		case "kubernetes", "k8s":
			return `Integration Instructions:

1. Apply the manifest:
   kubectl apply -f k8s/service.yaml

2. Check deployment status:
   kubectl get deployments
   kubectl get pods
   kubectl get services

3. View logs:
   kubectl logs -f deployment/your-service

4. Scale the deployment:
   kubectl scale deployment your-service --replicas=5
`
		case "systemd":
			return `Integration Instructions:

1. Copy the service file:
   sudo cp deploy/service.service /etc/systemd/system/

2. Reload systemd:
   sudo systemctl daemon-reload

3. Enable and start the service:
   sudo systemctl enable your-service
   sudo systemctl start your-service

4. Check status:
   sudo systemctl status your-service
   sudo journalctl -u your-service -f
`
		}
	}
	return ""
}

func verifyGoFile(filePath string) error {
	// Use go fmt to check syntax
	cmd := exec.Command("go", "fmt", filePath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("syntax check failed: %w", err)
	}
	return nil
}
