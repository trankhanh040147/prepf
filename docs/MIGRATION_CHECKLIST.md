# Crush to Prepf Migration Checklist

**Crush codebase**: Folder `crush` in the root of the project.
This document tracks the migration from `charmbracelet/crush` architecture to `prepf` (Technical Interview Coach CLI).

## Overview

**Key Finding**: `crush` is already a local CLI tool. No SSH/remote server code exists. The "server" references are for LSP (Language Server Protocol) and MCP (Model Context Protocol) - local development tools.

## Architecture Mappings

- `crush` session → `prepf` session (the main user session)
- `crush` agent → `prepf` agent (the AI assistant)
- `crush` config → `prepf` config
- `crush` message → `prepf` message
- `crush` tool → `prepf` tool

---

## Phase 1: Foundation (Core Infrastructure) ✅

### Copy Pure UI/Config Packages

- [x] Copy `crush/internal/tui/styles/` → `internal/ui/styles/`
  - [x] `theme.go` - Core theme definitions
  - [x] `charmtone.go` - Default theme
  - [x] `chroma.go` - Syntax highlighting
  - [x] `icons.go` - Icon definitions
  - [x] `markdown.go` - Markdown rendering

- [x] Copy `crush/internal/tui/components/core/` → `internal/ui/components/core/`
  - [x] `core.go`
  - [x] `layout/layout.go`
  - [x] `status/status.go`
  - [x] Test files

- [x] Copy `crush/internal/tui/components/dialogs/` → `internal/ui/components/dialogs/`
  - [x] `dialogs.go`
  - [x] `models/` (model selection dialogs)
  - [x] `permissions/` (simplified for session context)
  - [x] `quit/` (quit confirmation)
  - [x] `filepicker/` (if needed for solution files)
  - [x] Remove: `copilot/`, `hyper/` (service-specific OAuth)
  - [x] Remove: `sessions/` (will adapt to session context later)

- [x] Copy `crush/internal/tui/components/completions/` → `internal/ui/components/completions/`
  - [x] `completions.go`
  - [x] `keys.go`

- [x] Copy `crush/internal/tui/components/anim/` → `internal/ui/components/anim/`
  - [x] `anim.go`

- [x] Copy `crush/internal/tui/util/` → `internal/ui/util/`
  - [x] `util.go`
  - [x] `shell.go`

- [x] Copy `crush/internal/tui/keys.go` → `internal/ui/keys.go`

- [x] Copy utility packages:
  - [x] `csync/` → `internal/csync/`
  - [x] `pubsub/` → `internal/pubsub/`
  - [x] `event/` → `internal/event/`
  - [x] `format/` → `internal/format/`
  - [x] `stringext/` → `internal/stringext/`
  - [x] `filepathext/` → `internal/filepathext/`
  - [x] `ansiext/` → `internal/ansiext/`

- [x] Copy `crush/internal/config/` structure → `internal/config/` (skeleton)
  - [x] Copy `load.go`, `merge.go`, `resolve.go`, `provider.go` (structure)
  - [x] Copy `config.go` (as template for adaptation)
  - [x] Copy test files for reference

### Update Package Imports

- [x] Replace all `github.com/charmbracelet/prepf/internal` → `prepf/internal` (or appropriate package path)
- [x] Update all internal imports in copied files
- [x] Verify no broken imports remain

---

## Phase 2: Data Layer Adaptation

### Database Schema Migration

- [x] Create `internal/db/sql/sessions.sql` (copy from `crush/internal/db/sql/sessions.sql`)
  - [x] Remove `parent_session_id` column (sessions are standalone)
  - [x] Add fields relevant to session context if needed:
    - [x] `difficulty` (TEXT)
    - [x] `topic` (TEXT)
    - [x] `status` (TEXT)
  - [x] Keep: `id`, `title`, `message_count`, `prompt_tokens`, `completion_tokens`, `cost`, `summary_message_id`, `todos`, `created_at`, `updated_at`

- [x] Create initial migration: `internal/db/migrations/000001_initial_sessions.sql`
  - [x] Copy structure from crush's initial migration
  - [x] Adapt as needed

- [x] Update `internal/db/models.go`:
  - [x] Keep `Session` struct and update as needed
  - [x] Remove `ParentSessionID` field

- [x] Run `sqlc generate`
- [x] Verify generated code compiles

### Session Service

- [x] Copy `crush/internal/session/session.go` → `internal/session/session.go`
- [x] Keep package: `session`
- [x] Keep `Session` struct
- [x] Keep `Service` interface (implement for session)
- [x] Update methods as needed:
  - [x] `Create()`, `Get()`, `List()`, `Save()`, `UpdateTitleAndUsage()`, `Delete()` (keep)
  - [x] Remove agent-specific: `CreateAgentToolSessionID()`, `ParseAgentToolSessionID()`, `IsAgentToolSession()`
- [x] Update `Todo` struct (keep or adjust)
- [x] Keep references `session.Session`
- [x] Update pubsub: `pubsub.Broker[Session]`

---

## Phase 3: Agent System

### Agent Package

- [x] Create `internal/agent/` directory

- [x] Copy `crush/internal/agent/coordinator.go` → `internal/agent/coordinator.go`
  - [x] Keep package: `agent`
  - [x] Keep `Coordinator` interface and `coordinator` struct
  - [x] Remove `hyper` provider refs
  - [x] Keep providers: OpenAI, Anthropic, Google, Azure, Bedrock, OpenRouter

- [x] Copy `crush/internal/agent/agent.go` → `internal/agent/agent.go`
  - [x] Keep `SessionAgent` interface and `sessionAgent` struct
  - [x] Keep `SessionAgentOptions`
  - [x] Adjust logic as needed
  - [x] Keep session references

- [x] Copy `crush/internal/agent/errors.go` → `internal/agent/errors.go`
- [x] Copy `crush/internal/agent/event.go` → `internal/agent/event.go` (if exists)

### Template Updates

- [x] Create `internal/agent/templates/` directory

<Not done yet>
- [ ] Create `internal/agent/templates/agent.md.tpl` (use interview/coach context)
- [ ] Create `internal/agent/templates/title.md.tpl`
- [ ] Create `internal/agent/templates/summary.md`
- [ ] Remove any unnecessary task or agent_tool templates

- [ ] Copy `crush/internal/agent/prompts.go` → `internal/agent/prompts.go`
  - [ ] Update prompts as needed for session/agent

---

## Phase 4: UI Adaptation

### Chat UI

- [x] Copy `crush/internal/tui/page/chat/` → `internal/ui/page/chat/`

- [x] Copy `chat.go` → `chat.go`
  - [x] Keep package: `chat`
  - [x] Keep naming as is (`ChatPageID`, `chatCmp`, etc)
  - [x] UI text can stay "Chat", "Session"
  - [x] References: `app.Sessions`, `session.Session` remain

- [x] Copy chat/messages/, chat/editor/, chat/header/, chat/sidebar/, chat/todos/, chat/splash/ as-is
  - [x] Update onboarding and messaging as needed

- [x] Update `internal/ui/page/page.go` if needed

### Component Updates

- [ ] Update `internal/ui/tui.go` (copy from `crush/internal/tui/tui.go`) <Not done yet>
  - [x] Keep page/component names
  - [x] Keep references to sessions/chat as is

- [x] Update `internal/ui/components/dialogs/sessions/`
  - [x] Keep package: `sessions`
  - [x] Keep naming/references: `Session`
  - [x] Update UI text if needed

---

## Phase 5: Configuration Cleanup

### Config Simplification

- [x] Update `internal/config/config.go`
  - [ ] Remove `LSP`/`MCP` types, maps, and config structs
  - [ ] Remove `projects` config
  - [ ] Add session config as needed (difficulty, topics, time limits)
  - [x] Keep `Providers`, `Models`, `Permissions`, `TUIOptions`
  - [ ] Remove complex `Agent` struct if possible

- [ ] Update `internal/config/load.go`
  - [ ] Remove LSP/MCP/project loading
  - [x] Keep provider and model loading

- [ ] Update `internal/config/init.go`
  - [ ] Add session defaults

- [ ] Update `internal/config/provider.go`
  - [ ] Remove hyper/copilot providers

- [ ] Remove `internal/config/copilot.go`
- [ ] Remove `internal/config/hyper.go`
- [ ] Remove `internal/config/lsp_defaults_test.go`

### Provider Configuration

- [x] Keep `internal/oauth/claude/`
- [x] Keep `internal/oauth/token.go`
- [ ] Remove `internal/oauth/copilot/`
- [ ] Remove `internal/oauth/hyper/`

---

## Phase 6: Tool Cleanup

### Tool Trimming

- [ ] Copy `crush/internal/agent/tools/` → `internal/agent/tools/` (selective)

- [x] **KEEP**:
  - [x] `bash.go`, `edit.go`, `view.go`, `write.go`, `grep.go`, `ls.go`, `glob.go`, `todos.go`, `tools.go`, template files

- [ ] **DELETE**:
  - [ ] LSP-specific: `diagnostics.go`, `references.go`
  - [ ] Background jobs: `job_kill.go`, `job_output.go`
  - [ ] Agent/MCP: `agentic_fetch.go`, `agentic_fetch_tool.go`, `agent_tool.go`, `mcp/`
  - [ ] `sourcegraph.go`

- [ ] **EVALUATE** (copy, may remove):
  - [ ] `download.go`, `fetch.go`, `web_search.go`, `web_fetch.go`, `multiedit.go`

- [ ] Update `internal/agent/tools/tools.go`
  - [ ] Remove deleted tools from registration/imports
  - [ ] Adjust imports for package move

- [ ] Update tool files:
  - [ ] Keep package name: `tools`
  - [ ] Keep references: `session.Service`
  - [ ] Update tool descriptions/context only as needed

### Permission Simplification

- [ ] Copy `crush/internal/permission/` → `internal/permission/`
- [ ] Update for simpler (session-based) permissions
- [ ] Keep basic tool whitelisting/blacklisting

### Shell Package

- [ ] Copy `crush/internal/shell/` → `internal/shell/`
  - [ ] Keep as-is
  - [ ] Update package imports if needed

---

## Phase 7: Application Wiring

### App Structure

- [x] Copy `crush/internal/app/app.go` → `internal/app/app.go`
- [x] Update `App` struct as needed
  - [x] Keep `Sessions`, `AgentCoordinator`, `Messages`, `History`, `Permissions`
  - [x] Remove `LSPClients`

- [x] Update `New()`
  - [x] Use `session.NewService()`
  - [x] Remove any LSP logic
  - [ ] Remove any MSP logic
  - [x] Use `InitCoderAgent()` for agent setup

- [x] Keep `InitCoderAgent()` or equivalent for agent initialization

- [x] Remove `internal/app/lsp.go` and `internal/app/lsp_events.go`

- [ ] Keep references in code to `Sessions`/`AgentCoordinator`

### Command Updates

- [x] Update `internal/cmd/root.go` (or as needed)
  - [x] App name, help
  - [ ] Remove projects/dirs
  - [ ] Keep version
  - [ ] Adapt login (keep OpenAI/Anthropic)

- [ ] Add/Update `internal/cli/run.go` for non-interactive mode
  - [ ] Use agent

- [ ] Add/Update login command as needed
  - [ ] Remove copilot/hyper logins
  - [ ] Keep claude/openai

- [ ] Remove: projects.go, dirs.go, schema.go (if not needed), update_providers.go (if not needed)

### Main Application Entry

- [ ] Update `cmd/prepf/main.go` or create new
  - [ ] App/session init
  - [ ] UI setup

---

## Phase 8: Cleanup & Testing

### Remove Unused Packages

- [ ] Remove `internal/lsp/`, `internal/projects/`, `internal/skills/`
- [ ] Remove unused oauth/config files

### Update Dependencies

- [ ] Review `go.mod`
  - [ ] Keep core dependencies: bubbletea, lipgloss, fantasy, catwalk
  - [ ] Remove LSP/MCP libraries
  - [ ] Update module path if needed

### Testing

- [ ] UI runs
- [ ] Session creation
- [ ] Agent interaction
- [ ] Message persistence
- [ ] Tool execution (bash, edit, etc)
- [ ] Config loading and authentication works
- [ ] Migrations work
- [ ] Session listing/selection
- [ ] Non-interactive mode works

### Documentation

- [ ] Update README.md, DEVELOPMENT.md
- [ ] Document config, agent system, and API if any

---

## Notes & Considerations

### Key Files to Reference

- `crush/internal/tui/styles/theme.go` - Core styling system
- `crush/internal/config/config.go` - Configuration structure
- `crush/internal/agent/coordinator.go:708-754` - Provider setup (already uses standard providers)
- `crush/internal/session/session.go` - Session model
- `crush/internal/app/app.go:66-118` - Application initialization
- `crush/internal/tui/tui.go:688-710` - TUI initialization

### Important Adaptations

- **Provider Setup**: Already uses standard providers (OpenAI, Anthropic, etc.) - minimal changes needed
- **Session Creation**: Remove parent session logic
- **Agent Run**: Adjust prompt for interview context
- **UI Rendering**: All wording remains "Chat", "Session" unless needed

### Migration Order

1. Foundation (UI/styles/utils)
2. Data layer (session)
3. Agent system
4. UI adaptation
5. Configuration
6. Tool cleanup
7. App wiring
8. Testing

---

**Status**: 🚧 In Progress

**Last Updated**: 2025-01-27

**Phase 1 Status**: ✅ Complete - All foundation infrastructure copied and imports updated

