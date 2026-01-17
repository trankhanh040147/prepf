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

# v0.2.0 - Dual Mode Release (Mock & Gym)

**Status:** Complete

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

- [x] **Mode Selector Component:** Bubbletea list with two options (Mock, Gym)
- [x] **Pre-Mode Input:** Text area for user instructions before entering mode (Use current pre-chat)
- [x] **File Mention (@):** Support `@filepath` syntax to attach CV/resume/context files (implemented in current base)

### Mock Mode (The Gauntlet)

- [x] **System Prompt:** `mock.md.tpl` - Senior Architect persona
  - Penalizes fluff and vague answers
  - Tailors questions based on CV/experience gaps
  - Delivers "The Roast" with actionable feedback
- [x] **Session Flow:**
  - AI asks technical questions based on context
  - User responds (supports "I don't know" shortcut)
  - AI provides harsh but constructive feedback
- [x] **Roast Renderer:** Glamour markdown for feedback display (uses existing message renderer)

### Gym Mode (Training)

- [x] **System Prompt:** `gym.md.tpl` - Drill Instructor persona
  - Generates targeted practice questions
  - Identifies guessing vs. actual knowledge
  - Corrects misconceptions in real-time
- [x] **Session Flow:**
  - User chooses topic or lets AI suggest
  - AI generates random drill questions
  - Interactive Q&A with immediate feedback
- [x] **Topic Selector:** Optional topic input or "surprise me" (handled via user input)

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

# Roadmap Overview

| Version | Codename | Focus | Status |
|---------|----------|-------|--------|
| v0.1.0 | Foundation | Core infrastructure, CLI, TUI, AI client | ✅ Complete |
| v0.2.0 | Dual Modes | Mock & Gym MVP, mode selection | 🚧 In Progress |
| v0.3.0 | Enhanced Feedback | Transcripts, level scaling, progress tracking | 📋 Planned |
| v0.4.0 | Smart Retention | SRS integration, knowledge graphs | 📋 Planned |
| v0.5.0 | Document Intelligence | PDF/Resume parsing, context extraction | 📋 Planned |
| v1.0.0 | Production Ready | Polish, stability, cross-platform release | 📋 Planned |

---

# v0.3.0 - Enhanced Feedback

**Status:** Planned

**Goal:** Improve session value through persistent history, adaptive difficulty, and exportable learning artifacts.

## Features

### Transcript System
- [ ] **Session Export:** JSON/Markdown export of Q&A sessions
- [ ] **Structured Output:** Include question, user answer, AI feedback, knowledge gap tags
- [ ] **Export Commands:** `prepf export <session_id> --format=md|json`
- [ ] **Auto-Save Option:** Config toggle for automatic transcript saving

### Level Scaling
- [ ] **Experience Selector:** Junior → Mid → Senior → Staff → Principal
- [ ] **Adaptive Difficulty:** AI adjusts question complexity based on selected level
- [ ] **Level-Specific Prompts:** Different personas and expectations per level
- [ ] **Profile Integration:** Store preferred level in user profile

### Progress Tracking
- [ ] **Session History:** List past sessions with scores/outcomes
- [ ] **Knowledge Gaps:** Track weak areas across sessions
- [ ] **Progress Dashboard:** TUI view showing improvement over time
- [ ] **Stats Command:** `prepf stats` to view aggregate performance

### Tools & Skills (LLM Enhancement)
- [ ] **Code Execution:** Run code snippets for live technical demonstrations
- [ ] **Diagram Generation:** ASCII/Mermaid diagrams for system design questions
- [ ] **Reference Lookup:** Search documentation during explanations

---

# v0.4.0 - Smart Retention

**Status:** Planned

**Goal:** Implement spaced repetition and knowledge graph features to maximize long-term retention.

## Features

### SRS (Spaced Repetition System)
- [ ] **Anki-Style Algorithm:** SM-2 or FSRS algorithm for optimal review scheduling
- [ ] **Question Bank:** Store questions with retention metadata
- [ ] **Review Mode:** `prepf review` command for daily practice
- [ ] **Difficulty Ratings:** User self-rates difficulty (Easy/Medium/Hard/Again)
- [ ] **Due Queue:** Automatically surface questions due for review

### Knowledge Graph
- [ ] **Topic Mapping:** Track mastery across technical domains
- [ ] **Prerequisite Chains:** Suggest foundational topics when gaps detected
- [ ] **Visualization:** TUI-based topic tree showing coverage
- [ ] **Cross-Domain Links:** Identify connections between technologies

### Misconception Correction
- [ ] **Pattern Detection:** AI identifies recurring misconceptions
- [ ] **Correction Log:** Track corrected vs. persistent misconceptions
- [ ] **Targeted Drills:** Generate questions targeting known weak spots

---

# v0.5.0 - Document Intelligence

**Status:** Planned

**Goal:** Enable rich context extraction from resumes, job descriptions, and technical documents.

## Features

### PDF/Document Parsing
- [ ] **Resume Parser:** Extract experience, skills, projects from PDF/DOCX
- [ ] **Job Description Parser:** Extract requirements, tech stack, level expectations
- [ ] **Structured Extraction:** Convert unstructured docs to queryable data

### Context-Aware Tailoring
- [ ] **Gap Analysis:** Compare resume skills vs. job requirements
- [ ] **Targeted Questions:** Focus on experience gaps automatically
- [ ] **Seniority Detection:** Infer candidate level from experience
- [ ] **Project Deep-Dive:** Generate questions about specific projects listed

### File Management
- [ ] **Context Library:** Store and manage uploaded documents
- [ ] **Document Commands:** `prepf docs add/list/remove`
- [ ] **Session Context:** Attach specific documents to sessions

---

# v1.0.0 - Production Ready

**Status:** Planned

**Goal:** Stable release with polished UX, comprehensive testing, and cross-platform distribution.

## Features

### Polish & Stability
- [ ] **Error Recovery:** Graceful handling of network failures, API limits
- [ ] **Offline Mode:** Basic functionality without internet (review cached questions)
- [ ] **Performance:** Optimize startup time, memory usage
- [ ] **Accessibility:** Screen reader support, high contrast themes

### Distribution
- [ ] **Package Managers:** Homebrew, apt, dnf, pacman, chocolatey
- [ ] **Binary Releases:** Pre-built binaries for Linux/macOS/Windows
- [ ] **Container Image:** Docker image for isolated usage

### Enterprise Features
- [ ] **Team Mode:** Shared question banks, team progress tracking
- [ ] **Custom Prompts:** User-defined AI personas and evaluation criteria
- [ ] **API Mode:** Headless execution for CI/automation

---

# Future Considerations (Post v1.0)

These features are being considered for future releases based on user feedback:

- **Voice Mode:** Audio input/output for realistic interview simulation
- **Collaborative Sessions:** Pair mock interviews with peers
- **Interview Recording:** Record and playback sessions for review
- **Company-Specific Prep:** Question banks tailored to specific companies (FAANG, etc.)
- **Language Support:** Multi-language interview preparation
- **IDE Integration:** VS Code/Neovim extensions for embedded practice
- **Mobile Companion:** Quick review app for on-the-go practice