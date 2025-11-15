package docs_test

import (
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zeromicro/mcp-zero/tools"
)

func TestQueryDocs(t *testing.T) {
	ctx := context.Background()
	req := &mcp.CallToolRequest{}

	tests := []struct {
		name          string
		query         string
		expectSuccess bool
		expectConcept bool
		expectMigrate bool
		checkKeywords []string
	}{
		{
			name:          "query middleware concept",
			query:         "explain middleware in go-zero",
			expectSuccess: true,
			expectConcept: true,
			expectMigrate: false,
			checkKeywords: []string{"Middleware"},
		},
		{
			name:          "query jwt authentication",
			query:         "how to use JWT authentication",
			expectSuccess: true,
			expectConcept: true,
			expectMigrate: false,
			checkKeywords: []string{"JWT", "Auth"},
		},
		{
			name:          "query service context",
			query:         "what is service context",
			expectSuccess: true,
			expectConcept: true,
			expectMigrate: false,
			checkKeywords: []string{"Service Context", "injection"},
		},
		{
			name:          "query rpc services",
			query:         "how do RPC services work",
			expectSuccess: true,
			expectConcept: true,
			expectMigrate: false,
			checkKeywords: []string{"RPC", "gRPC"},
		},
		{
			name:          "query cache",
			query:         "how to use cache",
			expectSuccess: true,
			expectConcept: true,
			expectMigrate: false,
			checkKeywords: []string{"Cache", "Redis"},
		},
		{
			name:          "migration from gin",
			query:         "how to migrate from Gin to go-zero",
			expectSuccess: true,
			expectConcept: false,
			expectMigrate: true,
			checkKeywords: []string{"Gin", "migrate"},
		},
		{
			name:          "migration from spring boot",
			query:         "migrating from Spring Boot",
			expectSuccess: true,
			expectConcept: false,
			expectMigrate: true,
			checkKeywords: []string{"Spring Boot"},
		},
		{
			name:          "migration from echo",
			query:         "echo to go-zero migration",
			expectSuccess: true,
			expectConcept: false,
			expectMigrate: true,
			checkKeywords: []string{"Echo"},
		},
		{
			name:          "migration from express",
			query:         "nodejs express migration guide",
			expectSuccess: true,
			expectConcept: false,
			expectMigrate: true,
			checkKeywords: []string{"Express"},
		},
		{
			name:          "empty query",
			query:         "",
			expectSuccess: false,
			expectConcept: false,
			expectMigrate: false,
		},
		{
			name:          "no results query",
			query:         "some random query that doesn't match anything",
			expectSuccess: true,
			expectConcept: false,
			expectMigrate: false,
			checkKeywords: []string{"Available Topics"},
		},
		{
			name:          "general framework query",
			query:         "go-zero best practices",
			expectSuccess: true,
			expectConcept: true,
			expectMigrate: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := tools.QueryDocsParams{
				Query: tt.query,
			}

			result, _, err := tools.QueryDocs(ctx, req, params)

			if tt.expectSuccess {
				if err != nil {
					t.Errorf("Expected success, got error: %v", err)
					return
				}
				if result.IsError {
					t.Errorf("Expected IsError=false, got IsError=true with content: %v", result.Content)
					return
				}

				// Check response content
				if len(result.Content) == 0 {
					t.Error("Expected response content")
					return
				}

				// Get text from response (cast to TextContent)
				var responseText string
				if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
					responseText = textContent.Text
				} else {
					t.Error("Expected TextContent in response")
					return
				}

				// Check for expected concept or migration content
				if tt.expectConcept {
					if !strings.Contains(responseText, "Concepts") && !strings.Contains(responseText, "concept") {
						t.Error("Expected concept documentation in response")
					}
				}

				if tt.expectMigrate {
					if !strings.Contains(responseText, "Migration") && !strings.Contains(responseText, "Migrating") {
						t.Error("Expected migration guide in response")
					}
				}

				// Check for specific keywords
				for _, keyword := range tt.checkKeywords {
					if !strings.Contains(responseText, keyword) {
						t.Errorf("Expected response to contain '%s', but it doesn't", keyword)
					}
				}

			} else {
				if err == nil && !result.IsError {
					t.Errorf("Expected error or IsError=true, got success")
				}
			}
		})
	}
}

func TestQueryDocsKeywordExtraction(t *testing.T) {
	ctx := context.Background()
	req := &mcp.CallToolRequest{}

	tests := []struct {
		name            string
		query           string
		expectedKeyword string
	}{
		{
			name:            "middleware keyword",
			query:           "how does middleware work",
			expectedKeyword: "middleware",
		},
		{
			name:            "jwt keyword",
			query:           "JWT authentication guide",
			expectedKeyword: "jwt",
		},
		{
			name:            "cache keyword",
			query:           "using Redis cache",
			expectedKeyword: "cache",
		},
		{
			name:            "rpc keyword",
			query:           "RPC service communication",
			expectedKeyword: "rpc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := tools.QueryDocsParams{
				Query: tt.query,
			}

			result, _, err := tools.QueryDocs(ctx, req, params)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if len(result.Content) == 0 {
				t.Error("Expected response content")
				return
			}

			// Verify the response is relevant to the keyword
			if textContent, ok := result.Content[0].(*mcp.TextContent); ok {
				if !strings.Contains(strings.ToLower(textContent.Text), tt.expectedKeyword) {
					t.Errorf("Expected response to be relevant to keyword '%s'", tt.expectedKeyword)
				}
			}
		})
	}
}

func TestQueryDocsResponseFormat(t *testing.T) {
	ctx := context.Background()
	req := &mcp.CallToolRequest{}

	params := tools.QueryDocsParams{
		Query: "middleware",
	}

	result, _, err := tools.QueryDocs(ctx, req, params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(result.Content) == 0 {
		t.Fatal("Expected response content")
	}

	textContent, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatal("Expected TextContent in response")
	}

	responseText := textContent.Text

	// Check response structure
	expectedSections := []string{
		"Documentation Results",
		"Category",
		"Example",
		"Tips",
	}

	for _, section := range expectedSections {
		if !strings.Contains(responseText, section) {
			t.Errorf("Expected response to contain section '%s'", section)
		}
	}

	// Check code formatting
	if !strings.Contains(responseText, "```go") {
		t.Error("Expected response to contain formatted code examples")
	}
}
