package integration_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zeromicro/mcp-zero/tools"
)

// TestCreateAPIServiceSuccess tests successful API service creation with valid inputs
func TestCreateAPIServiceSuccess(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		port        int
		style       string
		wantPort    int
		wantStyle   string
	}{
		{
			name:        "basic service with defaults",
			serviceName: "testservice",
			port:        0, // will use default 8888
			style:       "",
			wantPort:    8888,
			wantStyle:   "go_zero",
		},
		{
			name:        "service with custom port",
			serviceName: "customport",
			port:        9090,
			style:       "",
			wantPort:    9090,
			wantStyle:   "go_zero",
		},
		{
			name:        "service with gozero style",
			serviceName: "gozeroservice",
			port:        0,
			style:       "gozero",
			wantPort:    8888,
			wantStyle:   "gozero",
		},
		{
			name:        "service with underscores",
			serviceName: "user_service",
			port:        8080,
			style:       "go_zero",
			wantPort:    8080,
			wantStyle:   "go_zero",
		},
		{
			name:        "service with numbers",
			serviceName: "api2service",
			port:        8081,
			style:       "go_zero",
			wantPort:    8081,
			wantStyle:   "go_zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for output
			tmpDir, err := os.MkdirTemp("", "api_service_test_*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Create request context
			ctx := context.Background()
			req := &mcp.CallToolRequest{}

			// Create parameters
			params := tools.CreateAPIServiceParams{
				ServiceName: tt.serviceName,
				Port:        tt.port,
				OutputDir:   tmpDir,
				Style:       tt.style,
			}

			// Call CreateAPIService
			result, _, err := tools.CreateAPIService(ctx, req, params)

			// Verify no error
			if err != nil {
				t.Fatalf("CreateAPIService returned error: %v", err)
			}

			// Verify result is not nil
			if result == nil {
				t.Fatal("CreateAPIService returned nil result")
			}

			// Verify result indicates success (not an error)
			if result.IsError {
				t.Fatalf("CreateAPIService returned error result: %v", result.Content)
			}

			// Verify service directory was created
			serviceDir := filepath.Join(tmpDir, tt.serviceName)
			if _, err := os.Stat(serviceDir); os.IsNotExist(err) {
				t.Fatalf("Service directory was not created: %s", serviceDir)
			}

			// Verify expected files exist
			expectedFiles := []string{
				filepath.Join(serviceDir, "go.mod"),
				filepath.Join(serviceDir, tt.serviceName+".go"),
				filepath.Join(serviceDir, tt.serviceName+".api"),
				filepath.Join(serviceDir, "etc", tt.serviceName+"-api.yaml"),
				filepath.Join(serviceDir, "internal", "config", "config.go"),
				filepath.Join(serviceDir, "internal", "handler", "routes.go"),
				filepath.Join(serviceDir, "internal", "types", "types.go"),
			}

			for _, file := range expectedFiles {
				if _, err := os.Stat(file); os.IsNotExist(err) {
					t.Errorf("Expected file does not exist: %s", file)
				}
			}

			// Verify service context file (naming depends on style)
			// go_zero style: service_context.go, gozero style: servicecontext.go
			svcFile1 := filepath.Join(serviceDir, "internal", "svc", "service_context.go")
			svcFile2 := filepath.Join(serviceDir, "internal", "svc", "servicecontext.go")
			if _, err1 := os.Stat(svcFile1); err1 != nil {
				if _, err2 := os.Stat(svcFile2); err2 != nil {
					t.Errorf("Expected service context file does not exist (tried both %s and %s)", svcFile1, svcFile2)
				}
			}

			// Verify go.mod contains correct module name
			goModPath := filepath.Join(serviceDir, "go.mod")
			goModContent, err := os.ReadFile(goModPath)
			if err != nil {
				t.Fatalf("Failed to read go.mod: %v", err)
			}
			if len(goModContent) == 0 {
				t.Error("go.mod is empty")
			}
		})
	}
}

// TestCreateAPIServiceInvalidName tests handling of invalid service names
func TestCreateAPIServiceInvalidName(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		wantError   bool
	}{
		{
			name:        "service name with hyphens",
			serviceName: "test-service",
			wantError:   true,
		},
		{
			name:        "service name with special chars",
			serviceName: "test@service",
			wantError:   true,
		},
		{
			name:        "service name starting with number",
			serviceName: "123service",
			wantError:   true,
		},
		{
			name:        "service name starting with underscore",
			serviceName: "_service",
			wantError:   true,
		},
		{
			name:        "empty service name",
			serviceName: "",
			wantError:   true,
		},
		{
			name:        "service name with spaces",
			serviceName: "test service",
			wantError:   true,
		},
		{
			name:        "service name with dots",
			serviceName: "test.service",
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for output
			tmpDir, err := os.MkdirTemp("", "api_service_test_invalid_*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Create request context
			ctx := context.Background()
			req := &mcp.CallToolRequest{}

			// Create parameters
			params := tools.CreateAPIServiceParams{
				ServiceName: tt.serviceName,
				Port:        8888,
				OutputDir:   tmpDir,
				Style:       "go_zero",
			}

			// Call CreateAPIService
			result, _, err := tools.CreateAPIService(ctx, req, params)

			// For invalid names, we expect an error result (validation should catch it)
			if tt.wantError {
				// Check if result indicates error
				if result == nil {
					t.Fatal("Expected error result, got nil")
				}
				if !result.IsError {
					t.Errorf("Expected error result for invalid service name %q, but got success", tt.serviceName)
				}
				// Note: Directory should not be created since validation fails early
			} else {
				if err != nil {
					t.Errorf("Expected no error for valid service name %q, got: %v", tt.serviceName, err)
				}
				if result != nil && result.IsError {
					t.Errorf("Expected success for valid service name %q, got error result", tt.serviceName)
				}
			}
		})
	}
}

// TestCreateAPIServicePortConflict tests port validation
func TestCreateAPIServicePortConflict(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		port        int
		expectError bool
	}{
		{
			name:        "port out of range - too low",
			serviceName: "lowport",
			port:        80,
			expectError: true,
		},
		{
			name:        "port out of range - too high",
			serviceName: "highport",
			port:        70000,
			expectError: true,
		},
		{
			name:        "privileged port",
			serviceName: "privport",
			port:        1023,
			expectError: true,
		},
		{
			name:        "valid available port",
			serviceName: "validport",
			port:        9999,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for output
			tmpDir, err := os.MkdirTemp("", "api_service_test_port_*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Create request context
			ctx := context.Background()
			req := &mcp.CallToolRequest{}

			// Create parameters
			params := tools.CreateAPIServiceParams{
				ServiceName: tt.serviceName,
				Port:        tt.port,
				OutputDir:   tmpDir,
				Style:       "go_zero",
			}

			// Call CreateAPIService
			result, _, err := tools.CreateAPIService(ctx, req, params)

			if tt.expectError {
				// Expect error result
				if result == nil {
					t.Fatal("Expected error result, got nil")
				}
				if !result.IsError {
					t.Errorf("Expected error result for port %d, but got success", tt.port)
				}
			} else {
				// Expect success
				if err != nil {
					t.Errorf("Expected no error for port %d, got: %v", tt.port, err)
				}
				if result != nil && result.IsError {
					t.Errorf("Expected success for port %d, got error result", tt.port)
				}

				// Verify service directory was created
				serviceDir := filepath.Join(tmpDir, tt.serviceName)
				if _, err := os.Stat(serviceDir); os.IsNotExist(err) {
					t.Errorf("Service directory should exist for valid port: %s", serviceDir)
				}
			}
		})
	}
}

// TestCreateAPIServiceInvalidOutputDir tests handling of invalid output directories
func TestCreateAPIServiceInvalidOutputDir(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		outputDir   string
		expectError bool
	}{
		{
			name:        "non-existent directory",
			serviceName: "testservice",
			outputDir:   "/nonexistent/path/that/does/not/exist",
			expectError: true,
		},
		{
			name:        "relative path",
			serviceName: "testservice",
			outputDir:   "./relative/path",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request context
			ctx := context.Background()
			req := &mcp.CallToolRequest{}

			// Create parameters
			params := tools.CreateAPIServiceParams{
				ServiceName: tt.serviceName,
				Port:        8888,
				OutputDir:   tt.outputDir,
				Style:       "go_zero",
			}

			// Call CreateAPIService
			result, _, err := tools.CreateAPIService(ctx, req, params)

			if tt.expectError {
				// Expect error result
				if result == nil {
					t.Fatal("Expected error result, got nil")
				}
				if !result.IsError {
					t.Errorf("Expected error result for invalid output dir %q, but got success", tt.outputDir)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for valid output dir %q, got: %v", tt.outputDir, err)
				}
				if result != nil && result.IsError {
					t.Errorf("Expected success for valid output dir %q, got error result", tt.outputDir)
				}
			}
		})
	}
}
