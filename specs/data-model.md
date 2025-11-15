# Data Model

**Feature**: MCP Tool for go-zero Framework
**Date**: November 14, 2025
**Phase**: 1 (Design & Contracts)

## Overview

This document defines the data structures, relationships, and validation rules for the MCP-zero tool system.

## Core Entities

### 1. MCPTool

Represents an individual tool exposed through the MCP protocol.

**Attributes**:

- `name` (string, required): Unique tool identifier (e.g., "create_api_service")
- `description` (string, required): Human-readable tool purpose
- `inputSchema` (JSON Schema, required): Parameter validation schema
- `handler` (function, required): Tool execution function

**Validation Rules**:

- Name must be lowercase with underscores
- Description must be non-empty and descriptive
- Input schema must be valid JSON Schema

**State Transitions**: N/A (stateless)

### 2. ServiceProject

Represents a generated go-zero service (API or RPC).

**Attributes**:

- `name` (string, required): Service name (validated Go identifier)
- `type` (enum, required): "api" | "rpc"
- `path` (string, required): Absolute filesystem path
- `port` (integer, optional): Service port number (1024-65535)
- `style` (enum, optional): "go_zero" | "gozero" (default: "go_zero")
- `config` (Config, required): Service configuration
- `status` (enum, required): "generating" | "fixing" | "validating" | "ready" | "failed"

**Validation Rules**:

- Name: No hyphens, valid Go identifier, starts with letter
- Path: Must be absolute, writable directory
- Port: If specified, must be 1024-65535 and not in use
- Type: Must be "api" or "rpc"

**State Transitions**:

```
generating → fixing → validating → ready
         ↓                    ↓        ↓
         └────────── failed ──────────┘
```

**Relationships**:

- Has one Config
- Contains multiple SourceFile entities

### 3. APISpecification

Represents a go-zero API specification file (.api).

**Attributes**:

- `path` (string, required): Absolute path to .api file
- `serviceName` (string, required): Service name from spec
- `endpoints` ([]Endpoint, required): List of API endpoints
- `types` ([]TypeDefinition, required): Request/response types
- `groups` ([]Group, optional): Endpoint groups

**Validation Rules**:

- Path must exist and be readable
- Must parse successfully with go-zero parser
- At least one endpoint required

**Relationships**:

- Has many Endpoint entities
- Has many TypeDefinition entities

### 4. Endpoint

Represents a single API endpoint.

**Attributes**:

- `method` (enum, required): "GET" | "POST" | "PUT" | "DELETE" | "PATCH"
- `path` (string, required): URL path (e.g., "/api/users/:id")
- `handler` (string, required): Handler function name
- `requestType` (string, optional): Request body type name
- `responseType` (string, optional): Response body type name

**Validation Rules**:

- Method must be valid HTTP method
- Path must start with "/"
- Handler must be valid Go identifier

### 5. RPCService

Represents a go-zero RPC service with protobuf definition.

**Attributes**:

- `name` (string, required): Service name
- `path` (string, required): Service directory path
- `protoContent` (string, required): Protobuf definition
- `methods` ([]RPCMethod, required): Service methods

**Validation Rules**:

- Proto content must be valid protobuf syntax
- At least one method required

**Relationships**:

- Has many RPCMethod entities

### 6. RPCMethod

Represents a single RPC method.

**Attributes**:

- `name` (string, required): Method name
- `requestType` (string, required): Request message type
- `responseType` (string, required): Response message type
- `streaming` (enum, required): "none" | "client" | "server" | "bidirectional"

### 7. DatabaseModel

Represents generated database access layer code.

**Attributes**:

- `tableName` (string, required): Database table name
- `modelName` (string, required): Generated struct name
- `fields` ([]ModelField, required): Table columns
- `dbType` (enum, required): "mysql" | "postgresql" | "mongo"
- `connectionInfo` (ConnectionInfo, required): Database credentials

**Validation Rules**:

- Table name must exist in target database
- Connection info must be valid for db type

**Relationships**:

- Has many ModelField entities
- Has one ConnectionInfo

### 8. ModelField

Represents a database column in generated model.

**Attributes**:

- `name` (string, required): Column name
- `goType` (string, required): Go type (string, int64, etc.)
- `dbType` (string, required): Database type (VARCHAR, INT, etc.)
- `nullable` (boolean, required): Whether null values allowed
- `primaryKey` (boolean, required): Is primary key
- `tags` (map[string]string, required): Struct tags (json, db)

### 9. Config

Represents service configuration.

**Attributes**:

- `name` (string, required): Service name
- `host` (string, optional): Listen host (default: "0.0.0.0")
- `port` (integer, required): Listen port
- `timeout` (integer, optional): Request timeout in ms (default: 3000)
- `customFields` (map[string]interface{}, optional): User custom config

**Validation Rules**:

- Port: 1024-65535
- Timeout: >0
- Custom fields: Valid YAML types

### 10. ConnectionInfo

Represents database connection information.

**Attributes**:

- `type` (enum, required): "connection_string" | "credential_file"
- `value` (string, required): Connection string or file path
- `encrypted` (boolean, optional): Whether credentials are encrypted

**Security Rules**:

- Never logged or persisted in generated code
- Cleared from memory after use
- Validated format before use

### 11. ProjectAnalysis

Represents analysis results of existing project.

**Attributes**:

- `projectPath` (string, required): Root project directory
- `serviceType` (enum, required): "api" | "rpc" | "mixed"
- `endpoints` ([]Endpoint, optional): Found API endpoints
- `rpcServices` ([]RPCService, optional): Found RPC services
- `dependencies` ([]Dependency, required): External dependencies
- `configFiles` ([]string, required): Configuration file paths

**Relationships**:

- Has many Endpoint entities (from API specs)
- Has many RPCService entities (from proto files)
- Has many Dependency entities

### 12. Dependency

Represents an external Go module dependency.

**Attributes**:

- `module` (string, required): Module path (e.g., "github.com/zeromicro/go-zero")
- `version` (string, required): Module version
- `usage` ([]string, optional): Where dependency is imported

### 13. Template

Represents a code template.

**Attributes**:

- `category` (enum, required): "middleware" | "error_handler" | "deployment"
- `name` (string, required): Template name
- `content` (string, required): Template content
- `parameters` ([]TemplateParameter, required): Required parameters

**Relationships**:

- Has many TemplateParameter entities

### 14. TemplateParameter

Represents a template parameter.

**Attributes**:

- `name` (string, required): Parameter name
- `type` (string, required): Parameter type (string, int, bool)
- `required` (boolean, required): Is parameter required
- `default` (string, optional): Default value
- `description` (string, required): Parameter description

## Validation Summary

### Cross-Entity Validation

1. **Port Uniqueness**: Within same workspace, no two services can use same port
2. **Name Uniqueness**: Within same workspace, service names must be unique
3. **Path Conflicts**: Services cannot be created in overlapping directories

### Common Validation Patterns

**Service Name**:

```go
func validateServiceName(name string) error {
    if !startsWithLetter(name) {
        return errors.New("service name must start with letter")
    }
    if containsHyphen(name) {
        suggestion := strings.ReplaceAll(name, "-", "_")
        return fmt.Errorf("service name cannot contain hyphens, try: %s", suggestion)
    }
    if !isValidGoIdentifier(name) {
        return errors.New("service name must be valid Go identifier")
    }
    return nil
}
```

**Port Validation**:

```go
func validatePort(port int) error {
    if port < 1024 || port > 65535 {
        return errors.New("port must be between 1024 and 65535")
    }
    if isPortInUse(port) {
        return fmt.Errorf("port %d is already in use", port)
    }
    return nil
}
```

**Path Validation**:

```go
func validatePath(path string) error {
    if !filepath.IsAbs(path) {
        return errors.New("path must be absolute")
    }
    if !isWritable(path) {
        return fmt.Errorf("path %s is not writable", path)
    }
    return nil
}
```

## Entity Lifecycle

### ServiceProject Lifecycle

1. **Creation**: User requests service creation
2. **Validation**: Inputs validated (name, port, path)
3. **Generation**: goctl executes, files created
4. **Fixing**: Import paths fixed, modules initialized
5. **Validation**: Build verification
6. **Ready**: Service ready to run

### Error States

- **Failed Generation**: goctl execution failed
- **Failed Fixing**: Post-generation fixes failed
- **Failed Validation**: Build verification failed

All failures preserve partial state for retry with correction.

## Performance Considerations

### Caching

- **goctl Path**: Cache discovered path for session
- **Project Analysis**: Cache parsed results (TTL: 5 minutes)
- **Validation Results**: Cache port availability checks (TTL: 1 minute)

### Optimization

- Batch file operations
- Parallel validation where possible
- Lazy loading of analysis results

## Security Considerations

### Sensitive Data

**ConnectionInfo**:

- Never included in logs
- Cleared from memory after use
- Not persisted to disk
- Validated before use

**Credentials**:

- Support for encrypted credential files
- Environment variable injection preferred
- Clear security guidelines in documentation
