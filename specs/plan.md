# Implementation Plan: MCP Tool for go-zero Framework

**Branch**: `001-mcp-tool-go-zero` | **Date**: November 14, 2025 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-mcp-tool-go-zero/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Build an MCP (Model Context Protocol) server that enables developers to use the go-zero framework through natural language conversations with AI assistants. The tool provides comprehensive code generation (API services, RPC services, database models, specifications), project analysis, configuration management, code templates, and documentation queries. All operations work through conversational interface, support monorepo architecture with multiple services, and follow go-zero best practices with automatic validation and error recovery.

## Technical Context

**Language/Version**: Go 1.19+ (required by go-zero framework and MCP SDK compatibility)

**Primary Dependencies**:

- mark3labs/mcp-go v0.37.0+ (Model Context Protocol SDK)
- zeromicro/go-zero (framework tools discovery and validation)
- goctl (go-zero CLI tool - external dependency, must be discovered at runtime)

**Storage**: N/A (tool generates code, doesn't persist state)

**Testing**: Go testing package (unit tests), integration tests with real goctl execution

**Target Platform**: macOS, Linux (MCP servers run as local processes invoked by AI assistants)

**Project Type**: Single binary (MCP server) with tool handlers

**Performance Goals**:

- Tool response time <5 seconds for code generation
- File system operations complete in <30 seconds
- Analysis of existing projects <1 minute

**Constraints**:

- Isolated execution environment (limited PATH, no shell configuration)
- Must discover goctl through multiple fallback strategies
- Cannot assume user's working directory
- Must work with Claude Desktop and other MCP-compatible AI assistants

**Scale/Scope**:

- 9 major tool categories (service creation, code generation, analysis, configuration, templates, documentation)
- Support for go-zero API services, RPC services, models, and specifications
- Monorepo support with multiple services in subdirectories

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Environment Resilience ✅

- **Requirement**: Tools MUST work reliably across diverse runtime environments
- **Status**: PASS - Design includes:
  - Multi-strategy goctl discovery (GOCTL_PATH env var, standard paths, fallbacks)
  - Absolute path usage for all external commands
  - Tool availability validation before execution
  - Actionable error messages with installation instructions

### Complete Automation ✅

- **Requirement**: Every tool MUST deliver a fully working, ready-to-run artifact
- **Status**: PASS - Design includes:
  - Automatic import path fixing
  - Go module initialization (`go mod init`)
  - Dependency resolution (`go mod tidy`)
  - Build verification before reporting success
  - Framework convention application

### Developer Experience First ✅

- **Requirement**: Tools MUST anticipate and prevent common pitfalls
- **Status**: PASS - Design includes:
  - Input validation (service names, paths, parameters)
  - Clear error examples and suggestions
  - Common mistake corrections (e.g., hyphen removal)
  - Comprehensive success messages with next steps
  - Sensible defaults throughout

### Validation & Safety ✅

- **Requirement**: Tools MUST verify their own outputs and fail fast
- **Status**: PASS - Design includes:
  - File existence and content verification
  - Build verification before success
  - Configuration syntax validation
  - Detailed error context (command, output, reason)
  - No silent partial failures

### Tool Composability ✅

- **Requirement**: Tools MUST be designed to work independently and in combination
- **Status**: PASS - Design includes:
  - Minimal and fully-specified parameter support
  - Customizable output directories
  - Consistent parameter naming
  - Stateless operation (no dependency on previous invocations)
  - Comprehensive parameter documentation

### Clear Error Communication ✅

- **Requirement**: Error messages MUST be actionable, not just informative
- **Status**: PASS - Design includes:
  - Plain language failure explanations
  - Specific command/operation that failed
  - Relevant output and logs
  - Concrete remediation steps
  - Documentation URL references

**Overall Gate Status**: ✅ PASS - All constitution requirements addressed in design

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
mcp-zero/                       # Repository root
├── main.go                     # MCP server entry point, tool registration
├── go.mod                      # Go module definition
├── go.sum                      # Dependency lock file
│
├── tools/                      # MCP tool handlers
│   ├── create_api_service.go   # User Story 1: API service creation
│   ├── create_rpc_service.go   # User Story 3: RPC service creation
│   ├── generate_from_spec.go   # User Story 2: Code from API spec
│   ├── generate_model.go       # User Story 4: Database models
│   ├── create_api_spec.go      # User Story 5: API spec creation
│   ├── analyze_project.go      # User Story 6: Project analysis
│   ├── manage_config.go        # User Story 7: Configuration management
│   ├── generate_template.go    # User Story 8: Code templates
│   └── query_docs.go           # User Story 9: Documentation query
│
├── internal/                   # Internal packages (not for external import)
│   ├── goctl/                  # goctl integration
│   │   ├── discovery.go        # Multi-strategy goctl path discovery
│   │   ├── executor.go         # Safe command execution with validation
│   │   └── validator.go        # Output validation
│   │
│   ├── validation/             # Input validation
│   │   ├── service_name.go     # Service name validation (no hyphens, etc.)
│   │   ├── port.go             # Port number validation
│   │   └── path.go             # File/directory path validation
│   │
│   ├── fixer/                  # Post-generation fixers
│   │   ├── imports.go          # Import path fixing
│   │   ├── config.go           # Configuration file updates
│   │   └── modules.go          # Go module initialization
│   │
│   ├── analyzer/               # Project analysis
│   │   ├── api_parser.go       # Parse .api files
│   │   ├── proto_parser.go     # Parse .proto files
│   │   └── project_scanner.go  # Scan project structure
│   │
│   ├── templates/              # Code template generation
│   │   ├── middleware.go       # Middleware templates
│   │   ├── error_handler.go    # Error handler templates
│   │   └── deployment.go       # Deployment config templates
│   │
│   └── docs/                   # Documentation query
│       ├── concepts.go         # Framework concept explanations
│       └── migration.go        # Migration guidance
│
├── pkg/                        # Public packages (if needed for future extensibility)
│
└── tests/                      # Test files
    ├── integration/            # Integration tests with real goctl
    │   ├── api_service_test.go
    │   ├── rpc_service_test.go
    │   └── model_gen_test.go
    │
    └── unit/                   # Unit tests
        ├── validation_test.go
        ├── discovery_test.go
        └── fixer_test.go
```

**Structure Decision**: Single Go binary (MCP server) with organized internal packages. This structure supports:

- Clear separation of tool handlers (one file per user story)
- Reusable internal packages for common operations (validation, fixing, analysis)
- Constitution compliance through dedicated modules (discovery, validation, fixing)
- Testability with separate unit and integration test directories
- Future extensibility via pkg/ directory if needed

## Complexity Tracking

**Constitutional Violations**: None identified

All constitution principles are satisfied in the current design:

- Environment Resilience: Multi-strategy goctl discovery addresses diverse runtime environments
- Complete Automation: Generate→Fix→Validate→Report workflow ensures working artifacts
- Developer Experience First: Proactive validation and correction of common mistakes
- Validation & Safety: Build verification and fail-fast on errors
- Tool Composability: Stateless tools with minimal, well-specified parameters
- Clear Error Communication: Actionable error messages with remediation steps

**Design Trade-offs**:

1. **goctl External Dependency**: Must discover and validate goctl at runtime (cannot bundle). Mitigated by multi-level discovery strategy and clear installation guidance.

2. **No State Persistence**: Tool operates statelessly (each invocation is independent). This simplifies reliability but means no cross-tool optimization. Acceptable trade-off for MCP architecture.

3. **Synchronous Operation**: All operations block until complete. For long-running operations (project analysis), provide progress indicators through MCP protocol.

**Complexity Estimate**: Medium

- 9 MCP tools with distinct logic
- External command execution with validation
- Multi-strategy discovery and fallback handling
- Schema-based configuration validation
- Template system with parameter injection

## Phase 0: Research

✅ **Completed** - See [research.md](./research.md) for full technical decisions

**Summary**: 10 technical decisions documented covering MCP SDK selection, goctl discovery strategy, code generation workflow, validation approaches, monorepo support, error recovery, configuration management, credentials handling, project analysis, and template system.

## Phase 1: Design & Contracts

✅ **Completed** - See artifacts below

### Data Model

✅ **Completed** - See [data-model.md](./data-model.md)

**Summary**: 14 entities defined including MCPTool, ServiceProject (with state machine), APISpecification, Endpoint, RPCService, RPCMethod, DatabaseModel, ModelField, Config, ConnectionInfo, ProjectAnalysis, Dependency, Template, TemplateParameter. All include validation rules, relationships, and lifecycle management.

### Tool Contracts

✅ **Completed** - See [contracts/mcp-tools.md](./contracts/mcp-tools.md)

**Summary**: 10 MCP tool contracts defined with JSON schemas for input validation and consistent error/success formats:

1. create_api_service - Create new go-zero API service
2. create_rpc_service - Create new go-zero RPC service
3. generate_api_from_spec - Generate code from .api specification
4. generate_model - Generate database models from schema/DDL
5. create_api_spec - Create sample .api specification file
6. analyze_project - Analyze existing go-zero project structure
7. validate_config - Validate go-zero configuration file
8. generate_config_template - Generate configuration template
9. generate_template - Generate code templates (middleware, handlers, etc.)
10. query_docs - Query go-zero documentation and best practices

### Quick Start Guide

✅ **Completed** - See [quickstart.md](./quickstart.md)

**Summary**: Developer guide covering installation (Go, goctl, mcp-zero), Claude Desktop configuration, tutorial (first API service, running service, adding endpoints), common use cases (microservices, spec-first development), tips & best practices, troubleshooting, and advanced usage.

## Phase 1 Complete

All planning artifacts have been created:

- ✅ research.md - 10 technical decisions
- ✅ data-model.md - 14 entity definitions
- ✅ contracts/mcp-tools.md - 10 tool contracts with JSON schemas
- ✅ quickstart.md - Developer quick start guide
- ✅ Agent context updated - GitHub Copilot now aware of Go 1.19+ and MCP architecture

**Next Step**: Run `/speckit.tasks` command to break down this plan into specific implementation tasks with acceptance criteria.

---

*This plan was generated by the `/speckit.plan` command. For task breakdown, use `/speckit.tasks`.*
