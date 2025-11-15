package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/zeromicro/mcp-zero/internal/docs"
	"github.com/zeromicro/mcp-zero/internal/responses"
)

type QueryDocsParams struct {
	Query string `json:"query"` // The documentation query
}

// QueryDocs queries go-zero framework documentation and migration guides
func QueryDocs(ctx context.Context, req *mcp.CallToolRequest, params QueryDocsParams) (*mcp.CallToolResult, any, error) {
	if params.Query == "" {
		return responses.FormatValidationError("query", "", "query is required",
			"Provide a question about go-zero concepts or migration guidance")
	}

	query := strings.TrimSpace(params.Query)
	keywords := extractKeywords(query)

	// Search concepts
	concepts := docs.SearchConcepts(query)

	// Search migration guides
	migrations := docs.SearchMigrationGuides(query)

	// Build response
	if len(concepts) == 0 && len(migrations) == 0 {
		return formatNoResultsResponse(query, keywords)
	}

	message := formatDocsResponse(query, concepts, migrations)

	data := map[string]any{
		"query":            query,
		"concepts_found":   len(concepts),
		"migrations_found": len(migrations),
		"keywords":         keywords,
	}

	return responses.FormatSuccessWithData(message, data)
}

func extractKeywords(query string) []string {
	query = strings.ToLower(query)

	// Common go-zero related keywords
	knownKeywords := []string{
		"middleware", "api", "rpc", "model", "cache", "jwt", "configuration",
		"error", "validation", "context", "service", "handler", "logic",
		"grpc", "protobuf", "redis", "mysql", "database",
		"gin", "echo", "spring", "express", "migrate", "migration",
	}

	var found []string
	for _, keyword := range knownKeywords {
		if strings.Contains(query, keyword) {
			found = append(found, keyword)
		}
	}

	return found
}

func formatDocsResponse(query string, concepts []docs.Concept, migrations []docs.MigrationGuide) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("📚 Documentation Results for: \"%s\"\n\n", query))

	// Format concepts
	if len(concepts) > 0 {
		sb.WriteString("## Concepts\n\n")
		for i, concept := range concepts {
			if i >= 3 {
				sb.WriteString(fmt.Sprintf("... and %d more concepts\n\n", len(concepts)-3))
				break
			}

			sb.WriteString(fmt.Sprintf("### %s\n", concept.Name))
			sb.WriteString(fmt.Sprintf("**Category**: %s\n\n", concept.Category))
			sb.WriteString(fmt.Sprintf("%s\n\n", concept.Description))

			if concept.Example != "" {
				sb.WriteString("**Example:**\n```go\n")
				sb.WriteString(concept.Example)
				sb.WriteString("\n```\n\n")
			}

			if len(concept.RelatedDocs) > 0 {
				sb.WriteString("**Related Documentation:**\n")
				for _, link := range concept.RelatedDocs {
					sb.WriteString(fmt.Sprintf("- %s\n", link))
				}
				sb.WriteString("\n")
			}

			sb.WriteString("---\n\n")
		}
	}

	// Format migration guides
	if len(migrations) > 0 {
		sb.WriteString("## Migration Guides\n\n")
		for i, guide := range migrations {
			if i >= 2 {
				sb.WriteString(fmt.Sprintf("... and %d more migration guides\n\n", len(migrations)-2))
				break
			}

			sb.WriteString(fmt.Sprintf("### Migrating from %s to %s\n", guide.FromFramework, guide.ToGoZero))
			sb.WriteString(fmt.Sprintf("**Difficulty**: %s\n\n", guide.Difficulty))
			sb.WriteString(fmt.Sprintf("**Key Differences:**\n%s\n\n", guide.KeyDifferences))

			if guide.Example != "" {
				sb.WriteString("**Example Comparison:**\n```go\n")
				sb.WriteString(guide.Example)
				sb.WriteString("\n```\n\n")
			}

			if len(guide.Steps) > 0 {
				sb.WriteString("**Migration Steps:**\n")
				for _, step := range guide.Steps {
					sb.WriteString(fmt.Sprintf("%s\n", step))
				}
				sb.WriteString("\n")
			}

			sb.WriteString("---\n\n")
		}
	}

	// Add helpful tips
	sb.WriteString("## 💡 Tips\n\n")
	sb.WriteString("- Use specific keywords like 'middleware', 'jwt', 'cache' for better results\n")
	sb.WriteString("- Ask about specific frameworks when looking for migration guides (e.g., 'gin', 'spring')\n")
	sb.WriteString("- Check the official documentation links for more detailed information\n")

	return sb.String()
}

func formatNoResultsResponse(query string, keywords []string) (*mcp.CallToolResult, any, error) {
	message := fmt.Sprintf("No documentation found for: \"%s\"\n\n", query)

	message += "## Available Topics\n\n"
	message += "**Core Concepts:**\n"
	message += "- Middleware - HTTP middleware pattern\n"
	message += "- API Definition - .api file syntax and structure\n"
	message += "- Service Context - Dependency injection container\n"
	message += "- RPC Services - gRPC microservices\n"
	message += "- Database Models - Model generation and usage\n"
	message += "- Configuration - YAML config management\n"
	message += "- Error Handling - Structured error handling\n"
	message += "- JWT Authentication - Built-in JWT support\n"
	message += "- Cache - Redis caching integration\n"
	message += "- Validation - Request validation\n\n"

	message += "**Migration Guides:**\n"
	message += "- Gin to go-zero\n"
	message += "- Echo to go-zero\n"
	message += "- Spring Boot to go-zero\n"
	message += "- Express.js to go-zero\n"
	message += "- vanilla gRPC to go-zero RPC\n\n"

	message += "💡 Try asking:\n"
	message += "- \"How does middleware work in go-zero?\"\n"
	message += "- \"How to migrate from Gin to go-zero?\"\n"
	message += "- \"Explain service context\"\n"
	message += "- \"How to use JWT authentication?\"\n"

	data := map[string]any{
		"query":    query,
		"keywords": keywords,
		"found":    false,
	}

	return responses.FormatSuccessWithData(message, data)
}
