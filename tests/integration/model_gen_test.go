package integration_test

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zeromicro/mcp-zero/tools"
)

func TestGenerateModelValidation(t *testing.T) {
	tests := []struct {
		name       string
		params     tools.GenerateModelParams
		shouldFail bool
	}{
		{
			name: "invalid source type",
			params: tools.GenerateModelParams{
				SourceType: "invalid",
				Source:     "user:pass@localhost/db",
				Table:      "users",
			},
			shouldFail: true,
		},
		{
			name: "missing source",
			params: tools.GenerateModelParams{
				SourceType: "mysql",
				Source:     "",
				Table:      "users",
			},
			shouldFail: true,
		},
		{
			name: "missing table",
			params: tools.GenerateModelParams{
				SourceType: "mysql",
				Source:     "user:pass@localhost/db",
				Table:      "",
			},
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _, err := tools.GenerateModel(context.Background(), &mcp.CallToolRequest{}, tt.params)

			if tt.shouldFail {
				if err == nil && (result == nil || !result.IsError) {
					t.Error("Expected validation error")
				}
			}
		})
	}
}
