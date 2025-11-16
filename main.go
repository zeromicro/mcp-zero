package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/zeromicro/mcp-zero/tools"
)

func main() {
	// Create MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "mcp-zero",
		Version: "1.0.0",
	}, nil)

	// Register create_api_service tool (T034 - User Story 1)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_api_service",
		Description: "Create a new go-zero API service with proper structure and configuration",
	}, tools.CreateAPIService)

	// Register generate_api_from_spec tool (T047 - User Story 2)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_api_from_spec",
		Description: "Generate go-zero API code from API specification file",
	}, tools.GenerateAPIFromSpec)

	// Register create_rpc_service tool (T058 - User Story 3)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_rpc_service",
		Description: "Create a new go-zero RPC service with protobuf definition",
	}, tools.CreateRPCService)

	// Register generate_model tool (T071 - User Story 4)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_model",
		Description: "Generate go-zero database model from table schema",
	}, tools.GenerateModel)

	// Register create_api_spec tool (T081 - User Story 5)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "create_api_spec",
		Description: "Create a sample API specification file for go-zero. IMPORTANT: Always define concrete types for request and response - do NOT use 'any' type in .api files as it's not supported by go-zero",
	}, tools.CreateAPISpec)

	// Register analyze_project tool (T097 - User Story 6)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "analyze_project",
		Description: "Analyze existing go-zero project structure and dependencies",
	}, tools.AnalyzeProject)

	// Register validate_config tool (T109 - User Story 7)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "validate_config",
		Description: "Validate go-zero service configuration file",
	}, tools.ValidateConfig)

	// Register generate_config_template tool (T109 - User Story 7)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_config_template",
		Description: "Generate configuration template for go-zero service",
	}, tools.GenerateConfigTemplate)

	// Register generate_template tool (T123 - User Story 8)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "generate_template",
		Description: "Generate common code templates (middleware, error handlers, deployment configs)",
	}, tools.GenerateTemplate)

	// Register query_docs tool (T134 - User Story 9)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "query_docs",
		Description: "Query go-zero framework documentation and migration guides",
	}, tools.QueryDocs)

	// Run the server over stdin/stdout using the StdioTransport
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
