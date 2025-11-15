# Feature Specification: MCP Tool for go-zero Framework

**Feature Branch**: `001-mcp-tool-go-zero`
**Created**: November 14, 2025
**Status**: Draft
**Input**: User description: "write a mcp tool to help users to use go-zero"

## Clarifications

### Session 2025-11-14

- Q: How will developers interact with this tool? → A: AI assistant integration (works through conversational interface like Claude Desktop via MCP protocol)
- Q: Should the tool support creating multiple services within a single workspace/repository? → A: Multiple services per workspace (monorepo approach, services in subdirectories)
- Q: When code generation fails partway through, how should the tool handle recovery? → A: Retry with correction (keep partial state, allow user to fix inputs and continue)
- Q: For configuration validation (User Story 7), how strict should validation be? → A: Schema-based (validate against defined schema but allow additional custom fields)
- Q: For database model generation (User Story 4), how should database credentials be provided? → A: Both connection string and credential file support (flexible, supports secure credential management)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Create New API Service (Priority: P1)

A developer wants to quickly create a new REST API service with proper structure and configuration without needing to know complex setup procedures or remember specific commands.

**Why this priority**: This is the most fundamental operation - creating a new service is the entry point for any project. Without this capability, developers cannot start building.

**Independent Test**: Can be fully tested by requesting "create a new API service called 'user-service' on port 8080" and verifying that a complete, runnable service project is generated with correct structure, configuration, and all dependencies resolved.

**Acceptance Scenarios**:

1. **Given** a developer needs a new API service, **When** they request "create API service named 'user-service' on port 8080", **Then** a complete service project is generated with correct directory structure, proper initialization, dependency resolution, and can be immediately run
2. **Given** a service name with invalid format, **When** they request "create API service named 'user-service'", **Then** the system validates the name and provides clear guidance on naming requirements
3. **Given** required tools are not easily accessible, **When** they request service creation, **Then** the tool automatically locates required utilities without requiring manual configuration

---

### User Story 2 - Generate Code from API Specification (Priority: P2)

A developer has written an API specification document and wants to generate the corresponding service implementation structure (request handlers, business logic stubs, data types) without manual coding.

**Why this priority**: This enables specification-first development workflow, which is a best practice for API design. It's the second most common workflow after initial service creation.

**Independent Test**: Can be fully tested by providing a valid API specification document and verifying that complete handler structure, logic stubs, and type definitions are generated and can be immediately run.

**Acceptance Scenarios**:

1. **Given** an existing API specification document, **When** developer requests "generate code from this API spec", **Then** complete handler structure, logic stubs, and type definitions are generated with correct dependencies and run successfully
2. **Given** a specification with multiple endpoints, **When** code generation is requested, **Then** all handlers are generated with proper routing and the project structure is maintained
3. **Given** invalid specification syntax, **When** generation is attempted, **Then** clear error messages indicate what's wrong in the specification

---

### User Story 3 - Create RPC Service (Priority: P3)

A developer needs to create a new microservice with defined remote procedure call interfaces for service-to-service communication.

**Why this priority**: Essential for microservice architectures but less common than REST API services. Can be implemented after API service creation is stable.

**Independent Test**: Can be fully tested by requesting "create RPC service named 'auth-service'" with interface definition and verifying that complete service structure is generated and runs successfully.

**Acceptance Scenarios**:

1. **Given** a developer needs microservice with RPC interfaces, **When** they request "create RPC service named 'auth-service' with interface definition", **Then** complete service structure is generated with interface files, service implementation, and configuration
2. **Given** existing interface definition, **When** RPC service generation is requested, **Then** the interface definition is used to generate service code and dependencies are properly configured

---

### User Story 4 - Generate Database Models (Priority: P4)

A developer has an existing database or table definition and wants to generate data access layer code to interact with the database.

**Why this priority**: Important for data-driven applications but depends on having database schema first. Can be implemented after core service generation features.

**Independent Test**: Can be fully tested by providing database connection details or table definition and verifying that data access code is generated with proper create, read, update, delete operations and can interact with the database.

**Acceptance Scenarios**:

1. **Given** database credentials (via connection string or credential file) and table name, **When** developer requests "generate model from database", **Then** complete data access code is generated with CRUD operations for the specified table
2. **Given** a table definition file, **When** model generation is requested, **Then** data access code is generated for all tables defined
3. **Given** different database systems, **When** model generation is requested, **Then** appropriate data access code is generated for the specific database type

---

### User Story 5 - Create Custom API Specification (Priority: P5)

A developer wants to create a properly formatted API specification document as a starting point for their API design.

**Why this priority**: Helpful for new users but experienced developers often write specifications manually. Can be implemented last as it's a convenience feature.

**Independent Test**: Can be fully tested by requesting "create API spec for user service with login and register endpoints" and verifying that a valid specification document is generated that can be used for code generation.

**Acceptance Scenarios**:

1. **Given** service name and endpoint descriptions, **When** developer requests "create API spec", **Then** a properly formatted specification document is generated following framework conventions
2. **Given** multiple endpoints with different request methods, **When** spec creation is requested, **Then** all endpoints are included with correct syntax in the generated specification

---

### User Story 6 - Analyze Existing Project (Priority: P6)

A developer working on an existing service project wants to understand the project structure, available endpoints, service definitions, and dependencies without manually reading through multiple files.

**Why this priority**: Valuable for onboarding new team members and understanding legacy projects, but not critical for new project creation workflows.

**Independent Test**: Can be fully tested by pointing the tool at an existing service project and verifying that it correctly identifies and lists all endpoints, service definitions, interface contracts, and project dependencies.

**Acceptance Scenarios**:

1. **Given** an existing service project, **When** developer requests "analyze this project", **Then** system lists all available endpoints, their request/response structures, and routing configuration
2. **Given** a project with multiple service definitions, **When** analysis is requested, **Then** all service interfaces and their methods are identified and documented
3. **Given** a project with dependencies, **When** analysis is requested, **Then** all external dependencies and their versions are listed with usage context

---

### User Story 7 - Manage Service Configuration (Priority: P7)

A developer needs to validate existing service configuration or generate configuration templates for different deployment scenarios (development, staging, production).

**Why this priority**: Important for deployment workflows but not blocking for development. Can be addressed after core generation features are stable.

**Independent Test**: Can be fully tested by requesting "validate my service configuration" or "generate production config template" and verifying that validation feedback is accurate or generated templates follow environment-specific best practices.

**Acceptance Scenarios**:

1. **Given** an existing service configuration, **When** developer requests "validate configuration", **Then** system checks configuration completeness against schema, allows custom fields, identifies potential issues, and suggests improvements
2. **Given** a deployment environment type, **When** developer requests "generate config template for production", **Then** appropriate configuration template is generated with environment-specific settings and security best practices
3. **Given** multiple configuration files, **When** validation is requested, **Then** system ensures consistency across configuration files and identifies conflicts

---

### User Story 8 - Generate Common Code Templates (Priority: P8)

A developer needs to add common patterns like middleware, error handlers, or deployment configurations to their service and wants pre-built templates that follow best practices.

**Why this priority**: Accelerates development but developers can write these manually. Lower priority than core generation features.

**Independent Test**: Can be fully tested by requesting "generate middleware for authentication" or "create deployment configuration" and verifying that generated templates are complete, follow best practices, and integrate properly with existing code.

**Acceptance Scenarios**:

1. **Given** a need for request middleware, **When** developer requests "generate authentication middleware", **Then** complete middleware template is generated with proper integration points and usage examples
2. **Given** a need for error handling, **When** developer requests "generate error handler", **Then** custom error handler template is generated following framework conventions
3. **Given** deployment requirements, **When** developer requests "generate deployment configuration", **Then** appropriate deployment templates are generated for the specified platform

---

### User Story 9 - Query Framework Documentation (Priority: P9)

A developer needs to understand framework concepts, find migration guidance, or get help with specific features without leaving their development environment.

**Why this priority**: Very helpful for learning and troubleshooting but not blocking for actual development work. Lowest priority as developers can access documentation separately.

**Independent Test**: Can be fully tested by asking "explain middleware in go-zero" or "how do I migrate from Express to go-zero" and verifying that responses are accurate, relevant, and include practical examples.

**Acceptance Scenarios**:

1. **Given** a framework concept question, **When** developer asks "explain how routing works", **Then** clear explanation is provided with practical examples and common use cases
2. **Given** a migration scenario, **When** developer asks "migrate from another framework", **Then** step-by-step migration guide is provided with code comparisons and best practices
3. **Given** a troubleshooting question, **When** developer asks about specific error or pattern, **Then** relevant documentation sections and solutions are provided

---

### Edge Cases

- What happens when required code generation tools are not installed or not accessible? System should automatically search standard installation locations and user-configured paths, providing clear error message with installation instructions if tools cannot be found.
- What happens when service name contains invalid characters or format? System should validate service names against framework requirements and provide clear error messages explaining naming conventions and requirements.
- What happens when target directory already exists? System should detect existing directories and either warn user about potential conflicts or suggest alternative locations.
- What happens when generated project has module reference issues? System should automatically correct module references to ensure proper dependency resolution.
- What happens when port number is invalid or already in use? System should validate port numbers and provide guidance on choosing appropriate available ports.
- What happens when project lacks proper initialization? System should automatically perform all necessary initialization steps to ensure project is ready to run.
- What happens when project dependencies are missing? System should automatically identify and install all required dependencies without user intervention.
- What happens when generated project is incomplete or has errors? System should verify project completeness and readiness before reporting success to user.
- What happens when creating multiple services in the same workspace? System should organize services in separate subdirectories with appropriate naming and ensure no conflicts between service configurations (e.g., port numbers).
- What happens when generation fails after partial completion? System should preserve successfully generated files, clearly report which step failed, and allow user to provide corrected inputs to continue from the failure point without regenerating successful parts.
- What happens with database credentials during model generation? System should support both connection strings (for quick development) and credential files (for secure production use), and should never log or persist credentials in generated code or tool output.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST create new REST API service projects with specified name and port configuration in user-specified workspace locations
- **FR-002**: System MUST validate service names against framework naming constraints and provide clear guidance when validation fails
- **FR-001a**: System MUST support creating multiple services within a single workspace using subdirectory structure
- **FR-003**: System MUST locate required code generation tools automatically without user configuration
- **FR-004**: System MUST generate service code from API specification documents
- **FR-005**: System MUST create new microservice RPC projects with interface definitions
- **FR-006**: System MUST generate data access layer code from database schemas or table definitions, accepting credentials via connection strings or secure credential files
- **FR-007**: System MUST create properly formatted API specification documents from endpoint descriptions
- **FR-008**: System MUST ensure generated projects have correct module references and dependencies
- **FR-009**: System MUST configure generated services with user-specified port numbers
- **FR-010**: System MUST ensure generated projects are properly initialized with dependency management
- **FR-011**: System MUST resolve and download all required dependencies automatically
- **FR-012**: System MUST verify generated projects are complete and ready to run before reporting success
- **FR-013**: System MUST provide clear error messages with installation guidance when required tools are missing
- **FR-014**: System MUST provide clear error messages explaining naming requirements when invalid names are provided
- **FR-015**: System MUST translate all technical errors into user-friendly messages with actionable guidance
- **FR-016**: System MUST support configurable naming conventions for generated code
- **FR-017**: System MUST create production-ready project structure with all necessary configuration and scaffolding
- **FR-018**: System MUST analyze existing service projects and extract structural information including endpoints, service definitions, and dependencies
- **FR-019**: System MUST parse and understand API specification documents and interface definition files
- **FR-020**: System MUST validate service configuration files against defined schema while allowing custom fields, and identify potential issues or inconsistencies
- **FR-021**: System MUST generate configuration templates appropriate for different deployment environments
- **FR-022**: System MUST provide common code templates including middleware, error handlers, and deployment configurations
- **FR-023**: System MUST ensure generated templates integrate properly with existing project structure
- **FR-024**: System MUST provide access to framework documentation and concept explanations
- **FR-025**: System MUST offer migration guidance from other frameworks with practical examples
- **FR-026**: System MUST operate through AI assistant conversational interface using Model Context Protocol (MCP)
- **FR-027**: System MUST accept natural language requests and translate them into appropriate code generation actions
- **FR-028**: System MUST preserve partial generation state when errors occur and allow users to correct inputs and retry without losing progress
- **FR-029**: System MUST provide clear indication of which generation steps succeeded and which failed during partial failures

### Key Entities

- **API Service**: A REST API service project containing request handlers, business logic, routing definitions, and runtime configuration
- **RPC Service**: A microservice project with defined interfaces for remote procedure calls and service implementation
- **API Specification**: A document defining REST endpoints, request/response structures, and service behavior
- **Model**: Database access layer providing create, read, update, and delete operations for data entities
- **Configuration**: Service settings including network ports, timeouts, database connections, and operational parameters
- **Template**: Pre-built code pattern for common functionality like middleware, error handling, or deployment configuration
- **Project Analysis**: Extracted structural information about an existing project including endpoints, services, dependencies, and configuration

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Developers can create a new API service and have it responding to requests in under 2 minutes from initial request
- **SC-002**: 100% of generated services are immediately runnable without requiring manual fixes or modifications
- **SC-003**: The tool successfully locates required code generation utilities in 95% of standard installation scenarios without user configuration
- **SC-004**: Service name validation catches 100% of invalid names before attempting project creation
- **SC-005**: Developers can generate complete service code from API specifications in under 30 seconds
- **SC-006**: All generated projects follow framework best practices and conventions automatically
- **SC-007**: Error messages provide actionable solutions in 100% of common failure scenarios (missing tools, invalid names, missing dependencies)
- **SC-008**: 90% of developers can use the tool successfully without consulting documentation (through intuitive interaction and clear error messages)
- **SC-009**: Developers can understand an existing project's structure and endpoints in under 1 minute through project analysis
- **SC-010**: Configuration validation identifies 95% of common configuration issues before deployment
- **SC-011**: Generated templates integrate with existing projects without requiring manual modifications in 90% of cases
- **SC-012**: Framework documentation queries return relevant, accurate responses in under 5 seconds
