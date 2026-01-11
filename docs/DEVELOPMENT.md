# Project: prepf

## Design Principles & Coding Standards

> **Reference:** All design principles, coding standards, and implementation guidelines are defined in [`.cursor/rules/rules.mdc`](../.cursor/rules/rules.mdc) and [`AGENTS.md`](../AGENTS.md)

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

### Code Quality
- [x] **Structure:** Small modular files, constants in `[...constants.go]`, isolated packages (AI client, stringutil)
- [x] **Standards:** `gofmt -s`, `go vet`, `staticcheck`, complexity ≤15, memory safety, `samber/lo` for functional ops
- [x] **Dependencies:** Pin versions in `go.mod` ✓, `govulncheck` before commit

### Testing
- [x] Interactive, Piped (`|`), Redirected (`<`), Non-interactive modes

### CI/CD 
- [x] GitHub Actions workflow with testing, linting, vulnerability scanning, and multi-platform builds

# v0.1.1 - Dual Mode Release (Mock & Gym)

**Status:** In Progress

**Goal:** Experimental release to validate the tool's capability. Both modes are MVP-level to gather feedback before enhancing in future releases.

## User Flow

```
prepf → Mode Selection → Pre-Mode Instructions → Enter Mode
        ├── Mock (The Gauntlet)
        └── Gym (Training Mode)
```

1. User runs `prepf`
2. Mode selection dialog (Bubbletea list): **Mock** or **Gym**
3. Pre-mode instruction input: User provides context/goals before entering
4. Mode-specific session begins

## Features

### Mode Selection (Shared)

- [ ] **Mode Selector Component:** Bubbletea list with two options (Mock, Gym)
- [x] **Pre-Mode Input:** Text area for user instructions before entering mode (Use current pre-chat)
- [x] **File Mention (@):** Support `@filepath` syntax to attach CV/resume/context files (implemented in current base)

### Mock Mode (The Gauntlet)

- [ ] **System Prompt:** `mock.md.tpl` - Senior Architect persona
  - Penalizes fluff and vague answers
  - Tailors questions based on CV/experience gaps
  - Delivers "The Roast" with actionable feedback
- [ ] **Session Flow:**
  - AI asks technical questions based on context
  - User responds (supports "I don't know" shortcut)
  - AI provides harsh but constructive feedback
- [ ] **Roast Renderer:** Glamour markdown for feedback display

### Gym Mode (Training)

- [ ] **System Prompt:** `gym.md.tpl` - Drill Instructor persona
  - Generates targeted practice questions
  - Identifies guessing vs. actual knowledge
  - Corrects misconceptions in real-time
- [ ] **Session Flow:**
  - User chooses topic or lets AI suggest
  - AI generates random drill questions
  - Interactive Q&A with immediate feedback
- [ ] **Topic Selector:** Optional topic input or "surprise me"

## Prompt Templates

Create new templates in `internal/agent/templates/`:

| Template | Purpose |
|----------|---------|
| `mock.md.tpl` | The Gauntlet - Mock interview system prompt |
| `gym.md.tpl` | The Gym - Training mode system prompt |
| `mode_select.md` | Mode selection instructions |

## Implementation Notes

- **Minimal Changes:** This release experiments with capability, not polish
- **Leverage Existing:** Use existing TUI components (chat, viewport, editor)
- **Mode State:** Track current mode in session metadata

## Out of Scope (Future Releases)

- PDF parsing (use markdown/text CV for now)
- SRS (Spaced Repetition) integration
- Transcript saving
- Level scaling (Junior → Principal)
- Cross-domain topic suggestions

---

# v0.1.2 - Enhanced Feedback

**Status:** Planned

- Tools & Skills: Implement Tools and Skills to grant LLM more capabilities

- Transcript saving (JSON/Markdown export)
- Level scaling selector (Junior → Senior → Principal)
- Session history and progress tracking