# Tasks: MCP Tool for go-zero Framework

**Input**: Design documents from `/specs/001-mcp-tool-go-zero/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

Based on plan.md structure:
- Main package: `main.go`
- Tool handlers: `tools/`
- Internal packages: `internal/`
- Tests: `tests/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Create project structure per implementation plan in repository root
- [X] T002 Initialize Go module with `go mod init github.com/zeromicro/mcp-zero` in go.mod
- [X] T003 [P] Install MCP SDK dependency (using modelcontextprotocol/go-sdk v1.1.0) in go.mod
- [X] T004 [P] Create main.go with MCP server initialization and stdio transport
- [X] T005 [P] Create tools/ directory for MCP tool handlers
- [X] T006 [P] Create internal/goctl/ directory for goctl integration
- [X] T007 [P] Create internal/validation/ directory for input validation
- [X] T008 [P] Create internal/fixer/ directory for post-generation fixers
- [X] T009 [P] Create internal/analyzer/ directory for project analysis
- [X] T010 [P] Create internal/templates/ directory for code templates
- [X] T011 [P] Create internal/docs/ directory for documentation query
- [X] T012 [P] Create tests/unit/ and tests/integration/ directories

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T013 [P] Implement goctl discovery with multi-level fallback in internal/goctl/discovery.go
- [X] T014 [P] Implement safe command executor with validation in internal/goctl/executor.go
- [X] T015 [P] Implement output validator for generated code in internal/goctl/validator.go
- [X] T016 [P] Implement service name validation in internal/validation/service_name.go
- [X] T017 [P] Implement port validation (1024-65535, not in use) in internal/validation/port.go
- [X] T018 [P] Implement path validation (absolute, writable) in internal/validation/path.go
- [X] T019 [P] Implement import path fixer in internal/fixer/imports.go
- [X] T020 [P] Implement Go module initializer in internal/fixer/modules.go
- [X] T021 [P] Implement configuration file updater in internal/fixer/config.go
- [X] T022 Create common error types and formatting in internal/errors/errors.go
- [X] T023 Create common response formatter in internal/responses/formatter.go

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Create New API Service (Priority: P1) 🎯 MVP

**Goal**: Enable developers to create new REST API service with proper structure through natural language

**Independent Test**: Request "create a new API service called 'userservice' on port 8080" and verify complete, runnable service is generated

### Implementation for User Story 1

- [X] T024 [US1] Define MCPTool struct for create_api_service in tools/create_api_service.go
- [X] T025 [US1] Implement input schema validation with JSON Schema in tools/create_api_service.go
- [X] T026 [US1] Implement service name validation and correction logic in tools/create_api_service.go
- [X] T027 [US1] Implement port validation and availability check in tools/create_api_service.go
- [X] T028 [US1] Implement path validation and directory creation in tools/create_api_service.go
- [X] T029 [US1] Execute goctl command with absolute paths in tools/create_api_service.go
- [X] T030 [US1] Apply import path fixes after generation in tools/create_api_service.go
- [X] T031 [US1] Apply Go module initialization (go mod init, go mod tidy) in tools/create_api_service.go
- [X] T032 [US1] Verify build succeeds (go build) before reporting success in tools/create_api_service.go
- [X] T033 [US1] Generate success response with next steps in tools/create_api_service.go
- [X] T034 [US1] Register create_api_service tool in main.go
- [X] T035 [US1] Create integration test for successful API service creation in tests/integration/api_service_test.go
- [X] T036 [US1] Create integration test for invalid service name handling in tests/integration/api_service_test.go
- [X] T037 [US1] Create integration test for port conflict handling in tests/integration/api_service_test.go

**Checkpoint**: User Story 1 complete - developers can create API services through MCP ✅

---

## Phase 4: User Story 2 - Generate Code from API Specification (Priority: P2)

**Goal**: Enable specification-first development by generating service code from .api files

**Independent Test**: Provide valid API specification and verify complete handler structure is generated and runs

### Implementation for User Story 2

- [X] T038 [P] [US2] Create APISpecification entity parsing in internal/analyzer/api_parser.go
- [X] T039 [P] [US2] Create Endpoint entity extraction in internal/analyzer/api_parser.go
- [X] T040 [US2] Define MCPTool struct for generate_api_from_spec in tools/generate_from_spec.go
- [X] T041 [US2] Implement input schema validation (api_file path) in tools/generate_from_spec.go
- [X] T042 [US2] Validate API specification syntax before generation in tools/generate_from_spec.go
- [X] T043 [US2] Execute goctl api go command with spec file in tools/generate_from_spec.go
- [X] T044 [US2] Apply post-generation fixes (imports, modules) in tools/generate_from_spec.go
- [X] T045 [US2] Verify build succeeds after generation in tools/generate_from_spec.go
- [X] T046 [US2] Generate success response with endpoint list in tools/generate_from_spec.go
- [X] T047 [US2] Register generate_api_from_spec tool in main.go
- [X] T048 [US2] Create integration test for spec-to-code generation in tests/integration/api_service_test.go

**Checkpoint**: User Story 2 complete - specification-first development enabled

---

## Phase 5: User Story 3 - Create RPC Service (Priority: P3)

**Goal**: Enable microservice creation with RPC interfaces for service-to-service communication

**Independent Test**: Request "create RPC service named 'authservice'" with interface and verify service runs

### Implementation for User Story 3

- [X] T049 [P] [US3] Create RPCService entity handling in internal/analyzer/proto_parser.go
- [X] T050 [P] [US3] Create RPCMethod entity extraction in internal/analyzer/proto_parser.go
- [X] T051 [US3] Define MCPTool struct for create_rpc_service in tools/create_rpc_service.go
- [X] T052 [US3] Implement input schema validation (service_name, proto_content) in tools/create_rpc_service.go
- [X] T053 [US3] Validate protobuf syntax before generation in tools/create_rpc_service.go
- [X] T054 [US3] Execute goctl rpc protoc command with proto definition in tools/create_rpc_service.go
- [X] T055 [US3] Apply post-generation fixes for RPC services in tools/create_rpc_service.go
- [X] T056 [US3] Verify RPC service builds successfully in tools/create_rpc_service.go
- [X] T057 [US3] Generate success response with service details in tools/create_rpc_service.go
- [X] T058 [US3] Register create_rpc_service tool in main.go
- [X] T059 [US3] Create integration test for RPC service creation in tests/integration/rpc_service_test.go

**Checkpoint**: User Story 3 complete - RPC microservices can be created

---

## Phase 6: User Story 4 - Generate Database Models (Priority: P4)

**Goal**: Enable data access layer generation from database schemas or table definitions

**Independent Test**: Provide database connection and table name, verify CRUD operations are generated

### Implementation for User Story 4

- [X] T060 [P] [US4] Create DatabaseModel entity structure in internal/analyzer/model_parser.go
- [X] T061 [P] [US4] Create ModelField entity extraction in internal/analyzer/model_parser.go
- [X] T062 [P] [US4] Create ConnectionInfo secure handling in internal/security/credentials.go
- [X] T063 [US4] Define MCPTool struct for generate_model in tools/generate_model.go
- [X] T064 [US4] Implement input schema validation (source_type, source, table) in tools/generate_model.go
- [X] T065 [US4] Support connection string format parsing in tools/generate_model.go
- [X] T066 [US4] Support credential file format parsing in tools/generate_model.go
- [X] T067 [US4] Execute goctl model command with database connection in tools/generate_model.go
- [X] T068 [US4] Clear credentials from memory after use in tools/generate_model.go
- [X] T069 [US4] Verify generated model compiles in tools/generate_model.go
- [X] T070 [US4] Generate success response with model details in tools/generate_model.go
- [X] T071 [US4] Register generate_model tool in main.go
- [X] T072 [US4] Create integration test for model generation in tests/integration/model_gen_test.go

**Checkpoint**: User Story 4 complete - database models can be generated

---

## Phase 7: User Story 5 - Create Custom API Specification (Priority: P5)

**Goal**: Help developers create properly formatted API specification documents

**Independent Test**: Request "create API spec for user service with login endpoint" and verify valid spec is generated

### Implementation for User Story 5

- [X] T073 [P] [US5] Create API spec templates in internal/templates/api_specs.go
- [X] T074 [US5] Define MCPTool struct for create_api_spec in tools/create_api_spec.go
- [X] T075 [US5] Implement input schema validation (service_name, endpoints_json) in tools/create_api_spec.go
- [X] T076 [US5] Parse endpoint descriptions from JSON input in tools/create_api_spec.go
- [X] T077 [US5] Generate .api file content with correct syntax in tools/create_api_spec.go
- [X] T078 [US5] Write specification file to output path in tools/create_api_spec.go
- [X] T079 [US5] Validate generated spec can be parsed by goctl in tools/create_api_spec.go
- [X] T080 [US5] Generate success response with spec location in tools/create_api_spec.go
- [X] T081 [US5] Register create_api_spec tool in main.go
- [X] T082 [US5] Create integration test for spec creation in tests/integration/api_service_test.go

**Checkpoint**: User Story 5 complete - API specifications can be created

---

## Phase 8: User Story 6 - Analyze Existing Project (Priority: P6)

**Goal**: Enable understanding of existing project structure, endpoints, and dependencies

**Independent Test**: Point tool at existing project and verify it identifies all endpoints and services

### Implementation for User Story 6

- [X] T083 [P] [US6] Create ProjectAnalysis entity structure in internal/analyzer/project_scanner.go
- [X] T084 [P] [US6] Create Dependency entity extraction in internal/analyzer/project_scanner.go
- [X] T085 [US6] Implement .api file discovery in internal/analyzer/project_scanner.go
- [X] T086 [US6] Implement .proto file discovery in internal/analyzer/project_scanner.go
- [X] T087 [US6] Implement config file discovery in internal/analyzer/project_scanner.go
- [X] T088 [US6] Implement go.mod dependency parsing in internal/analyzer/project_scanner.go
- [X] T089 [US6] Define MCPTool struct for analyze_project in tools/analyze_project.go
- [X] T090 [US6] Implement input validation (project_path) in tools/analyze_project.go
- [X] T091 [US6] Scan project directory structure in tools/analyze_project.go
- [X] T092 [US6] Parse all API specifications found in tools/analyze_project.go
- [X] T093 [US6] Parse all RPC services found in tools/analyze_project.go
- [X] T094 [US6] Extract all dependencies with versions in tools/analyze_project.go
- [X] T095 [US6] Generate comprehensive analysis report in tools/analyze_project.go
- [X] T096 [US6] Implement result caching (5 min TTL) in tools/analyze_project.go
- [X] T097 [US6] Register analyze_project tool in main.go
- [X] T098 [US6] Create integration test for project analysis in tests/integration/analyze_test.go

**Checkpoint**: User Story 6 complete - existing projects can be analyzed

---

## Phase 9: User Story 7 - Manage Service Configuration (Priority: P7)

**Goal**: Validate and generate service configuration files for different environments

**Independent Test**: Request "validate my config" or "generate production config" and verify accurate results

### Implementation for User Story 7

- [X] T099 [P] [US7] Create Config entity validation in internal/validation/config.go
- [X] T100 [P] [US7] Create config templates for different environments in internal/templates/configs.go
- [X] T101 [US7] Define MCPTool struct for validate_config in tools/manage_config.go
- [X] T102 [US7] Implement config file parsing (YAML/JSON) in tools/manage_config.go
- [X] T103 [US7] Validate config against schema (allow custom fields) in tools/manage_config.go
- [X] T104 [US7] Identify configuration issues and inconsistencies in tools/manage_config.go
- [X] T105 [US7] Generate validation report with suggestions in tools/manage_config.go
- [X] T106 [US7] Define MCPTool struct for generate_config_template in tools/manage_config.go
- [X] T107 [US7] Generate environment-specific config templates in tools/manage_config.go
- [X] T108 [US7] Include security best practices in production configs in tools/manage_config.go
- [X] T109 [US7] Register validate_config and generate_config_template tools in main.go
- [X] T110 [US7] Create integration test for config validation in tests/integration/config_test.go
- [X] T111 [US7] Create integration test for config generation in tests/integration/config_test.go

**Checkpoint**: User Story 7 complete - configuration management enabled

---

## Phase 10: User Story 8 - Generate Common Code Templates (Priority: P8)

**Goal**: Provide pre-built templates for middleware, error handlers, and deployment configs

**Independent Test**: Request "generate authentication middleware" and verify template integrates properly

### Implementation for User Story 8

- [X] T112 [P] [US8] Create Template entity structure in internal/templates/base.go
- [X] T113 [P] [US8] Create TemplateParameter handling in internal/templates/base.go
- [X] T114 [P] [US8] Create middleware templates (auth, logging, rate limiting) in internal/templates/middleware.go
- [X] T115 [P] [US8] Create error handler templates in internal/templates/error_handler.go
- [X] T116 [P] [US8] Create deployment config templates in internal/templates/deployment.go
- [X] T117 [US8] Define MCPTool struct for generate_template in tools/generate_template.go
- [X] T118 [US8] Implement input validation (template_type, parameters) in tools/generate_template.go
- [X] T119 [US8] Select appropriate template based on type in tools/generate_template.go
- [X] T120 [US8] Inject parameters into template in tools/generate_template.go
- [X] T121 [US8] Generate integration instructions in tools/generate_template.go
- [X] T122 [US8] Verify generated template compiles in tools/generate_template.go
- [X] T123 [US8] Generate success response with usage examples in tools/generate_template.go
- [X] T124 [US8] Register generate_template tool in main.go
- [X] T125 [US8] Create integration test for template generation in tests/integration/template_test.go

**Checkpoint**: User Story 8 complete - code templates can be generated ✅

---

## Phase 11: User Story 9 - Query Framework Documentation (Priority: P9)

**Goal**: Provide framework concept explanations and migration guidance without leaving dev environment

**Independent Test**: Ask "explain middleware in go-zero" and verify accurate, helpful response

### Implementation for User Story 9

- [X] T126 [P] [US9] Create concept documentation database in internal/docs/concepts.go
- [X] T127 [P] [US9] Create migration guide database in internal/docs/migration.go
- [X] T128 [US9] Define MCPTool struct for query_docs in tools/query_docs.go
- [X] T129 [US9] Implement input validation (query text) in tools/query_docs.go
- [X] T130 [US9] Implement keyword extraction from query in tools/query_docs.go
- [X] T131 [US9] Search concept database for relevant content in tools/query_docs.go
- [X] T132 [US9] Search migration guides for relevant content in tools/query_docs.go
- [X] T133 [US9] Format response with examples and links in tools/query_docs.go
- [X] T134 [US9] Register query_docs tool in main.go
- [X] T135 [US9] Create unit test for documentation queries in tests/unit/docs_test.go

**Checkpoint**: User Story 9 complete - documentation queries enabled ✅

---

## Phase 12: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [X] T136 [P] Add comprehensive logging across all tools
- [X] T137 [P] Add performance monitoring for tool response times
- [X] T138 [P] Create README.md with installation and usage instructions
- [X] T139 [P] Create CONTRIBUTING.md with development guidelines
- [X] T140 [P] Add unit tests for validation package in tests/unit/validation_test.go
- [ ] T141 [P] Add unit tests for goctl discovery in tests/unit/discovery_test.go
- [X] T142 [P] Add unit tests for fixer package in tests/unit/fixer_test.go
- [X] T143 [P] Add unit tests for analyzer package in tests/unit/analyzer_test.go
- [X] T144 Validate all tools against quickstart.md scenarios
- [X] T145 Add error recovery examples to documentation
- [X] T146 Security audit of credential handling
- [ ] T147 Performance optimization for project analysis caching
- [X] T148 Add CI/CD pipeline configuration
- [X] T149 Create release documentation

**Checkpoint**: Core polish tasks complete - logging, metrics, and documentation added ✅

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phases 3-11)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3 → ... → P9)
- **Polish (Phase 12)**: Depends on all desired user stories being complete

### User Story Dependencies

All user stories are **independently implementable** after Foundational phase:

- **User Story 1 (P1)**: No dependencies on other stories
- **User Story 2 (P2)**: No dependencies on other stories (uses API parser)
- **User Story 3 (P3)**: No dependencies on other stories (uses proto parser)
- **User Story 4 (P4)**: No dependencies on other stories (uses model generator)
- **User Story 5 (P5)**: No dependencies on other stories (creates specs for US2)
- **User Story 6 (P6)**: Uses parsers from US2/US3 but doesn't block them
- **User Story 7 (P7)**: No dependencies on other stories
- **User Story 8 (P8)**: No dependencies on other stories
- **User Story 9 (P9)**: No dependencies on other stories

### Within Each User Story

General pattern (varies by story):

1. Entity/parser creation (parallelizable with [P])
2. Tool struct definition
3. Input validation
4. Core implementation
5. Success response formatting
6. Tool registration in main.go
7. Integration tests

### Parallel Opportunities

#### Setup Phase (Phase 1)
All T003-T012 can run in parallel (directory creation)

#### Foundational Phase (Phase 2)
All T013-T021 can run in parallel (different packages)

#### After Foundational Phase
**All user stories can be developed in parallel** by different team members:
- Team Member 1: User Story 1 (T024-T037)
- Team Member 2: User Story 2 (T038-T048)
- Team Member 3: User Story 3 (T049-T059)
- Team Member 4: User Story 4 (T060-T072)
- And so on...

#### Within User Stories
Tasks marked [P] within same story can run in parallel

#### Polish Phase (Phase 12)
Most tasks T136-T143, T145-T149 can run in parallel

---

## Parallel Example: User Story 1

If working on User Story 1 with a team:

```bash
# Developer A: Entity and validation (parallel)
T024-T028 (validation logic)

# Developer B: Core generation logic (after A completes validation)
T029-T031 (goctl execution and fixes)

# Developer C: Verification and response (after B completes generation)
T032-T033 (build verification and response)

# Developer D: Integration (after C completes handler)
T034 (registration in main.go)

# Developer E: Testing (parallel with D, after handler complete)
T035-T037 (integration tests)
```

---

## Parallel Example: Multiple User Stories

After Foundational phase completes:

```bash
# Week 1: MVP - User Story 1 only
T024-T037 (Complete P1 feature)

# Week 2: Expand capabilities - User Stories 2, 3, 4 in parallel
Team A: T038-T048 (US2 - Generate from spec)
Team B: T049-T059 (US3 - RPC services)
Team C: T060-T072 (US4 - Database models)

# Week 3: Additional features - User Stories 5, 6, 7 in parallel
Team A: T073-T082 (US5 - Create specs)
Team B: T083-T098 (US6 - Project analysis)
Team C: T099-T111 (US7 - Config management)

# Week 4: Final features - User Stories 8, 9 in parallel + polish
Team A: T112-T125 (US8 - Templates)
Team B: T126-T135 (US9 - Documentation)
Team C: T136-T149 (Polish tasks)
```

---

## Implementation Strategy

### MVP Scope (Minimum Viable Product)

**Recommended MVP**: User Story 1 only (T001-T037)

This provides:
- Complete project setup
- Core infrastructure (goctl discovery, validation, fixing)
- API service creation (most common use case)
- Full end-to-end workflow demonstration
- Foundation for all other user stories

**MVP Delivery Time**: ~2-3 weeks for single developer

### Incremental Delivery

After MVP, add user stories in priority order:

1. **Release 1.0** (MVP): User Story 1
2. **Release 1.1**: Add User Story 2 (spec-to-code)
3. **Release 1.2**: Add User Story 3 (RPC services)
4. **Release 1.3**: Add User Story 4 (database models)
5. **Release 2.0**: Add User Stories 5-9 (full feature set)

Each release is independently valuable and testable.

---

## Task Summary

**Total Tasks**: 149

**Task Count by Phase**:
- Phase 1 (Setup): 12 tasks
- Phase 2 (Foundational): 11 tasks
- Phase 3 (US1): 14 tasks
- Phase 4 (US2): 11 tasks
- Phase 5 (US3): 11 tasks
- Phase 6 (US4): 13 tasks
- Phase 7 (US5): 10 tasks
- Phase 8 (US6): 16 tasks
- Phase 9 (US7): 13 tasks
- Phase 10 (US8): 14 tasks
- Phase 11 (US9): 10 tasks
- Phase 12 (Polish): 14 tasks

**Parallelization Potential**:
- Setup phase: 10 parallel tasks
- Foundational phase: 9 parallel tasks
- All 9 user stories can run in parallel after foundational phase
- Polish phase: 11 parallel tasks

**Format Validation**: ✅ All tasks follow required checklist format:
- [x] Checkbox present (- [ ])
- [x] Task ID sequential (T001-T149)
- [x] [P] marker on parallelizable tasks
- [x] [US#] marker on user story tasks
- [x] File paths included in descriptions

---

## Next Steps

1. **Start with Setup**: Complete Phase 1 (T001-T012) to establish project structure
2. **Build Foundation**: Complete Phase 2 (T013-T023) - CRITICAL for all user stories
3. **Deliver MVP**: Complete Phase 3 (T024-T037) - User Story 1 for immediate value
4. **Expand Features**: Add user stories in priority order (P2 → P3 → P4...)
5. **Polish**: Complete Phase 12 after desired user stories are done

**Recommended First Sprint**: T001-T037 (Setup + Foundational + MVP)
