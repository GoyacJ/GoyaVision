# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**GoyaVision** is an intelligent media processing platform built with Go and Vue 3. It provides media asset management, AI operator orchestration, and workflow automation for video/audio/image processing. The system integrates with MediaMTX for streaming (RTSP/RTMP/HLS/WebRTC) and uses a DAG-based workflow engine for complex media processing pipelines.

**Core Philosophy**: Business = Configuration, Capability = Plugin, Execution = Engine

## Information Completeness and Questioning Rules

Before executing any user request, perform an information completeness check:

**When to ask clarifying questions:**
- **Missing critical information**: If essential information needed to complete the task is missing, you MUST ask clarifying questions before proceeding
- **Ambiguity exists**: If there are multiple reasonable interpretations or execution paths, you MUST point out the ambiguity and ask for user preference
- **High-risk operations**: If continuing would lead to high-risk errors or irreversible consequences, you MUST confirm user intent first

**Questioning standards:**
- Ask only the minimum necessary questions at once (maximum 3)
- Questions must be specific, actionable, and answerable - avoid vague questions
- Do not repeat information already confirmed
- When information is sufficient, execute the task directly without asking

**Prohibited behaviors:**
- Making assumptions when critical information is missing
- Continuing to generate results blindly just to appear "helpful"

This ensures quality and prevents mistakes caused by insufficient information.

## Build Commands

### Backend
```bash
# Build Go binary
go build -o bin/goyavision ./cmd/server

# Run locally (reads configs/config.<env>.yaml)
GOYAVISION_ENV=dev ./bin/goyavision

# Or run directly
go run cmd/server/main.go

# Format code
go fmt ./...

# Vet code
go vet ./...
```

### Frontend
```bash
cd web

# Install dependencies (uses pnpm)
pnpm install

# Development server with hot reload
pnpm run dev

# Production build (outputs to web/dist/)
pnpm run build

# Preview production build
pnpm run preview
```

### Combined Build
```bash
# Build both frontend and backend
make build-all

# Or separately
make build-web    # Frontend only
make build        # Backend only

# Clean build artifacts
make clean
```

### Docker
```bash
# Start all services (GoyaVision, PostgreSQL, MediaMTX)
docker-compose up -d

# View logs
docker-compose logs -f goyavision

# Stop services
docker-compose down
```

## Architecture Overview

GoyaVision follows **Clean Architecture** with strict layered separation:

```
internal/
├── domain/      # Core entities (MediaSource, MediaAsset, Operator, Workflow, Task, Artifact, User, Role, Menu)
├── port/        # Interface definitions (Repository, MediaMTXClient, WorkflowEngine, OperatorPort)
├── app/         # Business services (MediaSourceService, WorkflowService, TaskService, AuthService, etc.)
├── adapter/     # Infrastructure implementations
│   ├── persistence/   # GORM + PostgreSQL repositories
│   ├── mediamtx/      # MediaMTX HTTP client
│   ├── engine/        # DAG workflow execution engine
│   └── ai/            # AI operator HTTP client
└── api/         # HTTP layer
    ├── handler/       # Request handlers
    ├── dto/           # Data Transfer Objects
    ├── middleware/    # JWT auth, CORS, logging
    └── router.go      # Route registration
```

### Dependency Rules (CRITICAL)
1. **Domain** → No dependencies (innermost layer)
2. **Port** → May depend on Domain
3. **App** → May depend on Domain + Port (not Adapter)
4. **Adapter** → Implements Port interfaces, may use Domain
5. **API** → May depend on App, Port, Domain (not Adapter directly)

**When adding features**: Start with Domain entities → Port interfaces → App services → Adapter implementations → API handlers + DTOs.

### App Layer Structure (CQRS)

The application layer follows **CQRS (Command Query Responsibility Segregation)** pattern with 39 handlers:

**Commands** (`internal/app/command/`) - Write operations:
- Asset library: `create_source.go`, `update_source.go`, `delete_source.go`, `create_asset.go`, `update_asset.go`, `delete_asset.go`
- Operator center: `create_operator.go`, `update_operator.go`, `delete_operator.go`, `enable_operator.go`
- Task center: `create_workflow.go`, `update_workflow.go`, `delete_workflow.go`, `enable_workflow.go`, `create_task.go`, `start_task.go`, `update_task.go`, `complete_task.go`, `fail_task.go`, `cancel_task.go`
- Auth: `login.go`

**Queries** (`internal/app/query/`) - Read operations:
- Asset library: `get_source.go`, `list_sources.go`, `get_asset.go`, `list_assets.go`, `get_asset_tags.go`, `list_asset_children.go`
- Operator center: `get_operator.go`, `get_operator_by_code.go`, `list_operators.go`
- Task center: `get_workflow.go`, `get_workflow_by_code.go`, `get_workflow_with_nodes.go`, `list_workflows.go`, `get_task.go`, `get_task_with_relations.go`, `get_task_stats.go`, `list_tasks.go`, `list_running_tasks.go`
- Auth: `get_profile.go`

**Ports** (`internal/app/port/`) - Application boundary interfaces:
- `media_gateway.go` - MediaMTX gateway interface
- `object_storage.go` - Object storage interface (MinIO/S3)
- `token_service.go` - Token service interface
- `event_bus.go` - Event bus interface
- `unit_of_work.go` - Transaction management interface

**Other services**:
- `artifact.go` - Artifact management (query, association)
- `file.go` - File management (upload, download, metadata)
- `user_management.go` - User management service (CRUD, role assignment)
- `workflow_scheduler.go` - Workflow scheduler (cron tasks, event triggers using gocron/v2)

## Core Concepts

The system centers around a media processing pipeline:

```
MediaSource → MediaAsset → Operator → Workflow → Task → Artifact
```

- **MediaSource**: Streaming sources (RTSP/RTMP pull/push) or file uploads
- **MediaAsset**: Video/image/audio resources with metadata and tagging
- **Operator**: AI/processing capability units (analyze/edit/generate/transform categories)
- **Workflow**: DAG-based orchestration of operators with conditional branching
- **Task**: Workflow execution instances with status tracking
- **Artifact**: Output products (new assets, structured results, timelines, diagnostics)

### Operator Standard Protocol

All operators MUST follow this I/O contract:

**Input**:
```json
{
  "asset_id": "uuid",
  "params": {"key": "value"}
}
```

**Output**:
```json
{
  "output_assets": [{"type": "video|image|audio", "path": "...", "format": "...", "metadata": {}}],
  "results": [{"type": "detection|classification|...", "data": {}, "confidence": 0.95}],
  "timeline": [{"start": 0.0, "end": 5.0, "event_type": "...", "confidence": 0.95, "data": {}}],
  "diagnostics": {"latency_ms": 150, "model_version": "v1.0", "device": "gpu"}
}
```

## Configuration

Primary config: `configs/config.yaml` (default) or `configs/config.<env>.yaml` for environment-specific settings.

Environment variable override pattern: `GOYAVISION_*` prefix
- `GOYAVISION_ENV` - Environment (dev/prod/test)
- `GOYAVISION_DB_DSN` - Database connection string
- `GOYAVISION_JWT_SECRET` - JWT signing secret (CHANGE IN PRODUCTION!)
- `GOYAVISION_JWT_ACCESS_EXPIRE` - Access token expiry (default: 2h)
- `GOYAVISION_JWT_REFRESH_EXPIRE` - Refresh token expiry (default: 168h = 7 days)
- `GOYAVISION_MEDIAMTX_API_ADDRESS` - MediaMTX API endpoint

Default credentials:
- Username: `admin`
- Password: `admin123` (⚠️ Change immediately in production)

**Configuration precedence**: Environment variables > config.<env>.yaml > config.yaml

## Key Technologies

**Backend**:
- Go 1.22+, Echo v4 (HTTP), GORM v1.25 (ORM), PostgreSQL 12+
- Viper (config), golang-jwt/jwt/v5 (auth), bcrypt (password hashing)
- gocron/v2 (scheduling), google/uuid (IDs)

**Media Stack**:
- MediaMTX (RTSP/RTMP/HLS/WebRTC streaming server)
- FFmpeg (frame extraction for AI processing)

**Frontend**:
- Vue 3.4 + TypeScript 5.3, Vite 5.0 (build)
- Element Plus 2.5 (UI), Tailwind CSS 3.4 (styling)
- Pinia 2.1 (state), Vue Router 4.2 (routing)
- Video.js 8.6 (HLS/video playback), Axios 1.6 (HTTP)

## API Structure

All endpoints under `/api/v1/`:

- `/auth` - Login, logout, token refresh (public)
- `/users`, `/roles`, `/menus` - RBAC management (protected)
- `/sources` - Media source CRUD, recording control
- `/assets` - Media asset CRUD, search, tagging
- `/operators` - AI operator management
- `/workflows` - DAG workflow orchestration
- `/tasks` - Task execution and monitoring
- `/artifacts` - Result artifact retrieval
- `/files` - File upload/download

Authentication: JWT with dual-token mechanism (Access Token 2h, Refresh Token 7d). Middleware in `internal/api/middleware/auth.go` validates tokens and enforces RBAC permissions.

## Frontend Integration

The Vue 3 frontend is embedded into the Go binary via `//go:embed` in `embed.go`. After `make build-web`, the `web/dist/` directory contents are compiled into the binary. The API layer serves static files and handles SPA routing (all non-API routes return `index.html`).

**Frontend structure**:
- `web/src/views/` - Page components (Dashboard, MediaSources, Assets, Workflows, Tasks, etc.)
- `web/src/stores/` - Pinia stores (user, auth, permissions)
- `web/src/api/` - API client modules per domain
- `web/src/components/` - Reusable UI components
- `web/src/composables/` - Reusable composition functions (Phase 2/3 refactoring)

**Frontend Composables Pattern** (Phase 2/3):

The frontend uses a unified composables pattern to eliminate repetitive state management code:

- **`useAsyncData.ts`** - Generic async data fetching with loading/error/data states
  - Handles API calls with automatic state management
  - Supports immediate execution and manual refresh
  - Provides `isLoading`, `error`, `data`, `execute()`

- **`usePagination.ts`** - Pagination state management
  - Manages `page`, `pageSize`, `total`
  - Provides `goToPage()`, `prevPage()`, `nextPage()`, `changePageSize()`
  - Computed properties: `totalPages`, `hasPrevPage`, `hasNextPage`, `startIndex`, `endIndex`

- **`useTable.ts`** - Complete table data management (combines above two)
  - Integrates `useAsyncData` + `usePagination`
  - Automatically reloads on page/pageSize changes
  - Supports reactive `extraParams` for filtering
  - Reduces page code by 60-70%

**Usage example**:
```typescript
const filterParams = computed(() => ({
  keyword: searchKeyword.value || undefined,
  status: selectedStatus.value || undefined
}))

const {
  items,
  isLoading,
  error,
  pagination,
  goToPage,
  changePageSize,
  refreshTable
} = useTable(
  (params) => assetApi.list(params),
  {
    immediate: true,
    initialPageSize: 20,
    extraParams: filterParams
  }
)
```

This pattern has been applied to all 5 list pages (asset, source, operator, workflow, task) with significant code reduction and improved maintainability.

## Workflow Engine

Workflows are defined as DAGs (Directed Acyclic Graphs) with:
- **Nodes**: Operator invocations with parameters, retry policies, timeouts
- **Edges**: Data flow connections with optional conditions
- **Triggers**: manual, schedule (cron), event (new asset, recording complete)

Execution flow:
1. Trigger activates → Create Task
2. Scheduler (gocron) → WorkflowEngine.Execute()
3. Parse DAG → Topological sort → Execute nodes sequentially/parallel
4. Each node: Fetch input asset → Call Operator → Save Artifact → Pass to downstream
5. All nodes complete → Task success; any node fails → Retry or mark failed

Engine implementation: `internal/infra/engine/dag_engine.go`

**DAG execution details**:
- Uses Kahn's algorithm for topological sorting
- Supports parallel execution of independent nodes
- Validates DAG for cycles before execution
- Each node execution: `OperatorPort.Execute()` → Save `Artifact` → Pass to downstream
- Retry mechanism with configurable attempts and timeout
- Comprehensive error handling with task status updates

## Development Workflow

### Pre-Development (MUST DO)

Before starting any new feature or bug fix:

1. **Read documentation in order** (use `/goya-dev-start` command):
   - `docs/development-progress.md` - Check feature status (✅/🚧/⏸️), avoid duplicate work
   - `CHANGELOG.md` - Focus on `[未发布]` section for recent changes
   - `docs/requirements.md` - Confirm functional requirements and acceptance criteria
   - `docs/architecture.md` - Understand layered architecture and dependency rules
   - `docs/api.md` - Check existing endpoints to avoid duplication

2. **Verify architectural constraints**:
   - Domain → No external dependencies
   - Port → May depend on Domain
   - App → May depend on Domain + Port (NEVER Adapter)
   - Adapter → Implements Port, may use Domain
   - API → May depend on App + Port + Domain (NEVER Adapter directly)

3. **Code style requirements**:
   - Go: Use `gofmt`/`goimports`, no end-of-line comments
   - Frontend: TypeScript strict mode, Vue 3 Composition API
   - Error handling: Use `internal/api/errors.go` error types

### During Development

- Follow existing code patterns and naming conventions
- Do not introduce new styles or patterns without discussion
- Never swallow errors - all errors must be returned or logged
- Use `context.Context` for cancellable operations
- Always use DTOs in API layer - never expose domain entities directly
- For frontend list pages, use the composables pattern (`useTable`, `useAsyncData`, `usePagination`)

### Post-Development (MUST DO)

After completing any feature or bug fix, execute in order (use `/goya-dev-done` command):

1. **Update development progress**:
   - Edit `docs/development-progress.md`
   - Update feature status (✅ Completed / 🚧 In Progress / ⏸️ Pending / ⚠️ Blocked / 🔄 Refactoring)
   - Add notes and completion dates for major features

2. **Update changelog**:
   - Edit `CHANGELOG.md` under `[未发布]` section
   - Categorize by type: 新增 (feat) / 变更 (change) / 修复 (fix) / 弃用 (deprecated) / 移除 (removed) / 安全 (security)
   - Write clear, user-facing descriptions

3. **Update related documentation** (if applicable):
   - If API changed: Update `docs/api.md` with request/response examples
   - If architecture changed: Update `docs/architecture.md` or `docs/requirements.md`
   - If user-facing: Update `README.md` or `docs/DEPLOYMENT.md`

4. **Pre-commit checklist**:
   - [ ] Code tested (unit tests or manual testing)
   - [ ] All documentation updated
   - [ ] Code formatted (`gofmt`/`goimports` or Prettier)
   - [ ] No debug code (console.log, temporary comments)
   - [ ] Linter errors fixed
   - [ ] Commit message follows Conventional Commits format

5. **Git commit** (use `/goya-commit` command):
   - Format: `<type>(<scope>): <subject>` (subject in Chinese)
   - Types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `perf`, `style`
   - Scopes: `asset`, `operator`, `workflow`, `task`, `auth`, `api`, `ui`
   - Example: `feat(asset): 实现媒体资产管理功能`

**Enforcement**: Both AI agents and developers MUST follow steps 1-5 after any code changes. Skipping documentation updates or using non-standard commit messages is prohibited.

## Claude Code Commands

Claude Code provides custom commands in `.claude/commands/` to streamline your development workflow:

### 🚀 Development Workflow Commands

- **`/goya-dev-start`** - Pre-development checklist
  - Reviews required documentation in order (progress → changelog → requirements → architecture → API)
  - Explains core architectural rules and dependencies
  - Provides quick reference for credentials, ports, and build commands
  - **Use when**: Starting any new feature or bug fix

- **`/goya-dev-done`** - Post-development checklist
  - Guides through updating development progress
  - Ensures CHANGELOG is updated under `[未发布]` section
  - Reminds to update API/architecture docs if applicable
  - Provides pre-commit self-check items
  - **Use when**: Feature/fix is complete, before committing

- **`/goya-commit`** - Create standardized Git commit
  - Explains Conventional Commits format
  - Lists commit types and scopes
  - Shows good vs bad commit examples
  - Includes commit message templates
  - **Use when**: Ready to commit code changes

### 📚 Project Context Command

- **`/goya-context`** - Get complete project context
  - Core concepts and data flow (MediaSource → MediaAsset → Operator → Workflow → Task → Artifact)
  - Clean Architecture layers and dependency rules
  - Operator standard I/O protocol
  - Complete API endpoint reference
  - Development status overview (✅/🚧/⏸️)
  - **Use when**: Need to understand architecture, entities, or API contracts

### 📝 Documentation Update Commands

- **`/goya-api-doc`** - API documentation update guide
  - Provides endpoint documentation template
  - Explains request/response format standards
  - Lists common error codes
  - **Use when**: Adding or modifying API endpoints

- **`/goya-progress`** - Development progress update guide
  - Explains status markers (✅/🚧/⏸️/⚠️/🔄)
  - Shows update examples for different scenarios
  - Clarifies when updates are required vs optional
  - **Use when**: Completing features or hitting milestones

### Typical Workflows

**Starting a new feature:**
```
1. /goya-dev-start     # Review pre-development checklist
2. /goya-context       # Understand architecture and APIs
3. [Implement feature]
4. /goya-dev-done      # Complete post-development checklist
5. /goya-api-doc       # Update API docs (if applicable)
6. /goya-progress      # Update development progress
7. /goya-commit        # Create standardized commit
```

**Fixing a bug:**
```
1. /goya-dev-start     # Review checklist
2. [Fix the issue]
3. /goya-dev-done      # Post-fix checklist
4. /goya-commit        # Commit with 'fix' type
```

**Quick context lookup:**
```
/goya-context          # View architecture, entities, APIs, status
```

These commands enforce the same development standards documented in `.cursor/skills/` but are optimized for Claude Code's command interface. See `.claude/commands/README.md` for more details.

## Common Development Patterns

### Creating a new entity

Follow this sequence for maximum architecture compliance:

1. **Domain entity** (`internal/domain/<module>/<entity>.go`):
   ```go
   type MediaAsset struct {
       ID          uuid.UUID `gorm:"type:uuid;primary_key"`
       Name        string
       Type        string
       // ... other fields
       CreatedAt   time.Time
       UpdatedAt   time.Time
   }
   ```

2. **Port interface** (`internal/port/repository.go`):
   ```go
   type Repository interface {
       // MediaAsset methods
       CreateMediaAsset(ctx context.Context, asset *domain.MediaAsset) error
       GetMediaAssetByID(ctx context.Context, id uuid.UUID) (*domain.MediaAsset, error)
       // ... other methods
   }
   ```

3. **App layer commands/queries** (`internal/app/command/` and `internal/app/query/`):
   ```go
   // command/create_asset.go
   type CreateAssetCommand struct {
       repo port.Repository
   }
   func (c *CreateAssetCommand) Execute(ctx context.Context, req CreateAssetRequest) (*domain.MediaAsset, error)

   // query/get_asset.go
   type GetAssetQuery struct {
       repo port.Repository
   }
   func (q *GetAssetQuery) Execute(ctx context.Context, id uuid.UUID) (*domain.MediaAsset, error)
   ```

4. **Adapter implementation** (`internal/adapter/persistence/`):
   ```go
   func (r *repository) CreateMediaAsset(ctx context.Context, asset *domain.MediaAsset) error {
       return r.db.WithContext(ctx).Create(asset).Error
   }
   ```

5. **API layer** (`internal/api/`):
   ```go
   // dto/asset.go - Request/Response DTOs
   type CreateAssetRequest struct {
       Name string `json:"name" validate:"required"`
       Type string `json:"type" validate:"required,oneof=video image audio"`
   }

   type AssetResponse struct {
       ID   string `json:"id"`
       Name string `json:"name"`
       Type string `json:"type"`
   }

   // handler/asset.go - HTTP handlers
   func (h *AssetHandler) Create(c echo.Context) error {
       var req dto.CreateAssetRequest
       // ... bind and validate
       asset, err := h.createAssetCmd.Execute(c.Request().Context(), req)
       // ... convert to DTO and return
   }
   ```

6. **Router registration** (`internal/api/router.go`):
   ```go
   assets := authenticated.Group("/assets")
   assets.POST("", handlers.Asset.Create)
   assets.GET("/:id", handlers.Asset.GetByID)
   ```

### Executing a workflow

1. Trigger activates → Call `command/create_task.go`
2. Create Task → Save via `UnitOfWork` to database
3. `WorkflowScheduler` dispatches → `WorkflowEngine.Execute()`
4. Parse DAG → Topological sort (Kahn's algorithm)
5. Execute nodes → `OperatorPort.Execute()` for each
6. Save `Artifact` → `artifact.go` service handles storage
7. Update Task status → `command/complete_task.go` or `command/fail_task.go`
8. Return execution result

## Deprecated Concepts (Do Not Use)

V1.0 is a complete rewrite. These legacy concepts are NO LONGER VALID:

- ❌ **`Stream`** - Replaced by `MediaSource`
- ❌ **`Algorithm`** - Replaced by `Operator`
- ❌ **`AlgorithmBinding`** - Replaced by `Workflow` (DAG-based)
- ❌ **`InferenceResult`** - Replaced by `Artifact` (with multiple types)

**Database migrations**: V1.0 uses entirely new table schemas. Legacy tables (`streams`, `algorithms`, `algorithm_bindings`, `inference_results`) are not compatible and must be migrated manually if needed.

## Code Style

**Go**:
- Use `gofmt`/`goimports` formatting
- No end-of-line comments (user preference)
- CamelCase for exported, camelCase for private
- Never swallow errors - return or log all errors
- Use context.Context for cancellable operations
- Unified error handling via `internal/api/errors.go` (ErrNotFound, ErrInvalidInput, ErrAlreadyExists, etc.)

**Frontend**:
- TypeScript strict mode
- Composition API (not Options API)
- Element Plus components with Tailwind utilities
- Consistent import ordering

## Database

PostgreSQL with GORM auto-migration on startup.

**Core tables**:
- `media_sources`, `media_assets` - Asset library
- `operators`, `workflows`, `workflow_nodes`, `workflow_edges`, `tasks`, `artifacts` - Processing engine
- `users`, `roles`, `permissions`, `menus`, `user_roles`, `role_permissions`, `role_menus` - RBAC

**Important indexes**: On `source_id`, `task_id`, `workflow_id`, `created_at` for performance.

## Service Ports

- **8080**: GoyaVision (web UI + API)
- **5432**: PostgreSQL
- **8554**: MediaMTX RTSP
- **1935**: MediaMTX RTMP
- **8888**: MediaMTX HLS
- **8889**: MediaMTX WebRTC
- **9997**: MediaMTX Control API

## Common Pitfalls

1. **Don't violate dependency rules**: App layer should never import from Adapter
2. **Always use DTOs in API layer**: Never expose domain entities directly in HTTP responses
3. **Workflow DAGs must be acyclic**: Engine validates on creation, will reject cycles
4. **Operator endpoints must follow standard protocol**: Input/output format is contract
5. **JWT secret in production**: Default in config.dev.yaml is for development only
6. **MediaMTX paths**: Recording paths follow pattern `%path/%Y-%m-%d_%H-%M-%S` (time-based segmentation)
7. **Asset parent_id tracking**: Use for derivative assets (original video → frame → detection result)

## Testing Strategy

- Unit tests: Domain/App layer business logic
- Integration tests: Adapter implementations (DB, HTTP clients)
- E2E tests: API layer full request/response cycles
- Workflow tests: DAG execution, data flow, error handling

## Documentation

Project docs in `docs/`:
- `requirements.md` - Feature specifications
- `architecture.md` - System design details
- `development-progress.md` - Implementation status tracking
- `api.md` - RESTful API reference
- `DEPLOYMENT.md` - Deployment and operations guide

Always keep documentation synchronized with code changes.

## Claude Code vs Cursor/Cline

This project supports multiple AI coding assistants with similar but adapted configurations:

| Tool | Configuration Location | Notes |
|------|----------------------|-------|
| **Claude Code** | `.claude/commands/`, `CLAUDE.md` | This file (CLAUDE.md) serves as project instructions |
| **Cursor** | `.cursor/rules/`, `.cursor/skills/`, `.cursor/hooks.json` | Full rules/skills/hooks system |
| **Cline** | `.cline/skills/`, `.clinerules/` | Simplified rules and skills |

**For Claude Code users**: All project rules, conventions, and workflows are documented in this file (CLAUDE.md). The custom commands in `.claude/commands/` provide quick access to common workflows. Use `/goya-dev-start` before beginning work and `/goya-dev-done` before committing.

**Shared standards**: All tools follow the same development workflow, code style, commit format, and documentation requirements described in this file.
