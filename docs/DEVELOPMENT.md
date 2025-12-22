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

# Philosophy (The "Why")
**pref** prioritizes **Action over Theory** and **Retention over Consumption**.
- **Tech Rot is Real:** Technology moves too fast to memorize everything. We focus on adaptability, not just static knowledge.
- **Kill Tutorial Hell:** You can't learn just by reading docs. You learn by breaking things and getting corrected.
- **Time Efficiency:** Designed for the busy engineer. Micro-learning sessions that fit into a coffee break, not hour-long lectures.
- **UX First:** Minimal friction. Get in, practice, get feedback, get out.
- **Radical Candor:** Most feedback is too nice. We believe in objective, harsh data to correct misconceptions immediately.
- **Breaking Silos:** We reject "that's not my job." A Backend Engineer must understand Frontend constraints. We demystify the "black boxes" of adjacent technologies you don't work with daily.
 Help developer didn understand all aspects of other technologies that they haven't worked with yet

# Core features
**1. The Gauntlet (Mock Interview)**
- **Real-world Simulation:** Interactive, voice/text-based interviews powered by AI contexts specific to your stack (e.g., Golang, Blockchain).
- **Context Aware:** Upload your CV/Resume. The AI tailors questions to your actual experience gaps.
- **The "Roast" (Feedback):** Post-interview reports that are frank, objective, and harsh. No sugar-coating—just distinct actionable advice on where you failed.
- **Level Scaling:** From "Junior sanity check" to "Principal Architect grilling."

**2. The Gym (Training Mode)**
- **Targeted Drills:** focused modules on CS Fundamentals, System Design, and niche tech stacks.
- **Cross-Domain Fluency:** fast-track learning for seniors moving to new stacks (e.g., a Backend dev learning Frontend concepts quickly).
- **Smart Retention:** Built-in **Spaced Repetition System (SRS)** (like Anki) to ensure you don't forget the edge cases of a language.
- **Misconception Hunter:** The AI actively identifies when you are guessing versus when you actually know the answer, correcting your intuition in real-time.

---


# v0.1.0 - Foundation

**Status:** In Progress

**Goal:** Establish core infrastructure following project principles (UX First, Action > Theory).

### Core CLI
- [ ] **Cobra Setup:**
    - [ ] Root command `prepf` with `RunE` pattern (no `os.Exit` in library code)
    - [ ] Unique short flag aliases per command (verify no redefinitions)
    - [ ] Error propagation to `main` only

### Config Engine
- [ ] **Viper Integration:**
    - [ ] Config path: `~/.config/prepf/config.yaml` (use `os.MkdirAll` for parent dirs)
    - [ ] Sensible defaults (API keys, timeouts, editor path)
    - [ ] Platform detection: `runtime.GOOS` for editor/open commands
    - [ ] Environment variable overrides

### TUI Shell (Bubbletea)
- [ ] **Base Model:**
    - [ ] Global keymap (Vim: `j/k`, `g/G`, `/`)
    - [ ] Help overlay (`?` key) with centralized constants
    - [ ] Viewport management (responsive width, default 80 if `Width() == 0`)
    - [ ] Width safety: guard `strings.Repeat` with `max(0, count)`, prevent negative widths
- [ ] **State Management:**
    - [ ] Spinner/loading states (non-blocking, `tea.Cmd` pattern)
    - [ ] Window resize handling (`WindowSizeMsg`)
    - [ ] Layout: center via style width only (never calculate string length)
- [ ] **TTY Detection:**
    - [ ] Check `isatty(stdout)` before colors/spinners
    - [ ] Respect `NO_COLOR` env var
    - [ ] Graceful degradation for piped/redirected output

### AI Client (Gemini)
- [ ] **Streaming Handler:**
    - [ ] Non-blocking updates via `tea.Cmd` (never block `Update` loop)
    - [ ] Real-time TUI updates during stream
    - [ ] Chunk handling with explicit error boundaries
- [ ] **Context Management:**
    - [ ] Token limit handling with clear errors
    - [ ] Token usage display in TUI
    - [ ] History management (conversation state)
- [ ] **Network Safety:**
    - [ ] Mandatory default timeouts (configurable via config)
    - [ ] Error wrapping: `fmt.Errorf("context: %w", err)`
    - [ ] Context cancellation propagation

### Storage Layer
- [ ] **User Profile (JSON):**
    - [ ] Simple JSON store (CV path, Experience Level)
    - [ ] File-based persistence
    - [ ] Safe mutation: cache IDs before editing
    - [ ] Index safety: verify bounds before access

### Code Quality
- [ ] **Structure:**
    - [ ] Small, modular files
    - [ ] Constants in `[...constants.go]` files (no hardcoding)
    - [ ] AI client isolation (separate package)
- [ ] **Standards:**
    - [ ] Run `gofmt -s`, `go vet`, `staticcheck`
    - [ ] Target complexity ≤15 per function
    - [ ] Memory safety: never `&m[k]`, assign to var first
- [ ] **Dependencies:**
    - [ ] Pin major versions in `go.mod`
    - [ ] Run `govulncheck` before commit

### Testing Requirements
- [ ] Verify in Interactive mode
- [ ] Verify in Piped mode (`|`)
- [ ] Verify in Redirected mode (`<`)
- [ ] Verify non-interactive mode (flags only)

---


# v0.1.1 - Mock Module (The Gauntlet)

**Status**: Planned, raw ideas

**Goal:** The MVP of "The Gauntlet." Focus on the _Roast_.
- [ ] **Context Loader:**
    - [ ] File reader (`os.ReadFile`) for CV/Resume (Markdown/PDF-text).
    - [ ] Tech Stack selector (Bubbletea list).
- [ ] **The Interview Loop (TUI):**
    - [ ] Split view: AI Question (top) vs. User Input (bottom textarea).
    - [ ] "I don't know" shortcut (Tab) to encourage honesty vs. guessing.
- [ ] **The Roast Renderer:**
    - [ ] Markdown renderer (Glamour) to display the "Frank/Harsh" feedback prettily.
    - [ ] Save transcript to local file.
- [ ] **System Prompt v1:** Implement the "Senior Architect" persona that penalizes fluff.

---

# v0.1.2

**Status**: In planning

- Gym mode: User choose what to learn --> AI generate random topics to choose --> start interactive learning
