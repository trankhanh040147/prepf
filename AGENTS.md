# State of the Art: Go CLI Agents

This codex documents the architectural patterns, libraries, and best practices for building high-performance, interactive AI Agents in Go, based on standout open-source projects in the ecosystem.

## 1. The Stack: "The Charm Ecosystem"

To build a modern, SOTA CLI Agent in Go, the **Charmbracelet** stack is the gold standard. It provides composable, Elm-inspired components for terminal user interfaces (TUIs).

| Library | Purpose | Key Usage in Agents |
| :--- | :--- | :--- |
| **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** | TUI Framework | The core event loop (The Elm Architecture). Handles input, rendering, and signals. |
| **[Lip Gloss](https://github.com/charmbracelet/lipgloss)** | Styling | CSS-like styling for terminal layouts. Essential for readability. |
| **[Bubbles](https://github.com/charmbracelet/bubbles)** | Components | Pre-made UI components: text inputs (editor), viewports, spinners, progress bars. |
| **[Glamour](https://github.com/charmbracelet/glamour)** | Markdown | Rendering rich Markdown (LLM responses) in the terminal. |
| **[Fantasy](https://github.com/charmbracelet/fantasy)** | AI/LLM | A unified interface for LLM providers (OpenAI, Anthropic) and tool execution. |
| **[Log](https://github.com/charmbracelet/log)** | Logging | Structured, colorful logging. |

### Auxiliary Tools
*   **Cobra**: For the CLI command structure (`cmd/`).
*   **SQLite** (`ncruces/go-sqlite3`): For local, zero-config state persistence (sessions, history).
*   **Chroma**: For syntax highlighting in code blocks.

---

## 2. Architecture: The Agent Loop

A SOTA Agent acts as a bridge between the **User (TUI)** and the **LLM (Model)**, mediated by **Tools**.

### Core Flow
1.  **User Input**: Captured via Bubble Tea `Editor` component.
2.  **Context Construction**:
    *   Fetch **Session History** from SQLite.
    *   Load **System Prompts** (persona, constraints).
    *   Attach **File Context** (loaded files, active buffers).
3.  **LLM Inference** (via `fantasy`):
    *   Stream response to TUI (for immediate feedback).
    *   Parse for **Tool Calls** (e.g., `edit`, `ls`).
4.  **Tool Execution**:
    *   **Check Permissions**: Verify if the user allows this action (especially `write`/`exec`).
    *   **Execute**: Run the tool (Go function).
    *   **Capture Output**: stdout/stderr/error.
5.  **Recursion**: Feed tool output back to LLM as a new message. Repeat until completion.

### Directory Structure Reference
Based on common patterns:

```text
internal/
├── agent/              # The Brain
│   ├── agent.go        # Core loop: Context -> LLM -> Tool -> Loop
│   ├── tools/          # Tool implementations (fs, bash, search)
│   └── templates/      # System prompts (markdown templates)
├── tui/                # The Face
│   ├── page/chat/      # Main chat interface (Bubble Tea model)
│   ├── styles/         # Lip Gloss definitions
│   └── components/     # Reusable UI bits (editor, sidebar)
├── session/            # State
│   └── session.go      # Session logic & DB interactions
└── config/             # Configuration
```

---

## 3. Package Architecture & Responsibilities

A robust Go Agent follows a clean separation of concerns. This structure prevents the "God Object" anti-pattern where the TUI or the Agent struct does everything.

### `cmd/` (The Entry Point)
*   **Responsibility**: CLI flag parsing, loading configuration, and initializing the application state.
*   **Best Practice**: Keep it thin. Use **Cobra**. Do not put business logic here. Initialize dependencies (DB, Logger) and pass them to `internal/app`.

### `internal/tui/` (The Presentation Layer)
*   **Responsibility**: Rendering the UI and handling user input (keystrokes).
*   **Key Concept**: The "View" in MVC. It should **never** call the LLM directly. Instead, it emits `tea.Cmd` messages that trigger the Agent.
*   **Sub-packages**:
    *   `pages/`: Full-screen views (Chat, Settings, Onboarding).
    *   `components/`: Reusable widgets (Editor, MarkdownRenderer).
    *   `styles/`: Centralized Lip Gloss styles.

### `internal/agent/` (The Orchestration Layer)
*   **Responsibility**: Managing the conversation loop.
    1.  Receives a prompt from TUI.
    2.  Assembles context (system prompt + history + relevant files).
    3.  Calls the LLM Provider.
    4.  Parses Tool Calls.
    5.  Executes Tools.
    6.  Returns the result (streamed) to the TUI.
*   **Key Concept**: State management. The Agent must track token usage and decide when to summarize history to stay within context windows.

### `internal/tools/` (The Capabilities Layer)
*   **Responsibility**: Providing specific capabilities (File I/O, Web Search, Shell execution).
*   **Contract**: Tools must be **stateless** and **sandboxed**. They receive parameters and return a string (result) or error.
*   **Security**: This layer implements permission checks (e.g., "Allow this tool to edit main.go?").

### `internal/session/` (The Persistence Layer)
*   **Responsibility**: Storing and retrieving chat history.
*   **Key Concept**: Use SQLite. Store messages as structured data (JSON), not just plain text strings, to preserve tool call metadata.

### `internal/config/` (The Configuration Layer)
*   **Responsibility**: Managing user preferences, LLM provider credentials, and model selection.
*   **Best Practice**: support `.env`, config files (YAML), and environment variables with precedence.

---

## 4. Tool Design & Documentation

### Tool Contract (Critical)
Tools are **stateless functions** with clear input/output contracts.
- **Definition**: Each tool lives in `internal/agent/tools/{name}.md` with:
  - Parameter schema (JSON Schema or structured description)
  - Output format (what the LLM sees)
  - Safety constraints (file restrictions, command whitelists)
  - Example invocations
- **Principle**: Tools do ONE thing. They receive params, return results. No guessing, no auto-recovery.
- **Fail Fast**: Return errors immediately. Let the LLM decide next steps, not the tool.
- **Output Strategy**:
  - Structural data (paths, search results): Return raw lists for LLM flexibility.
  - Verbose output (logs, API responses): Truncate/summarize to prevent context bloat.

### Safety & Sandboxing
- **Permissions**: Dangerous operations (`write`, `exec`, `delete`) require explicit user approval.
- **File Restrictions**: Tools cannot write outside project root or access system files.
- **Command Whitelisting**: Shell execution tools have an allowlist of safe commands.
- **Confirmation Workflow**: Pause execution and ask user before risky actions (inline or batch).

### Streaming & Real-Time Updates
- **Channel Pattern**: LLM writes tokens to channel; TUI reads and renders concurrently.
- **Batching**: Collect tokens until punctuation/newline (avoid excessive redraws).
- **Cancellation**: Pass `context.Context` through entire stack to enable Ctrl+C interrupts.
- **Cursor**: Show spinner while streaming to indicate "thinking" state.

### Testing & Mocking
- **Mock Providers**: Test agent logic without API calls. Use interfaces for LLM providers.
- **Golden Files**: Store expected outputs for regression testing (`.golden` files).
- **Debug Logging**: Log to file (`debug.log`) instead of stdout to preserve TUI rendering.
- **Deterministic Behavior**: Same inputs must produce same outputs (enables reliable testing).

---

## 5. Good vs. Bad Practices

### Architecture & Concurrency

| Feature | ❌ Bad Practice | ✅ Good Practice |
| :--- | :--- | :--- |
| **Blocking** | Calling LLM APIs directly in the Bubble Tea `Update()` function. This freezes the UI. | Return a `tea.Cmd` that runs the API call in a goroutine and returns a `Msg` when done. |
| **State** | Passing the entire `App` struct to every function. | Pass specific interfaces or context (e.g., `SessionService`, `Agent`) to components that need them. |
| **Cancellation** | Ignoring user interrupts. The agent keeps running after `Ctrl+C`. | Pass `context.Context` through the entire stack. Cancel the context when the user hits Esc/Ctrl+C. |

### Prompt Engineering & Context

| Feature | ❌ Bad Practice | ✅ Good Practice |
| :--- | :--- | :--- |
| **Prompts** | Hardcoding system prompts as Go string constants. | Store prompts as `embed` Markdown files (`templates/system.md`). Allows easy editing and versioning. |
| **Context** | Dumping the entire project file tree into the context. | Use **RAG** (Retrieval Augmented Generation) or explicit user "mentions" (@file) to select relevant context. |
| **Summarization**| Crashing when the conversation gets too long (Context Limit Exceeded). | Implement "Rolling Window" or "Auto-Summarization" when token count > 80% of window. |

---

## 6. Key Patterns & Best Practices

### A. TUI Integration (Bubble Tea)
*   **Streaming**: Don't wait for the full response. Update the UI on every token chunk.
    *   Use `tea.Batch` to handle UI updates and channel reading concurrently.
*   **State Machines**: Use `tea.Model` to manage states: `Idle`, `Thinking`, `ExecutingTool`, `StreamingResponse`.
*   **Components**:
    *   `Sidebar`: For session history and file trees.
    *   `Viewport`: For scrolling large chat logs.
    *   `Editor`: For multiline user input (often with syntax highlighting).

### B. Context Management (Token Budgeting & RAG)
*   **Token Budgeting**: Track per-request usage. When approaching limit (>80%), apply rolling window or auto-summarize.
    *   Strategy: Keep system prompt + last N messages + file references within budget.
    *   Model-specific: Claude 3.5 = 200K; GPT-4 = 128K. Configure per provider.
*   **RAG (Retrieval-Augmented Gen.)**: Don't embed entire codebases. Let LLM request files via explicit commands (e.g., `@file` syntax).
*   **System Prompts**: Embed as Markdown in `internal/agent/templates/`. Load at runtime for easy versioning.

### C. Persistence (SQLite)
*   **Zero Config**: Use `sqlite3` for a single-file database.
*   **Sessions**: Store `messages` (role, content, tool_calls) and `sessions` (metadata).
*   **FTS (Full Text Search)**: Enable FTS5 for searching past conversations.

---

## 7. Developer Experience (DX)

*   **Debug Logging**: Use `charm.land/log`.
    *   Log to a file (`debug.log`) instead of stdout to avoid messing up the TUI.
*   **Mocks**: Use interfaces for LLM providers to test the UI without burning API credits.

## 8. The Standard
To adhere to the standard set by modern CLI agents:
1.  **Beauty Matters**: A CLI tool should look as good as a web app. Use `lipgloss` extensively.
2.  **Keyboard First**: Full keyboard navigation (Vim bindings are a plus).
3.  **Transparency**: Show what the agent is "thinking" (reasoning traces) and what tools it is calling.
4.  **Local First**: Prefer local files and local DBs.
