# MCP Tool Contracts

**Feature**: MCP Tool for go-zero Framework
**Date**: November 14, 2025

## Overview

This document defines the MCP tool contracts (input/output schemas) for all tools exposed by the mcp-zero server.

## Tool 1: create_api_service

**Purpose**: Create a new go-zero API service with proper structure and configuration.

**Input Schema**:

```json
{
  "type": "object",
  "properties": {
    "service_name": {
      "type": "string",
      "description": "Name of the service (must be valid Go identifier, no hyphens)",
      "pattern": "^[a-zA-Z][a-zA-Z0-9_]*$"
    },
    "port": {
      "type": "integer",
      "description": "Port number for the service (1024-65535)",
      "minimum": 1024,
      "maximum": 65535
    },
    "output_dir": {
      "type": "string",
      "description": "Output directory path (absolute path recommended)"
    },
    "style": {
      "type": "string",
      "description": "Code style convention",
      "enum": ["go_zero", "gozero"],
      "default": "go_zero"
    }
  },
  "required": ["service_name"]
}
```

**Output**: Success message with service path and next steps, or detailed error.

**Example Request**:

```json
{
  "service_name": "userservice",
  "port": 8080,
  "output_dir": "/Users/kevin/projects/myapp",
  "style": "go_zero"
}
```

**Example Success Response**:

```text
✅ API service 'userservice' created successfully!

Location: /Users/kevin/projects/myapp/userservice
Port: 8080
Entry point: userservice.go

Next steps:
1. cd /Users/kevin/projects/myapp/userservice
2. go run userservice.go
3. Test: curl http://localhost:8080/ping

Service is ready to run!
```

## Tool 2: create_rpc_service

**Purpose**: Create a new go-zero RPC service with protobuf definitions.

**Input Schema**:

```json
{
  "type": "object",
  "properties": {
    "service_name": {
      "type": "string",
      "description": "Name of the RPC service",
      "pattern": "^[a-zA-Z][a-zA-Z0-9_]*$"
    },
    "proto_content": {
      "type": "string",
      "description": "Protobuf definition for the service"
    },
    "output_dir": {
      "type": "string",
      "description": "Output directory path"
    },
    "style": {
      "type": "string",
      "enum": ["go_zero", "gozero"],
      "default": "go_zero"
    }
  },
  "required": ["service_name", "proto_content"]
}
```

**Output**: Success message with RPC service details or error.

## Tool 3: generate_api_from_spec

**Purpose**: Generate go-zero code from existing API specification file.

**Input Schema**:

```json
{
  "type": "object",
  "properties": {
    "api_file": {
      "type": "string",
      "description": "Path to .api specification file"
    },
    "output_dir": {
      "type": "string",
      "description": "Output directory for generated code"
    },
    "style": {
      "type": "string",
      "enum": ["go_zero", "gozero"],
      "default": "go_zero"
    }
  },
  "required": ["api_file"]
}
```

**Output**: Generated code location and compilation status.

## Tool 4: generate_model

**Purpose**: Generate database model code from schema or DDL.

**Input Schema**:

```json
{
  "type": "object",
  "properties": {
    "source_type": {
      "type": "string",
      "description": "Source type for model generation",
      "enum": ["mysql", "postgresql", "mongo", "ddl"]
    },
    "source": {
      "type": "string",
      "description": "Database connection string or DDL file path"
    },
    "table": {
      "type": "string",
      "description": "Table name (for database sources)"
    },
    "output_dir": {
      "type": "string",
      "description": "Output directory for model files",
      "default": "./model"
    }
  },
  "required": ["source_type", "source"]
}
```

**Output**: Generated model files location.

## Tool 5: create_api_spec

**Purpose**: Create a properly formatted API specification file.

**Input Schema**:

```json
{
  "type": "object",
  "properties": {
    "service_name": {
      "type": "string",
      "description": "Name of the API service"
    },
    "endpoints_json": {
      "type": "string",
      "description": "JSON string containing array of endpoint objects"
    },
    "output_file": {
      "type": "string",
      "description": "Output file path for the API spec"
    }
  },
  "required": ["service_name", "endpoints_json"]
}
```

**Endpoints JSON Format**:

```json
[
  {
    "method": "POST",
    "path": "/api/user/login",
    "handler": "Login"
  },
  {
    "method": "GET",
    "path": "/api/user/:id",
    "handler": "GetUser"
  }
]
```

## Tool 6: analyze_project

**Purpose**: Analyze existing go-zero project structure.

**Input Schema**:

```json
{
  "type": "object",
  "properties": {
    "project_path": {
      "type": "string",
      "description": "Path to project directory"
    }
  },
  "required": ["project_path"]
}
```

**Output**: Structured analysis results including endpoints, services, dependencies.

**Example Output**:

```text
📊 Project Analysis: /path/to/project

Service Type: API
Endpoints Found: 12

API Endpoints:
- POST /api/user/login → LoginHandler
- GET /api/user/:id → GetUserHandler
- PUT /api/user/:id → UpdateUserHandler
...

Dependencies:
- github.com/zeromicro/go-zero v1.5.0
- github.com/go-sql-driver/mysql v1.7.1
...

Configuration Files:
- etc/user-api.yaml
```

## Tool 7: validate_config

**Purpose**: Validate service configuration files.

**Input Schema**:

```json
{
  "type": "object",
  "properties": {
    "config_file": {
      "type": "string",
      "description": "Path to configuration file"
    }
  },
  "required": ["config_file"]
}
```

**Output**: Validation results with issues and suggestions.

## Tool 8: generate_config_template

**Purpose**: Generate configuration template for deployment environment.

**Input Schema**:

```json
{
  "type": "object",
  "properties": {
    "environment": {
      "type": "string",
      "description": "Deployment environment",
      "enum": ["development", "staging", "production"]
    },
    "service_type": {
      "type": "string",
      "enum": ["api", "rpc"]
    },
    "output_file": {
      "type": "string",
      "description": "Output file path"
    }
  },
  "required": ["environment", "service_type"]
}
```

## Tool 9: generate_template

**Purpose**: Generate common code templates (middleware, error handlers, etc.).

**Input Schema**:

```json
{
  "type": "object",
  "properties": {
    "template_type": {
      "type": "string",
      "description": "Type of template to generate",
      "enum": ["middleware", "error_handler", "deployment"]
    },
    "template_name": {
      "type": "string",
      "description": "Specific template name (e.g., 'auth_middleware', 'k8s_deployment')"
    },
    "parameters": {
      "type": "object",
      "description": "Template-specific parameters"
    },
    "output_file": {
      "type": "string",
      "description": "Output file path"
    }
  },
  "required": ["template_type", "template_name"]
}
```

## Tool 10: query_docs

**Purpose**: Query go-zero documentation and get explanations.

**Input Schema**:

```json
{
  "type": "object",
  "properties": {
    "query": {
      "type": "string",
      "description": "Documentation query or concept to explain"
    },
    "query_type": {
      "type": "string",
      "description": "Type of documentation query",
      "enum": ["concept", "migration", "troubleshooting", "general"]
    }
  },
  "required": ["query"]
}
```

**Output**: Relevant documentation with examples.

## Error Response Format

All tools follow consistent error format:

```json
{
  "error": true,
  "message": "Human-readable error description",
  "details": {
    "command": "goctl api new userservice",
    "output": "Error: invalid service name",
    "suggestion": "Service names cannot contain hyphens. Try: userservice or user_service"
  }
}
```

## Success Response Format

All tools follow consistent success format:

```text
✅ [Operation] completed successfully!

[Key details]

Next steps:
1. [Step 1]
2. [Step 2]
...

[Additional context or tips]
```
