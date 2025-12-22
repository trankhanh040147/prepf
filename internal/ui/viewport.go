package ui

import (
	"math"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/trankhanh040147/prepf/internal/config"
)

// ViewportModel wraps viewport with safety checks
type ViewportModel struct {
	viewport viewport.Model
	width    int
	height   int
}

// NewViewport creates a new viewport model
func NewViewport() ViewportModel {
	vp := viewport.New(0, 0)
	return ViewportModel{
		viewport: vp,
		width:    config.DefaultMinWidth,
		height:   20,
	}
}

// SetContent sets viewport content
func (v *ViewportModel) SetContent(content string) {
	v.viewport.SetContent(content)
}

// SetSize sets viewport size with safety checks
func (v *ViewportModel) SetSize(width, height int) {
	// Guard against negative values
	width = int(math.Max(0, float64(width)))
	height = int(math.Max(0, float64(height)))

	// Default to minimum width if zero
	if width == 0 {
		width = config.DefaultMinWidth
	}
	if height == 0 {
		height = 20
	}

	v.width = width
	v.height = height
	v.viewport.Width = width
	v.viewport.Height = height
}

// Update handles viewport messages
func (v *ViewportModel) Update(msg tea.Msg) (*ViewportModel, tea.Cmd) {
	var cmd tea.Cmd
	v.viewport, cmd = v.viewport.Update(msg)
	return v, cmd
}

// View renders the viewport
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
	// Guard against negative values
	width = int(math.Max(0, float64(width)))
	height = int(math.Max(0, float64(height)))

	// Default to minimum width if zero
	if width == 0 {
		width = config.DefaultMinWidth
	}

	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}
