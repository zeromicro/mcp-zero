package integration

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zeromicro/mcp-zero/tools"
)

func TestGenerateTemplate(t *testing.T) {
	ctx := context.Background()
	req := &mcp.CallToolRequest{}

	tests := []struct {
		name          string
		params        tools.GenerateTemplateParams
		expectSuccess bool
		checkContent  []string // Strings that should appear in generated content
		checkFile     bool     // Whether to verify file was created
	}{
		{
			name: "generate auth middleware",
			params: tools.GenerateTemplateParams{
				TemplateType: "middleware",
				TemplateName: "auth",
				Parameters:   `{"MiddlewareName": "Auth", "SecretKey": "test-secret-key"}`,
				OutputPath:   "middleware/auth_middleware.go",
			},
			expectSuccess: true,
			checkContent:  []string{"AuthMiddleware", "test-secret-key", "Bearer"},
			checkFile:     true,
		},
		{
			name: "generate logging middleware",
			params: tools.GenerateTemplateParams{
				TemplateType: "middleware",
				TemplateName: "logging",
				Parameters:   `{"MiddlewareName": "Logging"}`,
				OutputPath:   "middleware/logging_middleware.go",
			},
			expectSuccess: true,
			checkContent:  []string{"LoggingMiddleware", "logx.Infof", "time.Since"},
			checkFile:     true,
		},
		{
			name: "generate rate-limiting middleware",
			params: tools.GenerateTemplateParams{
				TemplateType: "middleware",
				TemplateName: "rate-limiting",
				Parameters: `{
					"MiddlewareName": "RateLimit",
					"RequestsPerPeriod": 100,
					"PeriodSeconds": 60
				}`,
				OutputPath: "middleware/ratelimit_middleware.go",
			},
			expectSuccess: true,
			checkContent:  []string{"RateLimitMiddleware", "redis", "100", "60"},
			checkFile:     true,
		},
		{
			name: "generate basic error handler",
			params: tools.GenerateTemplateParams{
				TemplateType: "error_handler",
				TemplateName: "basic",
				OutputPath:   "handler/error_handler.go",
			},
			expectSuccess: true,
			checkContent:  []string{"ErrorHandler", "ErrorResponse", "httpx"},
			checkFile:     true,
		},
		{
			name: "generate detailed error handler",
			params: tools.GenerateTemplateParams{
				TemplateType: "error_handler",
				TemplateName: "detailed",
				OutputPath:   "handler/detailed_error_handler.go",
			},
			expectSuccess: true,
			checkContent:  []string{"ErrorHandler", "ValidationError", "BusinessError"},
			checkFile:     true,
		},
		{
			name: "generate Dockerfile",
			params: tools.GenerateTemplateParams{
				TemplateType: "deployment",
				TemplateName: "docker",
				Parameters:   `{"ServiceName": "user-api", "Port": 8888}`,
				OutputPath:   "Dockerfile",
			},
			expectSuccess: true,
			checkContent:  []string{"FROM golang", "user-api", "8888", "WORKDIR"},
			checkFile:     true,
		},
		{
			name: "generate Kubernetes manifest",
			params: tools.GenerateTemplateParams{
				TemplateType: "deployment",
				TemplateName: "kubernetes",
				Parameters: `{
					"ServiceName": "user-api",
					"Port": 8888,
					"Replicas": 3
				}`,
				OutputPath: "k8s/user-api.yaml",
			},
			expectSuccess: true,
			checkContent:  []string{"kind: Service", "kind: Deployment", "user-api", "replicas: 3"},
			checkFile:     true,
		},
		{
			name: "generate systemd service",
			params: tools.GenerateTemplateParams{
				TemplateType: "deployment",
				TemplateName: "systemd",
				Parameters: `{
					"ServiceName": "user-api",
					"ServiceDescription": "User API Service",
					"ExecStart": "/usr/local/bin/user-api"
				}`,
				OutputPath: "deploy/user-api.service",
			},
			expectSuccess: true,
			checkContent:  []string{"[Unit]", "[Service]", "user-api", "Restart=on-failure"},
			checkFile:     true,
		},
		{
			name: "missing template_type",
			params: tools.GenerateTemplateParams{
				TemplateName: "auth",
			},
			expectSuccess: false,
		},
		{
			name: "missing template_name",
			params: tools.GenerateTemplateParams{
				TemplateType: "middleware",
			},
			expectSuccess: false,
		},
		{
			name: "invalid template type",
			params: tools.GenerateTemplateParams{
				TemplateType: "invalid_type",
				TemplateName: "something",
			},
			expectSuccess: false,
		},
		{
			name: "invalid template name",
			params: tools.GenerateTemplateParams{
				TemplateType: "middleware",
				TemplateName: "nonexistent",
			},
			expectSuccess: false,
		},
		{
			name: "invalid parameters JSON",
			params: tools.GenerateTemplateParams{
				TemplateType: "middleware",
				TemplateName: "auth",
				Parameters:   `{invalid json}`,
			},
			expectSuccess: false,
		},
	}

	// Create temporary directory for test outputs
	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	defer os.Chdir(originalWd)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _, err := tools.GenerateTemplate(ctx, req, tt.params)

			if tt.expectSuccess {
				if err != nil {
					t.Errorf("Expected success, got error: %v", err)
					return
				}
				if result.IsError {
					t.Errorf("Expected IsError=false, got IsError=true with content: %v", result.Content)
					return
				}

				// Check file was created
				if tt.checkFile && tt.params.OutputPath != "" {
					outputPath := tt.params.OutputPath
					if !filepath.IsAbs(outputPath) {
						outputPath = filepath.Join(tempDir, outputPath)
					}

					if _, err := os.Stat(outputPath); os.IsNotExist(err) {
						t.Errorf("Expected file to be created at %s, but it doesn't exist", outputPath)
						return
					}

					// Read and check content
					content, err := os.ReadFile(outputPath)
					if err != nil {
						t.Errorf("Failed to read generated file: %v", err)
						return
					}

					for _, expected := range tt.checkContent {
						if !strings.Contains(string(content), expected) {
							t.Errorf("Expected content to contain '%s', but it doesn't. Content:\n%s",
								expected, string(content))
						}
					}
				}

				// Verify response has content
				if len(result.Content) == 0 {
					t.Error("Expected response content")
				}
			} else {
				if err == nil && !result.IsError {
					t.Errorf("Expected error or IsError=true, got success")
				}
			}
		})
	}
}

func TestGenerateTemplateDefaultPaths(t *testing.T) {
	ctx := context.Background()
	req := &mcp.CallToolRequest{}

	tests := []struct {
		name         string
		params       tools.GenerateTemplateParams
		expectedPath string // Path should contain this
	}{
		{
			name: "auth middleware default path",
			params: tools.GenerateTemplateParams{
				TemplateType: "middleware",
				TemplateName: "auth",
				Parameters:   `{"MiddlewareName": "Auth", "SecretKey": "test"}`,
			},
			expectedPath: "middleware/auth_middleware.go",
		},
		{
			name: "docker default path",
			params: tools.GenerateTemplateParams{
				TemplateType: "deployment",
				TemplateName: "docker",
				Parameters:   `{"ServiceName": "test-api"}`,
			},
			expectedPath: "Dockerfile",
		},
		{
			name: "kubernetes default path",
			params: tools.GenerateTemplateParams{
				TemplateType: "deployment",
				TemplateName: "kubernetes",
				Parameters:   `{"ServiceName": "test-api", "Port": 8888}`,
			},
			expectedPath: "k8s/test-api.yaml",
		},
	}

	tempDir := t.TempDir()
	originalWd, _ := os.Getwd()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}
	defer os.Chdir(originalWd)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := tools.GenerateTemplate(ctx, req, tt.params)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Just verify file was created with expected path
			expectedFullPath := filepath.Join(tempDir, tt.expectedPath)
			if _, err := os.Stat(expectedFullPath); os.IsNotExist(err) {
				t.Errorf("Expected file to exist at %s", expectedFullPath)
			}
		})
	}
}

func TestListAvailableTemplates(t *testing.T) {
	ctx := context.Background()
	req := &mcp.CallToolRequest{}

	tests := []struct {
		name             string
		templateType     string
		expectedInList   []string
		expectValidation bool
	}{
		{
			name:             "list middleware templates",
			templateType:     "middleware",
			expectedInList:   []string{"auth", "logging", "rate-limiting"},
			expectValidation: true,
		},
		{
			name:             "list error_handler templates",
			templateType:     "error_handler",
			expectedInList:   []string{"basic", "detailed"},
			expectValidation: true,
		},
		{
			name:             "list deployment templates",
			templateType:     "deployment",
			expectedInList:   []string{"docker", "kubernetes", "systemd"},
			expectValidation: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := tools.GenerateTemplateParams{
				TemplateType: tt.templateType,
				// No template_name - should list available templates
			}

			result, _, _ := tools.GenerateTemplate(ctx, req, params)

			if tt.expectValidation {
				if !result.IsError {
					t.Error("Expected validation error for missing template_name")
				}

				// Just verify we got an error response with content
				if len(result.Content) == 0 {
					t.Error("Expected error content")
				}
			}
		})
	}
}
