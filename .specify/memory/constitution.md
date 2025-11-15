<!--
Sync Impact Report:
- Version change: Template → 1.0.0 (Initial constitution)
- Modified principles: None (initial creation)
- Added sections:
  * Core Principles (6 principles: Environment Resilience, Complete Automation, Developer Experience First, Validation & Safety, Tool Composability, Clear Error Communication)
  * Tool Design Requirements
  * Quality Standards
  * Governance
- Removed sections: None
- Templates requiring updates:
  ✅ .specify/templates/plan-template.md (verified alignment)
  ✅ .specify/templates/spec-template.md (verified alignment)
  ✅ .specify/templates/tasks-template.md (verified alignment)
- Follow-up TODOs: None
-->

# MCP-Zero Constitution

## Core Principles

### I. Environment Resilience

**Tools MUST work reliably across diverse runtime environments.**

The MCP server operates in constrained environments (limited PATH, different working directories, variable permissions). Every tool MUST:

- Discover dependencies (goctl, go, etc.) through multiple fallback strategies
- Use absolute paths for all external command execution
- Validate availability of required tools before execution
- Provide actionable error messages when dependencies are missing, including installation instructions

**Rationale**: MCP servers run in isolated environments where typical shell configurations don't apply. Tools that assume specific PATH values or relative paths will fail unpredictably for users.

### II. Complete Automation

**Every tool MUST deliver a fully working, ready-to-run artifact.**

When a tool creates or generates code, it MUST:

- Fix all import paths to use local module names (no absolute remote paths)
- Initialize Go modules automatically
- Run `go mod tidy` to resolve dependencies
- Verify the generated code compiles successfully
- Apply any framework-specific conventions automatically

**Rationale**: Users invoke MCP tools expecting complete solutions. Requiring manual post-generation fixes creates friction and diminishes the value proposition of automation.

### III. Developer Experience First

**Tools MUST anticipate and prevent common pitfalls.**

Every tool MUST:

- Validate inputs before execution (service names, file paths, parameter formats)
- Provide clear examples in error messages
- Suggest corrections for common mistakes (e.g., "user-service" → "userservice")
- Generate comprehensive success messages with next steps
- Configure sensible defaults that work out-of-the-box

**Rationale**: The target audience includes developers unfamiliar with go-zero conventions. Proactive validation and helpful guidance accelerates learning and reduces support burden.

### IV. Validation & Safety

**Tools MUST verify their own outputs and fail fast with clear diagnostics.**

Every tool MUST:

- Check that generated files exist and contain expected content
- Run build verification before reporting success
- Validate configuration files are syntactically correct
- Provide detailed error context (command run, output, failure reason)
- Never leave partially-completed artifacts without warning

**Rationale**: Silent failures or incomplete generations erode trust. Verification catches issues immediately when context is fresh, rather than forcing users to debug later.

### V. Tool Composability

**Tools MUST be designed to work independently and in combination.**

Every tool MUST:

- Accept both minimal and fully-specified parameters
- Support customizable output directories
- Use consistent parameter naming conventions
- Not depend on previous tool invocations (stateless)
- Document all parameters with examples

**Rationale**: Users may want to scaffold a complete microservice architecture or generate individual components. Composable tools enable both workflows without special-casing.

### VI. Clear Error Communication

**Error messages MUST be actionable, not just informative.**

Every error message MUST:

- State what failed and why in plain language
- Provide the specific command or operation that failed
- Include relevant output or logs
- Suggest concrete remediation steps
- Reference documentation URLs when available

**Rationale**: Generic errors ("command failed") waste user time. Actionable errors enable self-service problem resolution and reduce frustration.

## Tool Design Requirements

### Service Generation Tools

Tools that generate services (API, RPC) MUST:

- Validate service names against go-zero constraints (no hyphens, valid identifiers)
- Support custom ports/configurations with automatic config file updates
- Generate project structure matching go-zero conventions
- Include example handlers with working business logic
- Document generated endpoints and their usage

### Code Generation Tools

Tools that generate code from specs MUST:

- Validate source files exist and are readable
- Support multiple code generation styles (go_zero, gozero)
- Preserve custom modifications in existing files when possible
- Generate idiomatic Go code following community standards
- Include comprehensive inline documentation

### Model Generation Tools

Tools that generate database models MUST:

- Support major databases (MySQL, PostgreSQL, MongoDB)
- Accept both connection strings and DDL files
- Generate type-safe model interfaces
- Include CRUD operations with proper error handling
- Support custom table selection for large schemas

## Quality Standards

### Testing Requirements

- All tool functions MUST have unit tests covering success and failure paths
- Integration tests MUST verify end-to-end workflows with real go-zero projects
- Test suites MUST run in CI/CD without external dependencies where possible
- Edge cases identified in user testing MUST be added to regression tests

### Documentation Requirements

- Every tool MUST have usage examples in README
- Parameter descriptions MUST be clear and unambiguous
- Common error scenarios MUST be documented with solutions
- Architecture decisions MUST be recorded in ADR format

### Performance Standards

- Tool invocations MUST complete within 30 seconds for typical projects
- Dependency resolution MUST use caching where applicable
- Progress feedback MUST be provided for operations exceeding 5 seconds

## Governance

This constitution defines the non-negotiable principles for mcp-zero development. All contributions MUST align with these principles.

**Amendment Process**:

1. Proposed changes MUST be discussed with rationale and impact analysis
2. Amendments require majority approval from project maintainers
3. Version bumps follow semantic versioning:
   - MAJOR: Principle removal/redefinition or backward-incompatible changes
   - MINOR: New principles or substantial expansions
   - PATCH: Clarifications, corrections, non-semantic improvements

**Compliance Verification**:

- All PRs MUST demonstrate adherence to applicable principles
- Code reviews MUST explicitly verify environment resilience and automation completeness
- User-reported issues violating principles MUST be prioritized

**Living Document**: This constitution evolves with the project. Principles proven impractical MUST be amended or removed rather than ignored.

**Version**: 1.0.0 | **Ratified**: 2025-11-14 | **Last Amended**: 2025-11-14
