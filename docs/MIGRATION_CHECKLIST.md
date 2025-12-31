# Crush to Prepf Migration Checklist

**Crush codebase**: Folder `crush` in the root of the project.
This document tracks the migration from `charmbracelet/crush` architecture to `prepf` (Technical Interview Coach CLI).

## Overview

**Key Finding**: `crush` is already a local CLI tool. No SSH/remote server code exists. The "server" references are for LSP (Language Server Protocol) and MCP (Model Context Protocol) - local development tools.

## Architecture Mappings

- `crush` Sessions → `prepf` Interviews (The Gauntlet)
- `crush` Agents → `prepf` AI Coaches
- `crush` Config → `prepf` Profile/Settings
- `crush` Messages → `prepf` Interview Messages/Responses
- `crush` Tools → `prepf` Interview Tools

---

## Phase 1: Foundation (Core Infrastructure)

### Copy Pure UI/Config Packages

- [ ] Copy `crush/internal/tui/styles/` → `internal/ui/styles/`
  - [ ] `theme.go` - Core theme definitions
  - [ ] `charmtone.go` - Default theme
  - [ ] `chroma.go` - Syntax highlighting
  - [ ] `icons.go` - Icon definitions
  - [ ] `markdown.go` - Markdown rendering

- [ ] Copy `crush/internal/tui/components/core/` → `internal/ui/components/core/`
  - [ ] `core.go`
  - [ ] `layout/layout.go`
  - [ ] `status/status.go`
  - [ ] Test files

- [ ] Copy `crush/internal/tui/components/dialogs/` → `internal/ui/components/dialogs/`
  - [ ] `dialogs.go`
  - [ ] `models/` (model selection dialogs)
  - [ ] `permissions/` (simplified for interviews)
  - [ ] `quit/` (quit confirmation)
  - [ ] `filepicker/` (if needed for solution files)
  - [ ] Remove: `copilot/`, `hyper/` (service-specific OAuth)
  - [ ] Remove: `sessions/` (will adapt to interviews later)

- [ ] Copy `crush/internal/tui/components/completions/` → `internal/ui/components/completions/`
  - [ ] `completions.go`
  - [ ] `keys.go`

- [ ] Copy `crush/internal/tui/components/anim/` → `internal/ui/components/anim/`
  - [ ] `anim.go`

- [ ] Copy `crush/internal/tui/util/` → `internal/ui/util/`
  - [ ] `util.go`
  - [ ] `shell.go`

- [ ] Copy `crush/internal/tui/keys.go` → `internal/ui/keys.go`

- [ ] Copy utility packages:
  - [ ] `csync/` → `internal/csync/`
  - [ ] `pubsub/` → `internal/pubsub/`
  - [ ] `event/` → `internal/event/`
  - [ ] `format/` → `internal/format/`
  - [ ] `stringext/` → `internal/stringext/`
  - [ ] `filepathext/` → `internal/filepathext/`
  - [ ] `ansiext/` → `internal/ansiext/`

- [ ] Copy `crush/internal/config/` structure → `internal/config/` (skeleton)
  - [ ] Copy `load.go`, `merge.go`, `resolve.go`, `provider.go` (structure)
  - [ ] Copy `config.go` (as template for adaptation)
  - [ ] Copy test files for reference

### Update Package Imports

- [ ] Replace all `github.com/charmbracelet/crush/internal` → `prepf/internal` (or appropriate package path)
- [ ] Update all internal imports in copied files
- [ ] Verify no broken imports remain

---

## Phase 2: Data Layer Adaptation

### Database Schema Migration

- [ ] Create `internal/db/sql/interviews.sql` (copy from `crush/internal/db/sql/sessions.sql`)
  - [ ] Rename table: `sessions` → `interviews`
  - [ ] Remove `parent_session_id` column (interviews are standalone)
  - [ ] Add interview-specific fields (if needed):
    - [ ] `difficulty` (TEXT: easy/medium/hard)
    - [ ] `topic` (TEXT: algorithms/system-design/etc.)
    - [ ] `status` (TEXT: in-progress/completed/abandoned)
  - [ ] Keep: `id`, `title`, `message_count`, `prompt_tokens`, `completion_tokens`, `cost`, `summary_message_id`, `todos`, `created_at`, `updated_at`

- [ ] Create initial migration: `internal/db/migrations/000001_initial_interviews.sql`
  - [ ] Copy structure from crush's initial migration
  - [ ] Adapt for interviews schema

- [ ] Update `internal/db/models.go`:
  - [ ] Rename `Session` → `Interview`
  - [ ] Update fields to match new schema
  - [ ] Remove `ParentSessionID` field

- [ ] Create `internal/db/sql/interviews.sql` queries file (copy from `sessions.sql`)
  - [ ] Update all query names: `CreateSession` → `CreateInterview`
  - [ ] Update all SQL to use `interviews` table
  - [ ] Remove parent session related queries

- [ ] Run `sqlc generate` to regenerate query code
- [ ] Verify generated code compiles

### Session → Interview Service

- [ ] Copy `crush/internal/session/session.go` → `internal/interview/interview.go`
- [ ] Rename package: `session` → `interview`
- [ ] Rename `Session` struct → `Interview`
- [ ] Rename `Service` interface → `InterviewService` (or keep `Service` in interview package)
- [ ] Update all method signatures:
  - [ ] `Create()` → keep (create interview)
  - [ ] Remove: `CreateTitleSession()`, `CreateTaskSession()` (not needed)
  - [ ] `Get()` → keep
  - [ ] `List()` → keep
  - [ ] `Save()` → keep
  - [ ] `UpdateTitleAndUsage()` → keep
  - [ ] `Delete()` → keep
  - [ ] Remove: `CreateAgentToolSessionID()`, `ParseAgentToolSessionID()`, `IsAgentToolSession()` (agent-specific)

- [ ] Update `Todo` struct (keep same structure, but in interview package)
- [ ] Update all references from `session.Session` → `interview.Interview`
- [ ] Update pubsub broker: `pubsub.Broker[Session]` → `pubsub.Broker[Interview]`

---

## Phase 3: Agent → Coach Transformation

### Agent Package Adaptation

- [ ] Create `internal/coach/` directory

- [ ] Copy `crush/internal/agent/coordinator.go` → `internal/coach/coordinator.go`
  - [ ] Rename package: `agent` → `coach`
  - [ ] Rename `Coordinator` interface → `CoachCoordinator`
  - [ ] Rename `coordinator` struct → `coachCoordinator`
  - [ ] Update method names (keep same signatures, update implementation)
  - [ ] Remove `hyper` provider references (lines referencing `hyper.Name`)
  - [ ] Keep standard providers: OpenAI, Anthropic, Google, Azure, Bedrock, OpenRouter

- [ ] Copy `crush/internal/agent/agent.go` → `internal/coach/coach.go`
  - [ ] Rename `SessionAgent` interface → `InterviewCoach`
  - [ ] Rename `sessionAgent` struct → `interviewCoach`
  - [ ] Update `SessionAgentOptions` → `InterviewCoachOptions`
  - [ ] Update all method implementations
  - [ ] Change references: `sessions.Service` → `interviews.Service`
  - [ ] Update prompt loading logic

- [ ] Copy `crush/internal/agent/errors.go` → `internal/coach/errors.go`
  - [ ] Update package name

- [ ] Copy `crush/internal/agent/event.go` → `internal/coach/event.go` (if exists)
  - [ ] Update package name

### Templates Rewrite

- [ ] Create `internal/coach/templates/` directory

- [ ] Create `internal/coach/templates/interview_coach.md.tpl`
  - [ ] Replace coding assistant prompts with interview coaching prompts
  - [ ] Add interview-specific context: difficulty, topic, time limit
  - [ ] Define coach role: asking questions, providing hints, evaluating solutions

- [ ] Create `internal/coach/templates/interview_title.md.tpl` (copy from `crush/internal/agent/templates/title.md`)
  - [ ] Adapt for interview title generation (e.g., "Two Sum Problem - Easy")

- [ ] Create `internal/coach/templates/summary.md` (copy from crush, adapt if needed)

- [ ] Remove/don't copy:
  - [ ] `coder.md.tpl` (replace with interview_coach.md.tpl)
  - [ ] `task.md.tpl` (not needed)
  - [ ] `agent_tool.md` (not needed)
  - [ ] `initialize.md.tpl` (not needed)
  - [ ] `agentic_fetch_*.tpl` (not needed)

- [ ] Copy `crush/internal/agent/prompts.go` → `internal/coach/prompts.go`
  - [ ] Update package name
  - [ ] Update prompt loading to use new templates
  - [ ] Update system prompt generation for interview context

---

## Phase 4: UI Adaptation

### Chat → Interview UI

- [ ] Copy `crush/internal/tui/page/chat/` → `internal/ui/page/interview/`

- [ ] Rename `chat.go` → `interview.go`
  - [ ] Rename package: `chat` → `interview`
  - [ ] Rename `ChatPageID` → `InterviewPageID`
  - [ ] Rename `chatCmp` → `interviewCmp`
  - [ ] Rename `chatPage` → `interviewPage`
  - [ ] Update all UI text: "Chat" → "Interview", "Session" → "Interview"
  - [ ] Update references: `app.Sessions` → `app.Interviews`
  - [ ] Update references: `session.Session` → `interview.Interview`

- [ ] Update `internal/ui/page/interview/messages/` (copy from chat/messages/)
  - [ ] Update package name
  - [ ] Keep message rendering logic (minimal changes needed)
  - [ ] Update tool rendering if needed

- [ ] Update `internal/ui/page/interview/editor/` (copy from chat/editor/)
  - [ ] Update package name
  - [ ] Keep editor logic (should work as-is)

- [ ] Update `internal/ui/page/interview/header/` (copy from chat/header/)
  - [ ] Update package name
  - [ ] Change header text to "Interview" context

- [ ] Update `internal/ui/page/interview/sidebar/` (copy from chat/sidebar/)
  - [ ] Update package name
  - [ ] Change "Sessions" → "Interviews" in UI
  - [ ] Remove LSP/MCP status indicators (if any)
  - [ ] Keep file list (for solution files)

- [ ] Update `internal/ui/page/interview/todos/` (copy from chat/todos/)
  - [ ] Update package name
  - [ ] Keep todo rendering (works for interview tasks)

- [ ] Update `internal/ui/page/interview/splash/` (copy from chat/splash/)
  - [ ] Update package name
  - [ ] Change onboarding text for interview context
  - [ ] Remove copilot/hyper OAuth flows
  - [ ] Keep OpenAI/Anthropic API key setup

- [ ] Update `internal/ui/page/page.go` (if exists)
  - [ ] Update page ID constants

### Component Updates

- [ ] Update `internal/ui/tui.go` (copy from `crush/internal/tui/tui.go`)
  - [ ] Update package name
  - [ ] Replace chat page with interview page
  - [ ] Update all references: `chat.ChatPageID` → `interview.InterviewPageID`
  - [ ] Update all references: `cmpChat.*` → `cmpInterview.*`
  - [ ] Update session-related messages → interview-related messages
  - [ ] Update dialogs: sessions dialog → interviews dialog

- [ ] Update `internal/ui/components/dialogs/sessions/` → rename to `interviews/`
  - [ ] Rename package: `sessions` → `interviews`
  - [ ] Update dialog component names
  - [ ] Update references: `Session` → `Interview`
  - [ ] Update UI text

---

## Phase 5: Configuration Cleanup

### Config Simplification

- [ ] Update `internal/config/config.go`
  - [ ] Remove `LSP` type and `LSPs` map
  - [ ] Remove `MCP` type and `MCPs` map
  - [ ] Remove `LSPConfig` struct
  - [ ] Remove `MCPConfig` struct
  - [ ] Remove `projects` related config
  - [ ] Add `InterviewOptions` struct:
    ```go
    type InterviewOptions struct {
        DefaultDifficulty string   `json:"default_difficulty,omitempty"`
        Topics            []string `json:"topics,omitempty"`
        TimeLimit         int      `json:"time_limit,omitempty"`
    }
    ```
  - [ ] Add `InterviewOptions` to `Options` struct
  - [ ] Keep `Providers` section (AI model configuration)
  - [ ] Keep `Models` section
  - [ ] Keep `Permissions` section (simplified)
  - [ ] Keep `TUIOptions` section
  - [ ] Remove `Agent` struct (replace with coach config if needed)

- [ ] Update `internal/config/load.go`
  - [ ] Remove LSP/MCP loading logic
  - [ ] Remove project-related loading
  - [ ] Keep provider loading
  - [ ] Keep model loading

- [ ] Update `internal/config/init.go`
  - [ ] Update default config for interviews
  - [ ] Remove LSP/MCP defaults
  - [ ] Add interview defaults

- [ ] Update `internal/config/provider.go`
  - [ ] Remove `hyper` provider references
  - [ ] Remove `copilot` provider references (keep generic OAuth)
  - [ ] Keep standard providers: OpenAI, Anthropic, Google, etc.

- [ ] Remove `internal/config/copilot.go`
- [ ] Remove `internal/config/hyper.go`
- [ ] Remove `internal/config/lsp_defaults_test.go`

### Provider Configuration

- [ ] Keep `internal/oauth/claude/` (Anthropic OAuth)
- [ ] Keep `internal/oauth/token.go` (generic OAuth token)
- [ ] Remove `internal/oauth/copilot/`
- [ ] Remove `internal/oauth/hyper/`

---

## Phase 6: Tool Cleanup

### Selective Tool Removal

- [ ] Copy `crush/internal/agent/tools/` → `internal/coach/tools/` (selective)

- [ ] **KEEP** (copy these):
  - [ ] `bash.go` - Code execution for interviews
  - [ ] `edit.go` - Code editing
  - [ ] `view.go` - File viewing
  - [ ] `write.go` - File writing
  - [ ] `grep.go` - Code search
  - [ ] `ls.go` - File listing
  - [ ] `glob.go` - File globbing
  - [ ] `todos.go` - Task tracking
  - [ ] `tools.go` - Tool registry/helpers
  - [ ] Template files: `*.md`, `*.tpl` for kept tools

- [ ] **DELETE** (don't copy):
  - [ ] `diagnostics.go` - LSP-specific
  - [ ] `references.go` - LSP-specific
  - [ ] `job_kill.go` - Background jobs
  - [ ] `job_output.go` - Background jobs
  - [ ] `sourcegraph.go` - External code search
  - [ ] `agentic_fetch.go` - Too complex
  - [ ] `agent_tool.go` - Agent-specific
  - [ ] `agentic_fetch_tool.go` - Agent-specific
  - [ ] `mcp/` directory - MCP integration

- [ ] **EVALUATE** (copy but may remove later):
  - [ ] `download.go` - May be useful for downloading problem files
  - [ ] `fetch.go` - May be useful for fetching problem context
  - [ ] `web_search.go` - May be useful if interviews need research
  - [ ] `web_fetch.go` - May be useful if interviews need research
  - [ ] `multiedit.go` - May be useful for batch edits

- [ ] Update `internal/coach/tools/tools.go`
  - [ ] Remove references to deleted tools
  - [ ] Update tool registration
  - [ ] Update package imports

- [ ] Update all tool files:
  - [ ] Change package: `tools` → `tools` (keep same)
  - [ ] Update imports: `crush/internal/session` → `prepf/internal/interview`
  - [ ] Update references: `session.Service` → `interview.Service`
  - [ ] Update tool descriptions for interview context

### Permission Simplification

- [ ] Copy `crush/internal/permission/` → `internal/permission/`
- [ ] Update `internal/permission/permission.go`
  - [ ] Simplify permission system (interviews are less risky)
  - [ ] Remove complex permission dialogs (keep basic tool permissions)
  - [ ] Update tool whitelist/blacklist for interview context

### Shell Package

- [ ] Copy `crush/internal/shell/` → `internal/shell/`
  - [ ] Keep as-is (needed for code execution)
  - [ ] Update package imports if needed

---

## Phase 7: Application Wiring

### App Structure

- [ ] Copy `crush/internal/app/app.go` → `internal/app/app.go`
- [ ] Update `App` struct:
  - [ ] `Sessions` → `Interviews` (type: `interview.Service`)
  - [ ] `AgentCoordinator` → `CoachCoordinator` (type: `coach.CoachCoordinator`)
  - [ ] Remove `LSPClients` field
  - [ ] Keep: `Messages`, `History` (if keeping file history), `Permissions`

- [ ] Update `New()` function:
  - [ ] `session.NewService()` → `interview.NewService()`
  - [ ] Remove `initLSPClients()` call
  - [ ] Remove MCP initialization
  - [ ] `InitCoderAgent()` → `InitInterviewCoach()`
  - [ ] Update coach initialization

- [ ] Create `InitInterviewCoach()` method:
  - [ ] Initialize coach coordinator
  - [ ] Set up interview coach with appropriate prompts
  - [ ] Configure tools for coach

- [ ] Remove `internal/app/lsp.go`
- [ ] Remove `internal/app/lsp_events.go`

- [ ] Update all references throughout app:
  - [ ] `app.Sessions` → `app.Interviews`
  - [ ] `app.AgentCoordinator` → `app.CoachCoordinator`

### Message Service

- [ ] Copy `crush/internal/message/` → `internal/message/`
  - [ ] Update package imports
  - [ ] Keep structure (messages work the same way)
  - [ ] Update references: `session.Session` → `interview.Interview` (if any)

### Database Connection

- [ ] Copy `crush/internal/db/connect.go` → `internal/db/connect.go`
- [ ] Copy `crush/internal/db/db.go` → `internal/db/db.go`
- [ ] Copy `crush/internal/db/embed.go` → `internal/db/embed.go`
- [ ] Update migrations path/embedding
- [ ] Update connection logic for interviews database

### Command Updates

- [ ] Update `internal/cli/root.go` (or create new cmd structure)
  - [ ] Update app name, help text
  - [ ] Remove: projects command, dirs command
  - [ ] Keep: version command
  - [ ] Adapt: login command (remove copilot/hyper, keep OpenAI/Anthropic)

- [ ] Create/update `internal/cli/run.go` (for non-interactive interview mode)
  - [ ] Copy from `crush/internal/cmd/run.go`
  - [ ] Adapt for interview context
  - [ ] Update to use coach instead of agent

- [ ] Create/update login command:
  - [ ] Copy `crush/internal/cmd/login.go` → `internal/cli/login.go`
  - [ ] Remove copilot/hyper login functions
  - [ ] Keep claude/anthropic login
  - [ ] Keep OpenAI API key setup

- [ ] Remove commands (don't copy):
  - [ ] `cmd/projects.go`
  - [ ] `cmd/dirs.go`
  - [ ] `cmd/schema.go` (or adapt if needed)
  - [ ] `cmd/update_providers.go` (or adapt if needed)

### Main Application Entry

- [ ] Update `cmd/prepf/main.go` (or create if needed)
  - [ ] Initialize app with interview service
  - [ ] Set up TUI with interview page
  - [ ] Wire everything together

---

## Phase 8: Cleanup & Testing

### Remove Unused Packages

- [ ] Delete `internal/lsp/` (entire directory)
- [ ] Delete `internal/projects/` (entire directory)
- [ ] Delete `internal/skills/` (entire directory)
- [ ] Remove unused OAuth packages (copilot, hyper)
- [ ] Remove unused config files

### Update Dependencies

- [ ] Review `go.mod`
  - [ ] Keep: `charm.land/bubbletea/v2`, `charm.land/lipgloss/v2`
  - [ ] Keep: `charm.land/fantasy` (AI framework)
  - [ ] Keep: `charm.land/catwalk` (provider management)
  - [ ] Remove: LSP client libraries (if any)
  - [ ] Remove: MCP client libraries
  - [ ] Update module path if needed

### Testing

- [ ] Verify UI compiles and runs
- [ ] Test interview creation
- [ ] Test coach interaction
- [ ] Test message persistence
- [ ] Test tool execution (bash, edit, view, etc.)
- [ ] Test configuration loading
- [ ] Test provider authentication (OpenAI, Anthropic)
- [ ] Verify database migrations work
- [ ] Test interview listing/selection
- [ ] Test non-interactive mode (if implemented)

### Documentation

- [ ] Update README.md with new architecture
- [ ] Update DEVELOPMENT.md with migration notes
- [ ] Document interview-specific configuration
- [ ] Document coach system
- [ ] Update API documentation (if any)

---

## Notes & Considerations

### Key Files to Reference

- `crush/internal/tui/styles/theme.go` - Core styling system
- `crush/internal/config/config.go` - Configuration structure
- `crush/internal/agent/coordinator.go:708-754` - Provider setup (already uses standard providers)
- `crush/internal/session/session.go` - Session model (adapt to Interview)
- `crush/internal/app/app.go:66-118` - Application initialization
- `crush/internal/tui/tui.go:688-710` - TUI initialization

### Important Adaptations

- **Provider Setup**: Already uses standard providers (OpenAI, Anthropic, etc.) - minimal changes needed
- **Session Creation**: Adapt to Interview creation (remove parent session logic)
- **Agent Run**: Adapt prompts for interview coaching context
- **UI Rendering**: Change terminology (Chat → Interview, Session → Interview)

### Migration Order

1. Foundation first (UI/styles/utils) - builds base
2. Data layer (database/interviews) - defines structure
3. Coach system - core logic
4. UI adaptation - user-facing changes
5. Config cleanup - simplification
6. Tool cleanup - remove unused
7. App wiring - connect everything
8. Testing - verify it works

---

**Status**: 🚧 In Progress

**Last Updated**: [Date]

