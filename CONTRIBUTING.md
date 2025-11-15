# Contributing to mcp-zero

Thank you for your interest in contributing to mcp-zero! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Pull Requests](#submitting-pull-requests)
- [Coding Standards](#coding-standards)
- [Tool Development Guidelines](#tool-development-guidelines)

## Code of Conduct

This project follows the [Go Community Code of Conduct](https://go.dev/conduct). Please be respectful and constructive in all interactions.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/your-username/mcp-zero.git
   cd mcp-zero
   ```
3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/zeromicro/mcp-zero.git
   ```
4. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Setup

### Prerequisites

- **Go 1.19 or later** (required by go-zero framework and MCP SDK)
- **goctl** (go-zero CLI tool):
  ```bash
  go install github.com/zeromicro/go-zero/tools/goctl@latest
  ```
- **Claude Desktop** (or other MCP-compatible client) for testing

### Install Dependencies

```bash
go mod download
```

### Build the Project

```bash
go build -o mcp-zero
```

### Run Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific test
go test -v ./tests/integration/api_test.go
```

## Project Structure

```
mcp-zero/
├── main.go                    # MCP server entry point, tool registration
├── internal/
│   ├── analyzer/              # Project structure analysis
│   ├── docs/                  # Documentation concepts and queries
│   ├── errors/                # Error response handling
│   ├── fixer/                 # Code fixing utilities (imports, modules)
│   ├── goctl/                 # goctl discovery and execution
│   ├── logging/               # Structured logging
│   ├── metrics/               # Performance monitoring
│   ├── responses/             # Tool response formatting
│   ├── security/              # Credential handling
│   ├── templates/             # Code generation templates
│   └── validation/            # Input validation
├── tools/
│   ├── api/                   # API service tools
│   ├── rpc/                   # RPC service tools
│   ├── model/                 # Database model tools
│   ├── spec/                  # API spec tools
│   ├── analyze/               # Project analysis tools
│   ├── config/                # Configuration tools
│   ├── template/              # Template generation tools
│   └── query_docs/            # Documentation query tools
├── tests/
│   ├── integration/           # End-to-end integration tests
│   └── unit/                  # Unit tests for internal packages
└── specs/                     # Feature specifications and planning
```

## Making Changes

### Branch Naming Convention

- `feature/` - New features (e.g., `feature/add-dockerfile-template`)
- `fix/` - Bug fixes (e.g., `fix/port-validation`)
- `docs/` - Documentation updates (e.g., `docs/update-readme`)
- `refactor/` - Code refactoring (e.g., `refactor/simplify-analyzer`)
- `test/` - Test additions or fixes (e.g., `test/add-validation-tests`)

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Test additions or changes
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `chore`: Maintenance tasks

**Examples:**

```
feat(api): add support for custom middleware templates

fix(validation): reject service names with hyphens

docs(readme): update installation instructions

test(fixer): add unit tests for import path fixing
```

### Code Changes

1. **Keep changes focused** - One feature or fix per pull request
2. **Write tests** - Add unit tests for new functionality
3. **Update documentation** - Update README.md if behavior changes
4. **Run tests locally** - Ensure all tests pass before submitting
5. **Format code** - Use `gofmt` or `go fmt ./...`

## Testing

### Unit Tests

Unit tests verify individual packages in isolation:

```bash
# Test validation package
go test -v ./internal/validation/

# Test fixer package
go test -v ./internal/fixer/

# Test analyzer package
go test -v ./internal/analyzer/
```

**Writing unit tests:**

```go
package validation

import "testing"

func TestValidateServiceName(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid name", "userservice", false},
        {"empty name", "", true},
        {"with hyphen", "user-service", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateServiceName(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateServiceName() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Integration Tests

Integration tests verify tools work end-to-end with real goctl execution:

```bash
# Run integration tests
go test -v ./tests/integration/

# Run specific integration test
go test -v ./tests/integration/api_test.go
```

**Writing integration tests:**

```go
package integration

import (
    "testing"
    "os"
)

func TestCreateAPIService(t *testing.T) {
    // Create temp directory
    tmpDir, err := os.MkdirTemp("", "mcp-zero-test-")
    if err != nil {
        t.Fatal(err)
    }
    defer os.RemoveAll(tmpDir)

    // Test service creation
    tool := &CreateAPIServiceTool{}
    result, err := tool.Execute(map[string]interface{}{
        "service_name": "testservice",
        "port": 8080,
        "output_dir": tmpDir,
    })

    if err != nil {
        t.Fatalf("Execute() failed: %v", err)
    }

    // Verify result contains expected fields
    if result["service_name"] != "testservice" {
        t.Errorf("Expected service_name=testservice, got %v", result["service_name"])
    }
}
```

### Testing with Claude Desktop

1. **Build the binary**:
   ```bash
   go build -o mcp-zero
   ```

2. **Update Claude Desktop config** with your local binary path:
   ```json
   {
     "mcpServers": {
       "go-zero-dev": {
         "command": "/path/to/your/mcp-zero",
         "env": {
           "GOCTL_PATH": "/path/to/goctl"
         }
       }
     }
   }
   ```

3. **Restart Claude Desktop**

4. **Test your changes** through natural language:
   ```
   Create a new API service called testservice on port 9090
   ```

## Submitting Pull Requests

1. **Push your branch** to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create a pull request** on GitHub

3. **Fill out the PR template**:
   - Description of changes
   - Related issue (if any)
   - Testing performed
   - Screenshots (if UI changes)

4. **Respond to feedback** - Maintainers may request changes

5. **Keep your PR updated**:
   ```bash
   git fetch upstream
   git rebase upstream/main
   git push -f origin feature/your-feature-name
   ```

## Coding Standards

### Go Style Guidelines

Follow the [Effective Go](https://go.dev/doc/effective_go) guidelines and [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments).

**Key points:**

- Use `gofmt` to format code
- Follow naming conventions (camelCase for unexported, PascalCase for exported)
- Write clear, concise comments for exported functions
- Use meaningful variable names
- Keep functions focused and small
- Handle errors explicitly (don't ignore them)

### Error Handling

```go
// Good: Return detailed errors
if err := validateInput(input); err != nil {
    return nil, fmt.Errorf("input validation failed: %w", err)
}

// Good: Provide context in error messages
if !fileExists(path) {
    return nil, fmt.Errorf("file not found: %s (checked: %s)", filename, path)
}
```

### Logging

Use the structured logger from `internal/logging`:

```go
logger := logging.NewLogger(logging.InfoLevel)
logger.Info("Starting service creation", "service", serviceName)
logger.Error("Failed to create service", "error", err, "service", serviceName)
```

**Log levels:**
- **Debug**: Detailed diagnostic information
- **Info**: General informational messages
- **Warn**: Warning messages for potential issues
- **Error**: Error messages for failures

### Metrics

Record metrics for tool operations:

```go
metrics := metrics.NewMetrics()
start := time.Now()
defer metrics.RecordRequest("create_api_service", time.Since(start), err == nil)
```

## Tool Development Guidelines

### Creating a New Tool

1. **Create tool package** in `tools/<category>/`
2. **Implement tool struct** with handler function
3. **Add input validation** using `internal/validation`
4. **Implement core logic** with error handling
5. **Format response** using `internal/responses`
6. **Register tool** in `main.go`
7. **Add integration tests** in `tests/integration/`
8. **Update README.md** with tool documentation

### Tool Structure Template

```go
package mytool

import (
    "github.com/zeromicro/mcp-zero/internal/validation"
    "github.com/zeromicro/mcp-zero/internal/responses"
)

// MyTool implements a sample MCP tool
type MyTool struct{}

// Execute runs the tool with given parameters
func (t *MyTool) Execute(params map[string]interface{}) (map[string]interface{}, error) {
    // 1. Extract and validate inputs
    serviceName, ok := params["service_name"].(string)
    if !ok {
        return responses.ErrorResponse("service_name is required and must be a string")
    }

    if err := validation.ValidateServiceName(serviceName); err != nil {
        return responses.ErrorResponse("invalid service name: " + err.Error())
    }

    // 2. Perform tool operation
    result, err := doWork(serviceName)
    if err != nil {
        return responses.ErrorResponse("operation failed: " + err.Error())
    }

    // 3. Return success response
    return responses.SuccessResponse(map[string]interface{}{
        "service_name": serviceName,
        "result": result,
    })
}
```

### Tool Registration in main.go

```go
// Register tool in main() function
server.AddTool(mcp.Tool{
    Name:        "my_tool_name",
    Description: "Brief description of what this tool does",
    InputSchema: mcp.ToolInputSchema{
        Type: "object",
        Properties: map[string]interface{}{
            "service_name": map[string]interface{}{
                "type":        "string",
                "description": "Name of the service",
            },
        },
        Required: []string{"service_name"},
    },
}, func(args map[string]interface{}) (*mcp.CallToolResult, error) {
    tool := &mytool.MyTool{}
    result, err := tool.Execute(args)
    if err != nil {
        return mcp.NewToolResultError(err.Error()), nil
    }
    return mcp.NewToolResultText(formatResult(result)), nil
})
```

### Input Validation Guidelines

- **Validate all inputs** before processing
- **Use existing validators** from `internal/validation`
- **Provide clear error messages** with examples
- **Check for required parameters** explicitly
- **Validate types** (string, int, bool)
- **Validate formats** (paths, ports, names)

### Response Formatting Guidelines

- **Use consistent structure** from `internal/responses`
- **Include all relevant fields** in success responses
- **Provide actionable errors** with suggestions
- **Include next steps** for users
- **Format paths** as absolute paths
- **Include status indicators** (✅ for success, ❌ for errors)

### Testing Guidelines for Tools

1. **Test happy path** - Valid inputs produce expected outputs
2. **Test error cases** - Invalid inputs produce clear errors
3. **Test edge cases** - Empty strings, special characters, etc.
4. **Test file operations** - Use temporary directories
5. **Test goctl integration** - Verify goctl commands execute correctly
6. **Test cleanup** - Remove temporary files after tests

## Questions or Issues?

- **Documentation**: Check [README.md](./readme.md) for usage instructions
- **Bug Reports**: Open an issue on GitHub
- **Feature Requests**: Open an issue with detailed description
- **Questions**: Open a discussion on GitHub Discussions

Thank you for contributing to mcp-zero! 🎉
