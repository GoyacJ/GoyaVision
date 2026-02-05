# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**GoyaVision** is an intelligent media processing platform built with Go and Vue 3. It provides media asset management, AI operator orchestration, and workflow automation for video/audio/image processing. The system integrates with MediaMTX for streaming (RTSP/RTMP/HLS/WebRTC) and uses a DAG-based workflow engine for complex media processing pipelines.

**Core Philosophy**: Business = Configuration, Capability = Plugin, Execution = Engine

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

Primary config: `configs/config.<env>.yaml`

Environment variable override pattern: `GOYAVISION_*` prefix
- `GOYAVISION_DB_DSN` - Database connection string
- `GOYAVISION_JWT_SECRET` - JWT signing secret (CHANGE IN PRODUCTION!)
- `GOYAVISION_MEDIAMTX_API_ADDRESS` - MediaMTX API endpoint

Default credentials:
- Username: `admin`
- Password: `admin123` (⚠️ Change immediately in production)

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

Engine implementation: `internal/adapter/engine/simple_engine.go`

## Development Workflow

**CRITICAL**: Before any coding, review:
1. `docs/requirements.md` - Functional specifications
2. `docs/architecture.md` - Detailed design decisions
3. `docs/development-progress.md` - Current feature status
4. `CHANGELOG.md` (especially `[未发布]` section) - Recent changes
5. `.cursor/rules/goyavision.mdc` - Project rules and conventions
6. `.cursor/rules/development-workflow.mdc` - Development process

**After completing any feature or bug fix, MUST**:
1. Update `docs/development-progress.md` - Mark status (✅/🚧/⏸️)
2. Update `CHANGELOG.md` - Add entry under `[未发布]` by type (新增/变更/修复/弃用/移除/安全)
3. Update API docs (`docs/api.md`) if endpoints changed
4. Update architecture/requirements docs if design changed
5. Commit with Conventional Commits format: `<type>(<scope>): <subject>`

**Commit types**: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`, `perf`, `style`

**Scopes**: `asset`, `operator`, `workflow`, `task`, `auth`, `api`, `ui`

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
