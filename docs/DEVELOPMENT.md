# Project: prepf

## Design Principles & Coding Standards

> **Reference:** All design principles, coding standards, and implementation guidelines are defined in [`.cursor/rules/rules.mdc`](../.cursor/rules/rules.mdc).

### How To Apply These Rules

Automatically loads rules from the `.cursor/rules/` directory. The `rules.mdc` file includes `alwaysApply: true` in its frontmatter, which ensures:

- **Automatic Application:** Rules are always active during coding sessions
- **Context Awareness:** Understands project-specific patterns (Vim navigation, TUI-first UX, Go conventions)
- **Consistency:** All code suggestions follow the defined principles without manual reminders

## Bug Fix Protocol

1. **Global Fix:** Search codebase (`rg`/`fd`) for similar patterns/implementations. Fix **all** occurrences, not just the reported one.
2. **Documentation:**
    - Update "Known Bugs" table (Status: Fixed).
    - Update coding standards in `.cursor/rules/rules.mdc` if the bug reflects a common anti-pattern.
3. **Testing:** Verify edge cases: Interactive, Piped (`|`), Redirected (`<`), and Non-interactive modes.
> **Reference:** Bug Fix Protocol are defined in [`.cursor/rules/rules.mdc`](../.cursor/rules/rules.mdc).

# Philosophy

**Action over Theory, Retention over Consumption.**
- **Tech Rot is Real:** Focus on adaptability, not static knowledge
- **Kill Tutorial Hell:** Learn by breaking things, not just reading
- **Time Efficiency:** Micro-learning sessions (coffee break length)
- **UX First:** Minimal friction - get in, practice, get feedback, get out
- **Radical Candor:** Objective, harsh feedback to correct misconceptions immediately
- **Breaking Silos:** Demystify adjacent technologies you don't work with daily

# Core Features

**1. The Gauntlet (Mock Interview)**
- Real-world simulation with AI contexts specific to your stack
- Context-aware: CV/Resume upload tailors questions to experience gaps
- The "Roast": Frank, objective, harsh feedback with actionable advice
- Level scaling: Junior sanity check → Principal Architect grilling

**2. The Gym (Training Mode)**
- Targeted drills: CS Fundamentals, System Design, niche tech stacks
- Cross-domain fluency: Fast-track learning for seniors moving to new stacks
- Smart retention: Built-in SRS (Spaced Repetition System) like Anki
- Misconception Hunter: AI identifies guessing vs. actual knowledge, corrects intuition in real-time

---


# v0.1.0 - Foundation

**Status:** ✅ Complete

**Goal:** Establish core infrastructure following project principles (UX First, Action > Theory).

## What's New

- **CLI Flags:** Added `--config/-c`, `--profile/-p`, `--verbose/-v`, `--quiet/-q` with path override support
- **AI Client:** Token limit management (1M default), cumulative usage tracking, `UsageDisplay()` method, conversation history for multi-turn interactions
- **Storage:** Profile validation (experience levels), safe mutation with ID caching, atomic saves
- **TUI:** BaseModel helpers (`Viewport()`, `SetViewportContent()`), verified width safety and non-blocking patterns
- **Testing:** Comprehensive coverage for interactive TUI, piped mode, non-interactive CLI, and AI client features
- **CI/CD:** GitHub Actions workflow with testing, linting (`gofmt`, `go vet`, `staticcheck`), vulnerability scanning (`govulncheck`), and multi-platform builds

## Code Quality Fixes

- **Config Loading:** Refactored to use context injection pattern. Config loaded in `rootCmd.PersistentPreRunE`, injected via context. Removed global flag getters.
- **Viper State:** Documented single-initialization requirement for Viper global state.
- **Validation:** Removed I/O (`os.Stat`) from `Profile.Validate()`, deferred to actual usage.
- **Token Estimation:** Added 20% safety margin to token estimation heuristic, documented limitation.
- **Code Refactoring:** Extracted `UsageStats` struct, refactored `SafeUpdate` duplication, extracted `ValidExperienceLevels` constant, simplified `InitialConfigContent` to use raw string literal.

### Core CLI
- [x] **Cobra Setup:** Root command with `RunE` pattern, unique flag aliases, error propagation to `main`

### Config Engine
- [x] **Viper Integration:** Config path `~/.config/prepf/config.yaml`, defaults, platform detection, env var overrides
- [x] **Config Command:** View all/specific keys, set values (`prepf config <key> <value>`), edit file (`prepf config edit`)
- [x] **UX:** Fuzzy matching for typos, helpful error messages, config save with validation, initial template
- [x] **Keys:** Writable (api_key, timeout, editor, token_limit), read-only (no_color, is_tty, config_dir, profile_path)

### TUI Shell (Bubbletea)
- [x] **Base Model:** Vim keymap (`j/k`, `g/G`, `/`), help overlay (`?`), viewport management, width safety
- [x] **State Management:** Non-blocking spinner/loading (`tea.Cmd`), window resize, layout via style width
- [x] **TTY Detection:** `isatty(stdout)` check, `NO_COLOR` support, graceful degradation for piped output

### AI Client (Gemini)
- [x] **Streaming:** Non-blocking updates via `tea.Cmd`, real-time TUI updates, chunk handling with error boundaries
- [x] **Context:** Token limit handling, usage display, history management
- [x] **Network:** Configurable timeouts, error wrapping, context cancellation

### Storage Layer
- [x] **User Profile (JSON):** CV path, Experience Level, file-based persistence, safe mutation, index safety

### Code Quality
- [x] **Structure:** Small modular files, constants in `[...constants.go]`, isolated packages (AI client, stringutil)
- [x] **Standards:** `gofmt -s`, `go vet`, `staticcheck`, complexity ≤15, memory safety, `samber/lo` for functional ops
- [x] **Dependencies:** Pin versions in `go.mod` ✓, `govulncheck` before commit

### Testing
- [x] Interactive, Piped (`|`), Redirected (`<`), Non-interactive modes

### CI/CD 
- [x] GitHub Actions workflow with testing, linting, vulnerability scanning, and multi-platform builds


# v0.1.1 - Mock Module (The Gauntlet MVP)

**Goal:** Refactor and improve code quality of existing Mock Module. No new features - pure refactoring release.

## Code Quality Fixes
**Status:** ✅ Complete

- **Viewport Padding:** Fixed manual padding anti-pattern, use `lipgloss.Style.Padding()`
- **Markdown Rendering:** Added markdown rendering with `glamour`
- **Slice Mutation Bug:** Fixed slice mutation bug: pass `*[]string` pointers to form
- **Enter/Tab Navigation:** Fixed form key handling: state-specific keys before global keys
- **Redundant Data Flow:** Simplified `ConfigSubmittedMsg` to empty struct
- **Global Quit Keys:** Fixed Ctrl+C/Esc handling across all states
- **Input State Escape:** Added explicit Esc handling to blur textinput before quit checks
- **Form Navigation Clarity:** Updated form descriptions for Space/Tab/Enter behavior
- **Roast Display:** Fixed micro-roast disappearing, added `showSurrenderFeedback` flag
- **Final Roast Generation:** Implemented AI-generated roast feedback with grade-based assessment
- **State-Specific Key Handling:** Established precedence: state intercepts → component updates → global fallbacks
- **Context Cancellation:** Fixed no-op cancel function: use `context.WithCancel()`
- **errgroup Pattern:** Removed detached goroutine pattern that could hide errors
- **Stream Channel Consolidation:** Replaced three channels with single `streamMsgChan`
- **Clipboard Implementation:** Added clipboard support (`yy`/`Y`) using `atotto/clipboard`
- **Dead Code Removal:** Removed unused stream command functions

## Refactor codebase 

**Status:** ✅ In Progress - Phase 1 Complete

### What We've Done

**Phase 1: Clone Crush Internal Packages** ✅ Complete

- **Cloned Utility Packages:** Successfully cloned and integrated reusable utility packages from `crush/internal/`:
  - `ansiext/` - ANSI escape utilities
  - `csync/` - Thread-safe concurrent slices and maps
  - `diff/` - Diff generation utilities
  - `env/` - Environment variable abstraction
  - `filepathext/` - Filepath extensions (`SmartJoin`, `SmartIsAbs`)
  - `fsext/` - File system utilities (fileutil, ignore, lookup, ls, expand, owner)
  - `format/` - Formatting utilities (spinner wrapper)
  - `home/` - Home directory utilities
  - `stringext/` - String extensions
  - `log/` - Logging setup (JSON handler, rotation, panic recovery)
  - `version/` - Version management utilities

- **Cloned TUI Components:** Cloned reusable TUI components from `crush/internal/tui/`:
  - `components/core/` - Core layout/status components
  - `components/anim/` - Animation utilities
  - `components/completions/` - Completions UI
  - `exp/diffview/` - Diff view component
  - `exp/list/` - List components
  - `highlight/` - Syntax highlighting
  - `styles/` - Styling utilities
  - `util/` - TUI utilities
  - Excluded crush-specific: `chat/`, `dialogs/`, `lsp/`, `mcp/`, `page/`

- **Import Adaptation:** Updated all imports from `github.com/charmbracelet/crush` to `github.com/trankhanh040147/prepf`
- **Dependency Management:** Added missing dependencies (`mvdan.cc/sh/v3`, `gopkg.in/natefinch/lumberjack.v2`, `github.com/zeebo/xxh3`)
- **Stub Packages:** Created compatibility stubs for missing dependencies:
  - `internal/uiutil/` - UI utility functions for TUI components
  - `internal/history/` - File history stub for file components
- **Build Status:** Main prepf command (`./cmd/prepf`) builds successfully ✅

### What We Need To Do Next

**Phase 2: Integration & Compatibility** 🔄 Next Steps

- [ ] **Fix TUI Component Compatibility:** Some cloned TUI components have API version mismatches:
  - `internal/tui/components/anim/` - Color API incompatibility with older lipgloss version
  - `internal/tui/styles/` - References `charm.land/glamour/v2` which doesn't exist
  - `internal/tui/tui.go` - References crush-specific packages (`app`, `event`, `permission`, `pubsub`, `agent/tools/mcp`)
  - **Action:** Either update prepf's bubbletea/lipgloss versions OR adapt components to work with current versions

- [ ] **Integrate Utility Packages:** Start using cloned utilities in prepf codebase:
  - Replace manual string operations with `stringext/` functions
  - Use `fsext/` for file operations instead of manual `filepath` usage
  - Integrate `csync/` for concurrent data structures where needed
  - Use `diff/` for any diff generation features
  - Consider `log/` package for structured logging

- [ ] **Merge String Utilities:** Merge `internal/util/stringutil/fuzzy.go` with `internal/stringext/` (prefer crush's naming)

- [ ] **TUI Component Integration:** Evaluate which TUI components from crush are useful:
  - `exp/list/` - Could replace custom list implementations
  - `exp/diffview/` - Useful for showing code diffs in mock interviews
  - `components/completions/` - Could enhance CLI autocomplete
  - `components/core/status/` - Already used by `components/files/`

- [ ] **Remove Unused Components:** Clean up components that won't be used:
  - `internal/tui/tui.go` - Crush-specific app model (not needed, prepf has own UI)
  - `internal/tui/components/files/` - Only used by crush's chat component
  - Any other crush-specific dependencies

- [ ] **Documentation:** Update codebase documentation to reflect new package structure

**Phase 3: Code Quality** 📋 Future

- [ ] **Adopt Library Patterns:** Replace custom implementations with crush's utilities where beneficial
- [ ] **Consistency:** Ensure all code follows crush's patterns for utilities (DRY, functional helpers)
- [ ] **Testing:** Add tests for integrated utility packages

## New Features 

**Status:** In planning

- [ ] **Pre-mock Topic Customization:** Added skippable configuration screen using `huh` forms for selecting topics to focus on and exclude. Users can choose from: Go, System Design, Algorithms, Data Structures, Databases, Networking, Concurrency, Testing. Press `Esc` to skip configuration.
- [ ] **Enhanced Prompts:** Improved AI prompts to generate realistic, varied interview questions. Prompts now include:
  - Instructions to ask questions that real interviewers would ask
  - Variety guidance to avoid repetitive questions
  - Conversation history awareness for question diversity
  - Topic-specific instructions based on user selections
- [ ] **Conversation History Integration:** Added variety instructions to each answer submission to ensure AI references conversation history and avoids repetition.

### 1. Sequential Interview Engine
- [ ] **Turn-Based Flow:** Strictly one question at a time. User input is locked while AI "speaks."
- [ ] **AI Orchestration:** AI decides when to follow up on an answer or pivot to a new topic via hidden `<NEXT>` signals.
- [ ] **Context Loader (v0):** Support for `.txt` and `.md` resume ingestion via `os.ReadFile`.
- [ ] **Protocol Engine:** `regexp` parser to intercept hidden `<NEXT>` and `<ROAST>` signals.

### 2. The "Roast" Mechanics
- [ ] **The "Surrender" Mechanic:** `Tab` key injects a Shadow Prompt: *"User surrenders. Give a snappy 1-2 sentence correction and move on."*
- [ ] **Inline Micro-Roast UI:** Mid-interview failures/surrenders styled in **Bold Red** via `lipgloss` for immediate feedback.
- [ ] **The Verdict:** 
    - **Visual Grade:** High-contrast `lipgloss` box displaying **Letter Grade (A-F)**.
    - **Persona Labels:** Descriptive status (e.g., `[A] - ARCHITECT MATERIAL`, `[F] - TERMINATED`).
- [ ] **The Roast:** AI-generated assessment with 3-point remediation plan rendered as **Interactive Buttons** (Placeholders for Gym Mode).

### 3. Session Governance
- [ ] **Safety Valve:** Hard limit (10 questions/15 mins).
- [ ] **Graceful Exit:** Status bar triggers an **Inverted Pulsing [FINAL QUESTION] Alert** (`tea.Tick`); system forces `<ROAST>` after the current turn.
- [ ] **Metadata Tracking:** Silently track "Surrender" count for future grading logic.

# v0.1.2 - The Scalable Standard
**Status:** In planning (Structured Logic & Persistence)

- [ ] **Function Calling:** Migrate to Gemini Tool Use for state transitions (`pivot_topic`, `finalize_roast`).
- [ ] **PDF Support:** Native resume parsing via `ledongthuc/pdf`.
- [ ] **Persistence:** Save transcripts to `~/.local/share/prepf/history/` as JSON.
- [ ] **Gym Integration:** Activate "Remediation Buttons" to launch targeted sessions in **The Gym**.

---

## Future Considerations (v0.2.0+)
- **Hybrid Grading:** Implement "Grade Ceilings" based on the number of `Tab` surrenders.
- **Dynamic Pressure:** TUI border colors shift (Green → Red) based on "Roast" severity.
- **Knowledge Graph:** Real-time side-panel visualization of technical "Weak Spots."
- **Gym Mode:** User chooses topic → AI generates random topics → interactive learning session
