package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zeromicro/mcp-zero/tools"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      string
		serviceType string
		expectValid bool
	}{
		{
			name:        "valid API config",
			config:      "Name: testapi\nHost: 0.0.0.0\nPort: 8888\n",
			serviceType: "api",
			expectValid: true,
		},
		{
			name:        "invalid API config",
			config:      "Name: testapi\n",
			serviceType: "api",
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yaml")
			os.WriteFile(configPath, []byte(tt.config), 0644)

			params := tools.ValidateConfigParams{
				ConfigPath:  configPath,
				ServiceType: tt.serviceType,
			}

			result, _, _ := tools.ValidateConfig(context.Background(), &mcp.CallToolRequest{}, params)

			if tt.expectValid && result.IsError {
				t.Errorf("Expected valid config")
			}
			if !tt.expectValid && !result.IsError {
				t.Errorf("Expected invalid config")
			}
		})
	}
}

func TestGenerateConfigTemplate(t *testing.T) {
	tmpDir := t.TempDir()

	params := tools.GenerateConfigParams{
		ServiceName: "testapi",
		ServiceType: "api",
		Environment: "development",
		Port:        8888,
		OutputPath:  filepath.Join(tmpDir, "config.yaml"),
	}

	result, _, err := tools.GenerateConfigTemplate(context.Background(), &mcp.CallToolRequest{}, params)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result.IsError {
		t.Fatalf("Tool returned error")
	}

	if _, err := os.Stat(params.OutputPath); os.IsNotExist(err) {
		t.Errorf("Output file not created")
	}
}
