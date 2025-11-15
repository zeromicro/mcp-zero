package integration

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zeromicro/mcp-zero/tools"
)

func TestCreateAPISpec(t *testing.T) {
	tests := []struct {
		name         string
		serviceName  string
		endpoints    []map[string]interface{}
		expectError  bool
		validateFile bool
	}{
		{
			name:        "simple REST API",
			serviceName: "userapi",
			endpoints: []map[string]interface{}{
				{
					"method":  "get",
					"path":    "/users/:id",
					"handler": "GetUserHandler",
				},
				{
					"method":  "post",
					"path":    "/users",
					"handler": "CreateUserHandler",
					"request": map[string]interface{}{
						"name":  "string",
						"email": "string",
					},
					"response": map[string]interface{}{
						"id":   0,
						"name": "string",
					},
				},
			},
			expectError:  false,
			validateFile: true,
		},
		{
			name:        "invalid service name with special char",
			serviceName: "user@api",
			endpoints: []map[string]interface{}{
				{
					"method":  "get",
					"path":    "/users",
					"handler": "ListUsersHandler",
				},
			},
			expectError:  true,
			validateFile: false,
		},
		{
			name:         "empty endpoints",
			serviceName:  "emptyapi",
			endpoints:    []map[string]interface{}{},
			expectError:  true,
			validateFile: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			outputPath := filepath.Join(tmpDir, tt.serviceName+".api")

			endpointsJSON, err := json.Marshal(tt.endpoints)
			if err != nil {
				t.Fatalf("Failed to marshal endpoints: %v", err)
			}

			params := tools.CreateAPISpecParams{
				ServiceName:   tt.serviceName,
				EndpointsJSON: string(endpointsJSON),
				OutputPath:    outputPath,
			}

			result, _, err := tools.CreateAPISpec(context.Background(), &mcp.CallToolRequest{}, params)

			if tt.expectError {
				if err == nil && !result.IsError {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result.IsError {
				t.Fatalf("Tool returned error: %v", result.Content)
			}

			if tt.validateFile {
				if _, err := os.Stat(outputPath); os.IsNotExist(err) {
					t.Errorf("Output file not created: %s", outputPath)
				}

				content, err := os.ReadFile(outputPath)
				if err != nil {
					t.Fatalf("Failed to read output file: %v", err)
				}

				contentStr := string(content)
				if contentStr == "" {
					t.Errorf("Output file is empty")
				}

				// Check for basic .api file structure
				mustContain := []string{"syntax", "info", "service"}
				for _, s := range mustContain {
					if !contains(contentStr, s) {
						t.Errorf("Output file missing required element: %s", s)
					}
				}
			}
		})
	}
}

func contains(content, substr string) bool {
	return len(content) > 0 && len(substr) > 0 && (content == substr || len(content) > len(substr) && findSubstring(content, substr))
}

func findSubstring(content, substr string) bool {
	for i := 0; i <= len(content)-len(substr); i++ {
		if content[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
