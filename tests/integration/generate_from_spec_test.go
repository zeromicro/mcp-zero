package integration_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zeromicro/mcp-zero/tools"
)

func TestGenerateAPIFromSpecSuccess(t *testing.T) {
	tmpDir := t.TempDir()
	outputDir := filepath.Join(tmpDir, "test-service")
	os.MkdirAll(outputDir, 0755)

	apiFile := filepath.Join(tmpDir, "test.api")
	apiContent := `syntax = "v1"

service test-api {
	@handler TestHandler
	get /test
}
`
	os.WriteFile(apiFile, []byte(apiContent), 0644)

	params := tools.GenerateAPIFromSpecParams{
		APIFile:   apiFile,
		OutputDir: outputDir,
		Style:     "go_zero",
	}

	result, _, err := tools.GenerateAPIFromSpec(context.Background(), &mcp.CallToolRequest{}, params)

	if err != nil {
		t.Errorf("Expected success, got error: %v", err)
	}
	if result == nil || result.IsError {
		t.Error("Expected successful result")
	}
}
