# prepf

[![Latest Release](https://img.shields.io/github/v/release/trankhanh040147/prepf.svg)](https://github.com/trankhanh040147/prepf/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/trankhanh040147/prepf)](https://goreportcard.com/report/github.com/trankhanh040147/prepf)
[![Go Reference](https://pkg.go.dev/badge/github.com/trankhanh040147/prepf.svg)](https://pkg.go.dev/github.com/trankhanh040147/prepf)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **The Gym for your Engineering Career.**
> Stop reading. Start answering. Get roasted. Get improving.

**prepf** is a local CLI tool that acts as your personal technical interview coach. It provides two distinct modes: **Mock Interview (The Gauntlet)** for realistic interview simulation and **Training Mode (The Gym)** for targeted skill practice. Both modes deliver harsh but constructive feedback to help you identify gaps and improve faster.

## Philosophy

**Action over Theory, Retention over Consumption.**
- **Tech Rot is Real:** Focus on adaptability, not static knowledge
- **Kill Tutorial Hell:** Learn by breaking things, not just reading
- **Time Efficiency:** Micro-learning sessions (coffee break length)
- **UX First:** Minimal friction - get in, practice, get feedback, get out
- **Radical Candor:** Objective, harsh feedback to correct misconceptions immediately
- **Breaking Silos:** Demystify adjacent technologies you don't work with daily

## Features

### 🎯 Mock Interview (The Gauntlet)

- **Real-world Simulation:** AI acts as a Senior Architect conducting a technical interview
- **Context-Aware:** Upload your CV/resume to tailor questions to your experience gaps
- **The Roast:** Frank, objective, harsh feedback with actionable advice
- **Penalizes Fluff:** Demands precision and directness in your answers
- **Experience-Based:** Questions adapt to your skill level and background

### 🏋️ Training Mode (The Gym)

- **Targeted Drills:** Practice CS Fundamentals, System Design, or niche tech stacks
- **Misconception Hunter:** AI identifies guessing vs. actual knowledge
- **Real-time Corrections:** Immediate feedback to correct intuition
- **Topic Flexibility:** Choose specific topics or let AI suggest areas to practice
- **Interactive Q&A:** Engaging practice sessions with instant feedback

### 🛠️ Core Capabilities

- **File Attachments:** Use `@filepath` syntax to attach CV/resume/context files
- **Beautiful TUI:** Terminal-first interface built with Bubble Tea
- **Multi-Provider Support:** Works with OpenAI, Anthropic, and other LLM providers
- **Session Management:** Track your progress across multiple practice sessions
- **Privacy-First:** Runs locally, all data stays on your machine
- **Keyboard-First:** Full keyboard navigation with Vim-style bindings

## Prerequisites

- **Go** (version 1.21 or higher)
- **LLM API Key:** OpenAI, Anthropic, or other supported provider
- **Terminal:** A modern terminal with TUI support

## Installation

Install directly using `go install`:

```bash
go install github.com/trankhanh040147/prepf@latest
```

Or build from source:

```bash
git clone https://github.com/trankhanh040147/prepf.git
cd prepf
go build .
```

## Quick Start

1. **Run prepf:**
   ```bash
   prepf
   ```

2. **Select a Mode:**
   - Choose **Mock (The Gauntlet)** for interview simulation
   - Choose **Gym (Training Mode)** for practice drills

3. **Provide Context (Optional):**
   - Type your goals or areas to focus on
   - Attach your CV/resume using `@filepath` syntax
   - Example: `@resume.pdf I want to practice system design questions`

4. **Start Practicing:**
   - Answer questions as they come
   - Receive immediate feedback
   - Learn from "The Roast" in Mock mode or corrections in Gym mode

## Usage

### Basic Commands

```bash
# Run in interactive mode (default)
prepf

# Run with debug logging
prepf -d

# Run in a specific directory
prepf -c /path/to/project

# Run a single non-interactive prompt
prepf run "Explain the use of context in Go"

# Print version
prepf -v
```

### File Attachments

Attach files to provide context for your session:

```bash
# In the pre-mode input, use @ syntax
@resume.pdf @cv.md I want to practice backend system design questions
```

Supported file types: PDF, Markdown, Text files, and code files.

### Session Management

- **New Session:** Press `Ctrl+N` or use the command palette
- **Switch Sessions:** Use `Ctrl+S` to open session selector
- **Session History:** All sessions are saved locally in `~/.config/prepf/`

## Configuration

Configuration is stored in `~/.config/prepf/config.json`. Key settings:

- **LLM Provider:** Configure your preferred model and API keys
- **Model Selection:** Choose between different models per mode
- **Data Directory:** Customize where sessions are stored
- **TUI Preferences:** Adjust interface settings

### Setting Up LLM Provider

1. **OpenAI:**
   ```bash
   export OPENAI_API_KEY="your-key-here"
   ```

2. **Anthropic:**
   ```bash
   export ANTHROPIC_API_KEY="your-key-here"
   ```

3. **Other Providers:**
   See [Configuration Guide](docs/DEVELOPMENT.md) for provider-specific setup.

## Interactive Mode

When running in interactive mode (default), you can:

- **Navigate:** Use Vim-style keys (`j/k` for up/down, `g/G` for top/bottom) or arrow keys
- **Search:** Press `/` to search within content, `n/N` for next/previous match
- **Copy:** Press `y` to copy content to clipboard
- **Help:** Press `?` to see all available keybindings
- **Exit:** Press `q` to quit, `Esc` to exit input mode

## What Gets Reviewed

### Mock Mode
- Technical questions based on your CV/experience
- System design scenarios
- Algorithm and data structure problems
- Language/framework-specific deep dives
- Architecture and best practices

### Gym Mode
- CS fundamentals (algorithms, data structures)
- System design concepts
- Technology-specific drills
- Problem-solving patterns
- Code review and optimization

## Example Session Flow

### Mock Interview Session

```
You: @resume.pdf I'm preparing for a senior backend engineer role

AI: [Senior Architect persona]
    "Looking at your resume, I see 5 years of Go experience but 
    limited distributed systems work. Let's start with: How would 
    you design a rate limiter that handles 1M requests/second?"

You: [Your answer]

AI: [The Roast]
    "Your answer shows you understand the basics, but you missed 
    critical considerations:
    1. Token bucket vs sliding window - you didn't justify your choice
    2. Distributed coordination - how do you handle multiple servers?
    3. ..."
```

### Gym Training Session

```
You: I want to practice Redis caching patterns

AI: [Drill Instructor persona]
    "Question 1: Explain the difference between cache-aside and 
    write-through patterns. When would you use each?"

You: [Your answer]

AI: "Good start, but you're missing a key point about consistency. 
    In write-through, what happens if the cache write fails but the 
    database write succeeds? Think about the failure modes..."
```

## Development

For development information, roadmap, and version-specific context:

- **[Development Roadmap](docs/DEVELOPMENT.md)** - Complete roadmap with all versions, features, and known bugs
- **[Prepf Development Guide](PREPF.md)** - Build, test, and development guidelines

The development documentation includes:
- Design principles and coding standards
- Feature implementation status
- Bug tracking and fixes
- Technical implementation notes
- Related file references

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

Before contributing, please review the [Development Roadmap](docs/DEVELOPMENT.md) to understand the project's direction and design principles.

## Roadmap

### v0.2.0 - Dual Mode Release (Mock & Gym) ✅

- Mode selection dialog
- Mock interview mode with Senior Architect persona
- Training mode with Drill Instructor persona
- File attachment support (`@filepath`)
- Session-based mode tracking

### Upcoming Features

- **Level Scaling:** Junior sanity check → Principal Architect grilling
- **SRS Integration:** Built-in spaced repetition system
- **Topic Suggestions:** Cross-domain topic recommendations
- **Transcript Export:** Save and review your practice sessions
- **PDF Parsing:** Automatic CV/resume parsing

See [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) for the complete roadmap.

## License

MIT License - see [LICENSE](LICENSE) for details.
