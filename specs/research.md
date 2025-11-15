# Research & Technical Decisions

**Feature**: MCP Tool for go-zero Framework
**Date**: November 14, 2025
**Phase**: 0 (Research & Discovery)

## Overview

This document captures research findings and technical decisions for implementing an MCP server that provides go-zero framework capabilities through conversational AI interfaces.

## Key Technical Decisions

### 1. MCP Protocol Implementation

**Decision**: Use mark3labs/mcp-go SDK v0.37.0+

**Rationale**:
- Official Go implementation of Model Context Protocol
- Mature library with stable API
- Supports stdio transport (required for Claude Desktop integration)
- Provides type-safe tool registration and execution
- Active maintenance and community support

**Alternatives Considered**:
- Building custom MCP implementation: Rejected due to protocol complexity and maintenance burden
- Using other language SDKs with cgo: Rejected due to deployment complexity and performance overhead

**Integration Pattern**:
```go
// Tool registration pattern
server := mcp.NewServer(mcp.WithStdio())
server.AddTool(mcp.Tool{
    Name: "create_api_service",
    Description: "Create new go-zero API service",
    InputSchema: /* JSON schema */,
}, handleCreateAPIService)
```

### 2. goctl Discovery Strategy

**Decision**: Multi-level fallback discovery with environment variable override

**Rationale**:
- MCP servers run in isolated environments without normal shell PATH
- Different users have different installation paths
- Must support both global and local goctl installations
- Constitution requirement: Environment Resilience

**Implementation Strategy**:

1. Check GOCTL_PATH environment variable (highest priority)
2. Search common installation locations:
   - `/usr/local/bin/goctl`
   - `$HOME/go/bin/goctl`
   - `$HOME/Develop/go/bin/goctl` (user's known path)
   - `$GOPATH/bin/goctl`
3. Search PATH directories
4. If not found: Return actionable error with installation instructions

**Error Handling**:
```go
if goctlPath == "" {
    return fmt.Errorf("goctl not found. Install with: go install github.com/zeromicro/go-zero/tools/goctl@latest\nOr set GOCTL_PATH environment variable")
}
```

### 3. Code Generation Workflow

**Decision**: Generate → Fix → Validate → Report pattern

**Rationale**:
- goctl generates working code but with some issues (import paths, config)
- Constitution requirements: Complete Automation, Validation & Safety
- Users expect ready-to-run artifacts without manual intervention

**Workflow Steps**:

1. **Generate**: Execute goctl command with absolute paths
2. **Fix**:
   - Replace absolute import paths with local module names
   - Update configuration files with user-specified parameters
   - Initialize go modules if missing
   - Run `go mod tidy` to resolve dependencies
3. **Validate**:
   - Verify expected files exist
   - Run `go build` to ensure code compiles
   - Check for common issues
4. **Report**: Success with next steps or detailed error with remediation

### 4. Service Name Validation

**Decision**: Proactive validation with helpful corrections

**Rationale**:
- go-zero requires valid Go identifiers (no hyphens, special characters)
- Common mistake: using "user-service" instead of "userservice"
- Constitution requirement: Developer Experience First

**Validation Rules**:
- Must start with letter
- Can contain letters, numbers, underscores only
- No hyphens, spaces, or special characters
- Provide suggestion: "user-service" → "userservice" or "user_service"

### 5. Monorepo Support

**Decision**: Services in subdirectories with independent go modules

**Rationale**:
- Clarification decision: Support multiple services per workspace
- Each service is independent deployable unit
- Avoid complex workspace-level go.mod management
- Aligns with microservice best practices

**Directory Structure**:
```
workspace/
├── api-gateway/
│   ├── go.mod
│   └── [service files]
├── user-service/
│   ├── go.mod
│   └── [service files]
└── order-service/
    ├── go.mod
    └── [service files]
```

### 6. Error Recovery Strategy

**Decision**: Preserve partial state, allow retry with corrections

**Rationale**:
- Clarification decision: Retry with correction
- Conversational interface allows iterative refinement
- Avoid forcing users to start over completely
- Show clear progress: what succeeded, what failed

**Recovery Pattern**:
- Keep successfully generated files
- Report specific failure point
- Accept corrected inputs
- Resume from failure point
- Don't regenerate successful parts

### 7. Configuration Validation

**Decision**: Schema-based validation with custom field support

**Rationale**:
- Clarification decision: Allow custom fields
- Balance safety with flexibility
- go-zero config evolves over time
- Users may have custom configuration needs

**Validation Approach**:
- Define schema for known configuration keys
- Validate types and required fields
- Warn about unknown fields (don't block)
- Check for common misconfigurations (conflicting ports, invalid timeouts)

### 8. Database Credential Handling

**Decision**: Support both connection strings and credential files

**Rationale**:
- Clarification decision: Flexibility for different scenarios
- Connection strings: Quick development and testing
- Credential files: Secure production deployments
- Never log or persist credentials

**Security Measures**:
- Credentials never written to generated code
- Credentials not logged in tool output
- Support for encrypted credential files
- Clear documentation on security best practices

### 9. Project Analysis Approach

**Decision**: Multi-stage parsing with caching

**Rationale**:
- Large projects can have many files
- Parsing can be expensive
- Need to meet <1 minute success criteria

**Analysis Stages**:

1. **Quick Scan**: Find all .api and .proto files
2. **Parse Structure**: Extract endpoints, services, types
3. **Dependency Analysis**: Scan go.mod, imports
4. **Cache Results**: In-memory cache for repeated queries

### 10. Template System

**Decision**: Embedded Go templates with parameter injection

**Rationale**:
- Templates need to integrate with existing code
- Must follow go-zero conventions
- Support customization through parameters

**Template Categories**:
- Middleware (auth, logging, rate limiting)
- Error handlers (custom error types, HTTP responses)
- Deployment configs (Docker, Kubernetes, systemd)

Each template includes:
- Working code with sensible defaults
- Parameter placeholders
- Integration instructions
- Usage examples

## Integration Patterns

### MCP Tool Registration

All tools follow consistent pattern:

```go
type ToolHandler func(ctx context.Context, arguments map[string]interface{}) (string, error)

// Each tool:
// 1. Validates inputs
// 2. Performs operation
// 3. Returns structured result or detailed error
```

### Command Execution

All external commands use safe execution wrapper:

```go
func executeCommand(ctx context.Context, cmd string, args []string, workDir string) (string, error)
    // - Use absolute paths
    // - Set working directory explicitly
    // - Capture stdout/stderr separately
    // - Include timeout
    // - Return structured error with command context
```

### File Operations

All file operations:

```go
// - Check permissions before writing
// - Use atomic writes (write temp, rename)
// - Verify write succeeded
// - Clean up on failure
```

## Testing Strategy

### Unit Tests

- All validation functions
- All fixer functions
- goctl discovery logic
- Template generation
- Error message formatting

### Integration Tests

- End-to-end service creation
- Code generation from real .api files
- Model generation from test databases
- Project analysis on sample projects

### Test Data

- Sample .api specifications
- Sample .proto files
- Test database schemas
- Example configurations

## Performance Considerations

### Optimization Areas

1. **goctl Discovery**: Cache discovered path for session
2. **Project Analysis**: Cache parsed results
3. **File Operations**: Batch related operations
4. **Validation**: Fail fast on first error

### Success Criteria Alignment

- Code generation: <5 seconds (target: 2-3 seconds)
- File operations: <30 seconds (target: 10-20 seconds)
- Project analysis: <1 minute (target: 30-45 seconds)

## Documentation Strategy

### User Documentation

- README with quick start
- Tool-by-tool usage examples
- Common workflows (create service, add endpoint, etc.)
- Troubleshooting guide

### Developer Documentation

- Architecture overview
- Adding new tools guide
- Testing guide
- Contributing guidelines

## Deployment Considerations

### Distribution

- Single binary (go build)
- No external dependencies except goctl
- Configuration via environment variables
- Installation via go install or binary download

### Claude Desktop Integration

1. User installs mcp-zero binary
2. Adds to Claude Desktop config:
```json
{
  "mcpServers": {
    "go-zero": {
      "command": "/path/to/mcp-zero",
      "env": {
        "GOCTL_PATH": "/path/to/goctl"
      }
    }
  }
}
```
3. Restart Claude Desktop
4. Tools available in conversation

## Open Questions & Future Enhancements

### Phase 1+ Considerations

- **Batch Operations**: Create multiple services in one request
- **Project Templates**: Full project scaffolding with multiple services
- **Custom Templates**: User-provided template support
- **Plugin System**: Extensibility for custom tools
- **Web UI**: Optional web interface for non-AI usage

### Known Limitations

- Requires goctl installation (can't bundle due to licensing)
- go-zero version compatibility (document supported versions)
- Platform-specific paths (Windows support needs additional work)

## References

- [Model Context Protocol Specification](https://modelcontextprotocol.io/)
- [mark3labs/mcp-go SDK](https://github.com/mark3labs/mcp-go)
- [go-zero Documentation](https://go-zero.dev/)
- [goctl Tool Reference](https://go-zero.dev/docs/tasks)
