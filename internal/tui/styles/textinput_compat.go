package styles

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/lipgloss"
)

// TextInputStyles is a compatibility type for textinput styling
// that matches the structure expected by cloned crush components
type TextInputStyles struct {
	Focused TextInputStyleState
	Blurred TextInputStyleState
	Cursor  TextInputCursorStyle
}

type TextInputStyleState struct {
	Text        lipgloss.Style
	Placeholder lipgloss.Style
	Prompt      lipgloss.Style
	Suggestion  lipgloss.Style
}

type TextInputCursorStyle struct {
	Color color.Color
	Shape tea.CursorShape
	Blink bool
}

// TextAreaStyles is a compatibility type for textarea styling
type TextAreaStyles struct {
	Focused TextAreaStyleState
	Blurred TextAreaStyleState
	Cursor  TextAreaCursorStyle
}

type TextAreaStyleState struct {
	Base             lipgloss.Style
	Text             lipgloss.Style
	LineNumber       lipgloss.Style
	CursorLine       lipgloss.Style
	CursorLineNumber lipgloss.Style
	Placeholder      lipgloss.Style
	Prompt           lipgloss.Style
}

type TextAreaCursorStyle struct {
	Color color.Color
	Shape tea.CursorShape
	Blink bool
}

// SetTextInputStyles applies TextInputStyles to a textinput.Model
func SetTextInputStyles(ti *textinput.Model, styles TextInputStyles) {
	ti.PromptStyle = styles.Focused.Prompt
	ti.TextStyle = styles.Focused.Text
	ti.PlaceholderStyle = styles.Focused.Placeholder
	ti.CompletionStyle = styles.Focused.Suggestion
	// Note: Cursor styling may need manual application via Cursor.Style field
}

// SetTextAreaStyles applies TextAreaStyles to a textarea.Model
func SetTextAreaStyles(ta *textarea.Model, styles TextAreaStyles) {
	ta.FocusedStyle.Base = styles.Focused.Base
	ta.FocusedStyle.Text = styles.Focused.Text
	ta.FocusedStyle.LineNumber = styles.Focused.LineNumber
	ta.FocusedStyle.CursorLine = styles.Focused.CursorLine
	ta.FocusedStyle.CursorLineNumber = styles.Focused.CursorLineNumber
	ta.FocusedStyle.Placeholder = styles.Focused.Placeholder
	ta.FocusedStyle.Prompt = styles.Focused.Prompt

	ta.BlurredStyle.Base = styles.Blurred.Base
	ta.BlurredStyle.Text = styles.Blurred.Text
	ta.BlurredStyle.LineNumber = styles.Blurred.LineNumber
	ta.BlurredStyle.CursorLine = styles.Blurred.CursorLine
	ta.BlurredStyle.CursorLineNumber = styles.Blurred.CursorLineNumber
	ta.BlurredStyle.Placeholder = styles.Blurred.Placeholder
	ta.BlurredStyle.Prompt = styles.Blurred.Prompt
	// Note: Cursor styling may need manual application
}

