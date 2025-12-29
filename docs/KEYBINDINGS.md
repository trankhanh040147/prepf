# Keybinding Rules & Guidelines

## Global Keybinding Principles

### 1. **Precedence Order**
All TUI states must handle keys in this order:
1. **State-specific intercepts** (e.g., Tab for surrender in interview)
2. **Component updates** (e.g., textinput, form processing)
3. **Global fallbacks** (quit, help)

### 2. **Global Keys (Must Work Everywhere)**
- `q` / `Ctrl+C` → Quit application
- `?` → Toggle help overlay
- `Esc` → Back/cancel/blur (context-dependent)

### 3. **Focused Component Handling**
When components like `textinput` or `huh.Form` are focused:
- They consume most keypresses
- **MUST** explicitly handle `Esc` and `Ctrl+C` BEFORE passing to component
- Check `component.Focused()` before handling quit keys
- Example:
  ```go
  case InterviewUserInput:
      // Handle Esc FIRST
      if key.Matches(msg, m.keys.Back) {
          if m.answerInput.Focused() {
              m.answerInput.Blur()
              return m, nil
          }
      }
      // Handle Ctrl+C BEFORE component
      if key.Matches(msg, m.keys.Quit) {
          m.cancelCtx()
          return m, tea.Quit
      }
      // NOW pass to component
      var cmd tea.Cmd
      m.answerInput, cmd = m.answerInput.Update(msg)
      return m, cmd
  ```

### 4. **Form Navigation (`huh.Form`)**
- Forms handle Enter/Tab/Space internally
- Let form process ALL keys except global quit/help
- Form completes when user presses Enter on **last field**
- Tab navigates between fields
- Space toggles multi-select items
- **Never** intercept Enter/Tab when form is active

### 5. **Transient Content Display**
When showing temporary content (micro-roasts, notifications):
- Use persistent **flags** (e.g., `showSurrenderFeedback`)
- Clear flag when **new content arrives**, not on first render
- Prevents content from disappearing immediately

## Mock Interview Keybindings

### Configuration State
- `↑/↓` → Navigate options
- `Space` → Toggle selection
- `Tab` → Move to next field
- `Enter` → Submit (on last field only)
- `Esc` → Skip configuration
- `q/Ctrl+C` → Quit

### User Input State
- `Enter` → Submit answer
- `Tab` → Surrender question
- `Esc` → Blur input (first press) / Back (when blurred)
- `q/Ctrl+C` → Quit
- `?` → Help

### Roasting State
- `q/Ctrl+C` → Quit
- `?` → Help

## Anti-Patterns to Avoid

### ❌ **DON'T: Let components block quit keys**
```go
// BAD: Input consumes Esc, user can't quit
case tea.KeyMsg:
    m.input, cmd = m.input.Update(msg)
    return m, cmd
```

### ✅ **DO: Handle critical keys first**
```go
// GOOD: Check quit keys before component
case tea.KeyMsg:
    if key.Matches(msg, m.keys.Quit) {
        return m, tea.Quit
    }
    m.input, cmd = m.input.Update(msg)
    return m, cmd
```

### ❌ **DON'T: Clear transient content immediately**
```go
// BAD: Micro-roast disappears on first render
if m.microRoast != "" {
    display := RenderMicroRoast(m.microRoast)
    m.microRoast = "" // Cleared too early!
    return display
}
```

### ✅ **DO: Use persistent flags**
```go
// GOOD: Content stays until explicitly cleared
if m.showMicroRoast && m.microRoast != "" {
    display := RenderMicroRoast(m.microRoast)
    return display
}
// Clear flag when new content arrives
m.showMicroRoast = false
```

### ❌ **DON'T: Intercept form navigation keys**
```go
// BAD: Breaks form navigation
case InterviewConfiguring:
    if key.Matches(msg, m.keys.Enter) {
        return m.handleSubmit() // Form never gets Enter!
    }
    m.form.Update(msg)
```

### ✅ **DO: Let form handle its own keys**
```go
// GOOD: Form processes Enter/Tab/Space
case InterviewConfiguring:
    if key.Matches(msg, m.keys.Quit) {
        return m, tea.Quit
    }
    // Form handles Enter/Tab internally
    m.form, cmd = m.form.Update(msg)
    if m.form.State == huh.StateCompleted {
        return m.handleSubmit()
    }
    return m, cmd
```

## Testing Checklist

When implementing new keybindings:
- [ ] `q` works in all states
- [ ] `Ctrl+C` works in all states
- [ ] `Esc` works when input is focused
- [ ] `?` toggles help in all states
- [ ] Forms can navigate with Tab
- [ ] Forms can submit with Enter
- [ ] Transient content stays visible long enough to read
- [ ] No key conflicts between states
- [ ] Help text matches actual keybindings

## Reference Implementation

See `internal/mock/update.go` for complete implementation of these principles.
