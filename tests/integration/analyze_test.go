package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zeromicro/mcp-zero/tools"
)

func TestAnalyzeProject(t *testing.T) {
	tests := []struct {
		name             string
		setupFunc        func(string) error
		expectError      bool
		checkServices    bool
		expectedAPICount int
		expectedRPCCount int
	}{
		{
			name: "project with API service",
			setupFunc: func(dir string) error {
				apiSpec := "syntax = \"v1\"\n\ninfo(\n\ttitle: \"Test API\"\n\tversion: \"1.0\"\n)\n\ntype Request {\n\tName string `json:\"name\"`\n}\n\ntype Response {\n\tMessage string `json:\"message\"`\n}\n\nservice test {\n\t@handler TestHandler\n\tpost /test (Request) returns (Response)\n}\n"
				return os.WriteFile(filepath.Join(dir, "test.api"), []byte(apiSpec), 0644)
			},
			expectError:      false,
			checkServices:    true,
			expectedAPICount: 1,
			expectedRPCCount: 0,
		},
		{
			name: "project with RPC service",
			setupFunc: func(dir string) error {
				protoSpec := "syntax = \"proto3\";\n\npackage test;\n\noption go_package = \"./test\";\n\nmessage Request {\n  string name = 1;\n}\n\nmessage Response {\n  string message = 1;\n}\n\nservice TestService {\n  rpc Test(Request) returns (Response);\n}\n"
				return os.WriteFile(filepath.Join(dir, "test.proto"), []byte(protoSpec), 0644)
			},
			expectError:      false,
			checkServices:    true,
			expectedAPICount: 0,
			expectedRPCCount: 1,
		},
		{
			name: "project with go.mod",
			setupFunc: func(dir string) error {
				goMod := "module test\n\ngo 1.19\n\nrequire (\n\tgithub.com/zeromicro/go-zero v1.6.0\n\tgoogle.golang.org/grpc v1.60.0\n)\n"
				return os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goMod), 0644)
			},
			expectError:   false,
			checkServices: false,
		},
		{
			name: "empty project",
			setupFunc: func(dir string) error {
				return nil
			},
			expectError:   false,
			checkServices: false,
		},
		{
			name: "non-existent directory",
			setupFunc: func(dir string) error {
				os.RemoveAll(dir)
				return nil
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			if err := tt.setupFunc(tmpDir); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}

			params := tools.AnalyzeProjectParams{
				ProjectPath: tmpDir,
			}

			result, _, err := tools.AnalyzeProject(context.Background(), &mcp.CallToolRequest{}, params)

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

			if tt.checkServices {
				if len(result.Content) == 0 {
					t.Errorf("Result content is empty")
				}
			}
		})
	}
}

func TestAnalyzeProjectCache(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := "module testcache\n\ngo 1.19\n\nrequire github.com/zeromicro/go-zero v1.6.0\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	params := tools.AnalyzeProjectParams{
		ProjectPath: tmpDir,
	}

	result1, data1, err := tools.AnalyzeProject(context.Background(), &mcp.CallToolRequest{}, params)
	if err != nil {
		t.Fatalf("First call failed: %v", err)
	}
	if result1.IsError {
		t.Fatalf("First call returned error: %v", result1.Content)
	}

	result2, data2, err := tools.AnalyzeProject(context.Background(), &mcp.CallToolRequest{}, params)
	if err != nil {
		t.Fatalf("Second call failed: %v", err)
	}
	if result2.IsError {
		t.Fatalf("Second call returned error: %v", result2.Content)
	}

	if data1 == nil || data2 == nil {
		t.Error("Expected data to be returned")
	}

	dataMap, ok := data2.(map[string]any)
	if !ok {
		t.Fatalf("Expected data to be a map")
	}

	fromCache, ok := dataMap["from_cache"].(bool)
	if !ok {
		t.Errorf("Expected from_cache field in data")
	}

	if !fromCache {
		t.Errorf("Second call should be from cache")
	}
}
