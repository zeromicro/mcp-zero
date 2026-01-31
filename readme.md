# Go-Zero MCP Tool

A Model Context Protocol (MCP) tool that helps developers quickly scaffold and generate go-zero projects with ease.

## Quick Start

**New to mcp-zero?** Check out our [Quick Start Guide](QUICKSTART.md) for a step-by-step tutorial on getting started!

The quickstart covers:

- Installation and configuration
- Creating your first API service
- Common use cases and workflows
- Integration with Claude Desktop

## Features

### Core Service Generation

- **Create API Services**: Generate new REST API services with customizable ports and styles
- **Create RPC Services**: Generate gRPC services from protobuf definitions
- **Generate API Code**: Convert API specification files to Go code
- **Generate Models**: Create database models from various sources (MySQL, PostgreSQL, MongoDB, DDL)
- **Create API Specs**: Generate sample API specification files

### Advanced Features

- **Analyze Projects**: Analyze existing go-zero projects to understand structure and dependencies
- **Manage Configuration**: Generate configuration files with proper structure validation
- **Generate Templates**: Create middleware, error handlers, and deployment templates
- **Query Documentation**: Access go-zero concepts and migration guides from other frameworks
- **Validate Input**: Comprehensive validation for API specs, protobuf definitions, and configurations

## Prerequisites

1. **Go** (1.19 or later)

2. **go-zero CLI (goctl)**: Install with `go install github.com/zeromicro/go-zero/tools/goctl@latest`

3. **Claude Desktop** (or other MCP-compatible client)

For detailed installation instructions, see the [Quick Start Guide](QUICKSTART.md).

## Installation

1. Create a new directory for your MCP tool:

   ```bash
   mkdir go-zero-mcp && cd go-zero-mcp
   ```

2. Initialize Go module:

   ```bash
   go mod init go-zero-mcp
   ```

3. Install dependencies:

   ```bash
   go get github.com/modelcontextprotocol/go-sdk
   go get gopkg.in/yaml.v3
   ```

4. Save the main tool code as `main.go`

5. Build the tool:

   ```bash
   go build -o go-zero-mcp main.go
   ```

## Configuration for Claude Desktop

Add this configuration to your Claude Desktop MCP settings:

### macOS

Edit `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "mcp-zero": {
      "command": "/path/to/your/mcp-zero",
      "env": {
        "GOCTL_PATH": "/Users/yourname/go/bin/goctl"
      }
    }
  }
}
```

### Linux

Edit `~/.config/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "mcp-zero": {
      "command": "/path/to/your/mcp-zero",
      "env": {
        "GOCTL_PATH": "/usr/local/bin/goctl"
      }
    }
  }
}
```

### Windows

Edit `%APPDATA%\Claude\claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "mcp-zero": {
      "command": "C:\\path\\to\\your\\mcp-zero.exe",
      "env": {
        "GOCTL_PATH": "C:\\Go\\bin\\goctl.exe"
      }
    }
  }
}
```

## Available Tools

### 1. create_api_service

Creates a new go-zero API service.

**Parameters:**

- `service_name` (required): Name of the API service
- `port` (optional): Port number (default: 8888)
- `style` (optional): Code style - "go_zero" or "gozero" (default: "go_zero")
- `output_dir` (optional): Output directory (default: current directory)

### 2. create_rpc_service

Creates a new go-zero RPC service from protobuf definition.

**Parameters:**

- `service_name` (required): Name of the RPC service
- `proto_content` (required): Protobuf definition content
- `output_dir` (optional): Output directory (default: current directory)

### 3. generate_api_from_spec

Generates go-zero API code from an API specification file.

**Parameters:**

- `api_file` (required): Path to the .api specification file
- `output_dir` (optional): Output directory (default: current directory)
- `style` (optional): Code style - "go_zero" or "gozero" (default: "go_zero")

### 4. generate_model

Generates database model code from database schema.

**Parameters:**

- `source_type` (required): Source type - "mysql", "postgresql", "mongo", or "ddl"
- `source` (required): Database connection string or DDL file path
- `table` (optional): Specific table name (for database sources)
- `output_dir` (optional): Output directory (default: "./model")

### 5. create_api_spec

Creates a sample API specification file.

**Parameters:**

- `service_name` (required): Name of the API service
- `endpoints` (required): Array of endpoint objects with method, path, and handler
- `output_file` (optional): Output file path (default: service_name.api)

### 6. analyze_project

Analyzes an existing go-zero project structure and dependencies.

**Parameters:**

- `project_dir` (required): Path to the project directory
- `analysis_type` (optional): Type of analysis - "api", "rpc", "model", or "full" (default: "full")

### 7. generate_config

Generates configuration files for go-zero services.

**Parameters:**

- `service_name` (required): Name of the service
- `service_type` (required): Service type - "api" or "rpc"
- `config_type` (optional): Configuration type - "dev", "test", or "prod" (default: "dev")
- `output_file` (optional): Output file path (default: etc/{service_name}.yaml)

### 8. generate_template

Generates common code templates for go-zero services.

**Parameters:**

- `template_type` (required): Template type - "middleware", "error_handler", "dockerfile", "docker_compose", or "kubernetes"
- `service_name` (required): Name of the service
- `output_path` (optional): Output file path (uses defaults based on template type)

### 9. query_docs

Queries go-zero documentation and migration guides.

**Parameters:**

- `query` (required): Natural language query about go-zero concepts or migration
- `doc_type` (optional): Documentation type - "concept", "migration", or "both" (default: "both")

### 10. validate_input

Validates API specs, protobuf definitions, or configuration files.

**Parameters:**

- `input_type` (required): Input type - "api_spec", "proto", or "config"
- `content` (required): Content to validate
- `strict` (optional): Enable strict validation mode (default: false)

## Usage Examples

### Creating a New API Service

```text
Please create a new go-zero API service called "user-service" on port 8080
```

### Creating an RPC Service

```text
Create a new go-zero RPC service called "auth-service" with this protobuf definition:

syntax = "proto3";

package auth;

option go_package = "./auth";

service AuthService {
  rpc Login(LoginRequest) returns (LoginResponse);
  rpc Logout(LogoutRequest) returns (LogoutResponse);
}

message LoginRequest {
  string username = 1;
  string password = 2;
}

message LoginResponse {
  string token = 1;
  int64 expires_at = 2;
}

message LogoutRequest {
  string token = 1;
}

message LogoutResponse {
  bool success = 1;
}
```

### Generating Models from Database

```text
Generate go-zero models from my MySQL database with connection string "user:password@tcp(localhost:3306)/mydb"
```

### Creating an API Specification

```text
Create an API specification for a "blog-service" with these endpoints:
- GET /api/posts (handler: GetPostsHandler)
- POST /api/posts (handler: CreatePostHandler)
- GET /api/posts/:id (handler: GetPostHandler)
- PUT /api/posts/:id (handler: UpdatePostHandler)
- DELETE /api/posts/:id (handler: DeletePostHandler)
```

### Analyzing a Project

```text
Analyze my go-zero project in /path/to/myproject to understand its structure
```

### Generating Configuration

```text
Generate a production configuration file for my "order-service" API service
```

### Generating Templates

```text
Generate a middleware template for my "auth-service"
```

### Querying Documentation

```text
How do I implement JWT authentication in go-zero?
```

```text
How do I migrate from Express.js to go-zero?
```

### Validating Input

```text
Validate this API spec file at /path/to/service.api with strict mode enabled
```

## Project Structure

After building, your MCP server will have the following structure:

```text
mcp-zero/
├── main.go                    # Entry point and tool registration
├── tools/                     # Tool implementations
│   ├── create_api_service.go
│   ├── create_rpc_service.go
│   ├── generate_api.go
│   ├── generate_model.go
│   ├── create_api_spec.go
│   ├── analyze_project.go
│   ├── generate_config.go
│   ├── generate_template.go
│   ├── query_docs.go
│   └── validate_input.go
├── internal/                  # Internal packages
│   ├── analyzer/             # Project analysis
│   ├── validation/           # Input validation
│   ├── security/             # Credential handling
│   ├── templates/            # Code templates
│   ├── docs/                 # Documentation database
│   ├── logging/              # Structured logging
│   └── metrics/              # Performance metrics
└── tests/                     # Test suites
    ├── integration/
    └── unit/
```

## Architecture

The MCP server is built with:

- **MCP SDK**: Uses github.com/modelcontextprotocol/go-sdk for protocol implementation
- **Transport**: stdio-based communication with Claude Desktop
- **Code Generation**: Leverages go-zero's goctl CLI tool for generating production-ready code
- **Validation**: Comprehensive input validation for safety and correctness
- **Security**: Safe credential handling with environment variable substitution
- **Observability**: Built-in logging and metrics for monitoring tool performance

## Best Practices

1. **Service Naming**: Use lowercase with hyphens (e.g., "user-service", "auth-api")
2. **Port Configuration**: Choose unique ports for each service (8080-8090 range recommended)
3. **Code Style**: Stick to "go_zero" style for consistency with official conventions
4. **Configuration**: Use environment-specific configs (dev, test, prod)
5. **Documentation**: Query docs regularly to stay aligned with go-zero best practices
6. **Validation**: Always validate inputs before generation to catch errors early

## Troubleshooting

### Common Issues

1. **goctl command not found**: Make sure goctl is installed and in your PATH
2. **Permission denied**: Ensure the MCP tool executable has proper permissions
3. **Database connection errors**: Verify connection strings and database accessibility

### Debug Mode

To enable debug logging, set the environment variable:

```bash
export MCP_DEBUG=1
```

## Contributing

Feel free to extend this tool with additional go-zero features such as:

- Dockerfile generation
- Kubernetes manifest generation
- Docker Compose file creation
- API documentation generation
- Testing template creation

## License

MIT License - see [LICENSE](LICENSE) file for details.

Copyright (c) 2025 go-zero team
