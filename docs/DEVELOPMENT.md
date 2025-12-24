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

**Status:** In review

**Goal:** Establish core infrastructure following project principles (UX First, Action > Theory).

### Core CLI
- [ ] **Cobra Setup:** Root command with `RunE` pattern, unique flag aliases, error propagation to `main`

### Config Engine
- [x] **Viper Integration:** Config path `~/.config/prepf/config.yaml`, defaults, platform detection, env var overrides
- [x] **Config Command:** View all/specific keys, set values (`prepf config <key> <value>`), edit file (`prepf config edit`)
- [x] **UX:** Fuzzy matching for typos, helpful error messages, config save with validation, initial template
- [x] **Keys:** Writable (api_key, timeout, editor), read-only (no_color, is_tty, config_dir, profile_path)

### TUI Shell (Bubbletea)
- [ ] **Base Model:** Vim keymap (`j/k`, `g/G`, `/`), help overlay (`?`), viewport management, width safety
- [ ] **State Management:** Non-blocking spinner/loading (`tea.Cmd`), window resize, layout via style width
- [ ] **TTY Detection:** `isatty(stdout)` check, `NO_COLOR` support, graceful degradation for piped output

### AI Client (Gemini)
- [x] **Streaming:** Non-blocking updates via `tea.Cmd`, real-time TUI updates, chunk handling with error boundaries
- [ ] **Context:** Token limit handling, usage display, history management
- [ ] **Network:** Configurable timeouts, error wrapping, context cancellation

### Storage Layer
- [ ] **User Profile (JSON):** CV path, Experience Level, file-based persistence, safe mutation, index safety

### Code Quality
- [x] **Structure:** Small modular files, constants in `[...constants.go]`, isolated packages (AI client, stringutil)
- [x] **Standards:** `gofmt -s`, `go vet`, `staticcheck`, complexity ≤15, memory safety, `samber/lo` for functional ops
- [ ] **Dependencies:** Pin versions in `go.mod` ✓, `govulncheck` before commit

### Testing
- [ ] Interactive, Piped (`|`), Redirected (`<`), Non-interactive modes

# v0.1.1 - Mock Module (The Gauntlet)

**Status:** Planned

**Goal:** MVP of "The Gauntlet" - focus on the _Roast_.
- [ ] **Context Loader:** CV/Resume reader (Markdown/PDF-text), tech stack selector (Bubbletea list)
- [ ] **Interview Loop (TUI):** Split view (AI Question top, User Input bottom), "I don't know" shortcut (Tab)
- [ ] **Roast Renderer:** Markdown renderer (Glamour) for harsh feedback, save transcript
- [ ] **System Prompt v1:** "Senior Architect" persona that penalizes fluff

---

# v0.1.2 - Gym Mode

**Status:** In planning

- User chooses topic → AI generates random topics → interactive learning session
