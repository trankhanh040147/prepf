package ui

import (
	"math"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/trankhanh040147/prepf/internal/config"
	"github.com/trankhanh040147/prepf/internal/util/stringutil"
)

// ViewportModel wraps viewport with safety checks
type ViewportModel struct {
	viewport viewport.Model
	width    int
	height   int
}

// NewViewport creates a new viewport model
func NewViewport() *ViewportModel {
	vp := viewport.New(0, 0)
	return &ViewportModel{
		viewport: vp,
		width:    config.DefaultMinWidth,
		height:   20,
	}
}

// sanitizeDimensions guards against negative values and applies default width
func sanitizeDimensions(width, height int) (int, int) {
	w := int(math.Max(0, float64(width)))
	h := int(math.Max(0, float64(height)))
	if w == 0 {
		w = config.DefaultMinWidth
	}
	return w, h
}

// SetContent sets viewport content with markdown rendering, text wrapping, and padding
func (v *ViewportModel) SetContent(content string) {
	// Use viewport width for rendering, fallback to DefaultMinWidth if width is 0
	wrapWidth := v.width
	if wrapWidth == 0 {
		wrapWidth = config.DefaultMinWidth
	}
	// Account for horizontal padding when rendering (reduce width by padding on both sides)
	effectiveWidth := wrapWidth - (2 * config.ViewportContentPadding)
	if effectiveWidth < 1 {
		effectiveWidth = 1
	}

	// Render markdown content (gracefully degrades to plain text if rendering fails)
	renderedContent := stringutil.RenderMarkdown(content, effectiveWidth)

	// Apply padding to rendered content using lipgloss
	paddingStyle := lipgloss.NewStyle().Padding(config.ViewportTopPadding, config.ViewportContentPadding)
	paddedContent := paddingStyle.Render(renderedContent)

	v.viewport.SetContent(paddedContent)
}

// GotoBottom scrolls the viewport to the bottom
func (v *ViewportModel) GotoBottom() {
	v.viewport.GotoBottom()
}

// GotoTop scrolls the viewport to the top
func (v *ViewportModel) GotoTop() {
	v.viewport.SetYOffset(0)
}

// SetSize sets viewport size with safety checks
func (v *ViewportModel) SetSize(width, height int) {
	w, h := sanitizeDimensions(width, height)
	if h == 0 {
		h = 20
	}

	v.width = w
	v.height = h
	v.viewport.Width = w
	v.viewport.Height = h
}

// Update handles viewport messages
func (v *ViewportModel) Update(msg tea.Msg) (*ViewportModel, tea.Cmd) {
	var cmd tea.Cmd
	v.viewport, cmd = v.viewport.Update(msg)
	return v, cmd
}

// View renders the viewport (padding is applied to content in SetContent)
func (v *ViewportModel) View() string {
	return v.viewport.View()
}

// Width returns viewport width
func (v *ViewportModel) Width() int {
	return v.width
}

// Height returns viewport height
func (v *ViewportModel) Height() int {
	return v.height
}

// Center centers content using lipgloss.Place
func Center(content string, width, height int) string {
	w, h := sanitizeDimensions(width, height)
	return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, content)
}
