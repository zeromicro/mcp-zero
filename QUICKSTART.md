# Quick Start Guide

**Feature**: MCP Tool for go-zero Framework
**Version**: 1.0.0
**Date**: November 14, 2025

## Overview

mcp-zero is an MCP (Model Context Protocol) server that brings go-zero framework capabilities to AI assistants like Claude Desktop. Use natural language to create services, generate code, and manage go-zero projects.

## Prerequisites

- Go 1.19 or later
- goctl (go-zero CLI tool)
- Claude Desktop (or other MCP-compatible AI assistant)

## Installation

### 1. Install goctl

```bash
go install github.com/zeromicro/go-zero/tools/goctl@latest
```

Verify installation:

```bash
goctl --version
```

### 2. Install mcp-zero

```bash
go install github.com/zeromicro/mcp-zero@latest
```

Or build from source:

```bash
git clone https://github.com/zeromicro/mcp-zero.git
cd mcp-zero
go build -o mcp-zero
```

### 3. Configure Claude Desktop

Add to your Claude Desktop configuration (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "go-zero": {
      "command": "/Users/yourname/go/bin/mcp-zero",
      "env": {
        "GOCTL_PATH": "/Users/yourname/go/bin/goctl"
      }
    }
  }
}
```

**Note**: Adjust paths to match your installation locations.

### 4. Restart Claude Desktop

Quit and restart Claude Desktop to load the new MCP server.

## Quick Tutorial

### Creating Your First API Service

Open Claude Desktop and try:

```
Create a new API service called "userservice" on port 8080
```

Claude will use mcp-zero to:

1. Validate the service name
2. Execute goctl to generate project structure
3. Fix import paths and initialize go modules
4. Verify the build succeeds
5. Provide next steps

Example response:

```
✅ API service 'userservice' created successfully!

Location: /your/current/directory/userservice
Port: 8080
Entry point: userservice.go

Next steps:
1. cd userservice
2. go run userservice.go
3. Test: curl http://localhost:8080/ping

Service is ready to run!
```

### Running the Service

```bash
cd userservice
go run userservice.go
```

Test it:

```bash
curl http://localhost:8080/ping
# Response: {"message":"pong"}
```

### Adding an Endpoint

Ask Claude:

```
Add a new endpoint GET /api/user/:id to my userservice
```

mcp-zero will update the API specification and regenerate code.

### Creating an RPC Service

```
Create an RPC service called "authservice" with methods Login and Verify
```

### Analyzing a Project

```
Analyze the project in /path/to/my/go-zero/project
```

You'll get:

- List of all endpoints
- Service dependencies
- Configuration files
- Suggestions for improvements

### Generating Database Models

```
Generate models for the users table in mysql://user:pass@localhost:3306/mydb
```

### Creating Middleware

```
Generate authentication middleware for JWT tokens
```

## Common Use Cases

### 1. Starting a New Microservice Project

```
Create these services in ./services directory:
- API gateway on port 8080
- User service RPC on port 9001
- Order service RPC on port 9002
```

### 2. Spec-First Development

```
Create an API spec for a user service with these endpoints:
- POST /login
- POST /register
- GET /user/:id
- PUT /user/:id
```

Then generate code:

```
Generate code from userservice.api
```

### 3. Adding Features to Existing Service

```
I have a service at ./userservice. Add rate limiting middleware.
```

### 4. Configuration Management

```
Generate a production configuration template for my API service
```

### 5. Migration from Another Framework

```
How do I migrate my Express.js API to go-zero?
```

## Tips & Best Practices

### Service Naming

- ✅ Good: `userservice`, `orderservice`, `apigateway`
- ❌ Bad: `user-service` (hyphens not allowed), `123service` (must start with letter)

### Port Selection

- Use ports 8080-8089 for API services
- Use ports 9000-9099 for RPC services
- mcp-zero will warn if ports are already in use

### Project Organization

For monorepo:

```
myproject/
├── services/
│   ├── api-gateway/
│   ├── user-service/
│   └── order-service/
└── shared/
    └── types/
```

Ask Claude:

```
Create services in ./services directory
```

### Error Recovery

If generation fails:

```
The service generation failed. Try again with corrected parameters.
```

mcp-zero preserves partial state, so you won't lose progress.

### Getting Help

```
Explain go-zero middleware
```

```
What's the difference between API and RPC services in go-zero?
```

```
Show me best practices for error handling in go-zero
```

## Troubleshooting

### "goctl not found"

**Problem**: mcp-zero can't find goctl executable.

**Solution**:

1. Verify goctl is installed: `goctl --version`
2. Add GOCTL_PATH to Claude config (see Installation step 3)
3. Or add goctl to your PATH

### "Service name validation failed"

**Problem**: Service name contains invalid characters.

**Solution**: Use only letters, numbers, and underscores. Start with a letter.

- `user-service` → `userservice` or `user_service`

### "Port already in use"

**Problem**: Specified port is occupied.

**Solution**: Choose a different port or stop the service using that port.

### "Build failed"

**Problem**: Generated code doesn't compile.

**Solution**:

1. Check the error message for missing dependencies
2. Run `go mod tidy`
3. Report the issue if it persists

### "Permission denied"

**Problem**: Can't write to output directory.

**Solution**: Use a directory where you have write permissions, or create the directory first.

## Advanced Usage

### Custom Templates

```
Use this middleware template for my service:
[paste your template]
```

### Batch Operations

```
Create 5 microservices: user, order, product, payment, notification
All RPC services, ports starting from 9001
```

### Configuration Validation

```
Validate my configuration file at ./etc/config.yaml
```

### Project Documentation

```
Generate documentation for my project structure
```

## Environment Variables

- `GOCTL_PATH`: Override goctl executable location
- `MCP_ZERO_DEBUG`: Enable debug logging (development only)

## Next Steps

- Read the [full documentation](https://github.com/zeromicro/mcp-zero)
- Explore [go-zero documentation](https://go-zero.dev/)
- Join the [go-zero community](https://github.com/zeromicro/go-zero/discussions)

## Support

- GitHub Issues: [github.com/zeromicro/mcp-zero/issues](https://github.com/zeromicro/mcp-zero/issues)
- go-zero Discord: [discord.gg/go-zero](https://discord.gg/go-zero)

---

**Happy coding with mcp-zero! 🚀**
