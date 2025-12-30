package diffview

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/charmtone"
)

// LineStyle defines the styles for a given line type in the diff view.
type LineStyle struct {
	LineNumber lipgloss.Style
	Symbol     lipgloss.Style
	Code       lipgloss.Style
}

// Style defines the overall style for the diff view, including styles for
// different line types such as divider, missing, equal, insert, and delete
// lines.
type Style struct {
	DividerLine LineStyle
	MissingLine LineStyle
	EqualLine   LineStyle
	InsertLine  LineStyle
	DeleteLine  LineStyle
}

// DefaultLightStyle provides a default light theme style for the diff view.
func DefaultLightStyle() Style {
	return Style{
		DividerLine: LineStyle{
			LineNumber: lipgloss.NewStyle().
				Foreground(lipgloss.Color(charmtone.Iron.Hex())).
				Background(lipgloss.Color(charmtone.Thunder.Hex())),
			Code: lipgloss.NewStyle().
				Foreground(lipgloss.Color(charmtone.Oyster.Hex())).
				Background(lipgloss.Color(charmtone.Anchovy.Hex())),
		},
		MissingLine: LineStyle{
			LineNumber: lipgloss.NewStyle().
				Background(lipgloss.Color(charmtone.Ash.Hex())),
			Code: lipgloss.NewStyle().
				Background(lipgloss.Color(charmtone.Ash.Hex())),
		},
		EqualLine: LineStyle{
			LineNumber: lipgloss.NewStyle().
				Foreground(lipgloss.Color(charmtone.Charcoal.Hex())).
				Background(lipgloss.Color(charmtone.Ash.Hex())),
			Code: lipgloss.NewStyle().
				Foreground(lipgloss.Color(charmtone.Pepper.Hex())).
				Background(lipgloss.Color(charmtone.Salt.Hex())),
		},
		InsertLine: LineStyle{
			LineNumber: lipgloss.NewStyle().
				Foreground(lipgloss.Color(charmtone.Turtle.Hex())).
				Background(lipgloss.Color("#c8e6c9")),
			Symbol: lipgloss.NewStyle().
				Foreground(lipgloss.Color(charmtone.Turtle.Hex())).
				Background(lipgloss.Color("#e8f5e9")),
			Code: lipgloss.NewStyle().
				Foreground(lipgloss.Color(charmtone.Pepper.Hex())).
				Background(lipgloss.Color("#e8f5e9")),
		},
		DeleteLine: LineStyle{
			LineNumber: lipgloss.NewStyle().
				Foreground(lipgloss.Color(charmtone.Cherry.Hex())).
				Background(lipgloss.Color("#ffcdd2")),
			Symbol: lipgloss.NewStyle().
				Foreground(lipgloss.Color(charmtone.Cherry.Hex())).
				Background(lipgloss.Color("#ffebee")),
			Code: lipgloss.NewStyle().
				Foreground(lipgloss.Color(charmtone.Pepper.Hex())).
				Background(lipgloss.Color("#ffebee")),
		},
	}
}

// DefaultDarkStyle provides a default dark theme style for the diff view.
func DefaultDarkStyle() Style {
	return Style{
		DividerLine: LineStyle{
			LineNumber: lipgloss.NewStyle().
				Foreground(lipgloss.Color(charmtone.Smoke.Hex())).
				Background(lipgloss.Color(charmtone.Sapphire.Hex())),
			Code: lipgloss.NewStyle().
				Foreground(lipgloss.Color(charmtone.Smoke.Hex())).
				Background(lipgloss.Color(charmtone.Ox.Hex())),
		},
		MissingLine: LineStyle{
			LineNumber: lipgloss.NewStyle().
				Background(lipgloss.Color(charmtone.Charcoal.Hex())),
			Code: lipgloss.NewStyle().
				Background(lipgloss.Color(charmtone.Charcoal.Hex())),
		},
		EqualLine: LineStyle{
			LineNumber: lipgloss.NewStyle().
				Foreground(lipgloss.Color(charmtone.Ash.Hex())).
				Background(lipgloss.Color(charmtone.Charcoal.Hex())),
			Code: lipgloss.NewStyle().
				Foreground(lipgloss.Color(charmtone.Salt.Hex())).
				Background(lipgloss.Color(charmtone.Pepper.Hex())),
		},
		InsertLine: LineStyle{
			LineNumber: lipgloss.NewStyle().
				Foreground(lipgloss.Color(charmtone.Turtle.Hex())).
				Background(lipgloss.Color("#293229")),
			Symbol: lipgloss.NewStyle().
				Foreground(lipgloss.Color(charmtone.Turtle.Hex())).
				Background(lipgloss.Color("#303a30")),
			Code: lipgloss.NewStyle().
				Foreground(lipgloss.Color(charmtone.Salt.Hex())).
				Background(lipgloss.Color("#303a30")),
		},
		DeleteLine: LineStyle{
			LineNumber: lipgloss.NewStyle().
				Foreground(lipgloss.Color(charmtone.Cherry.Hex())).
				Background(lipgloss.Color("#332929")),
			Symbol: lipgloss.NewStyle().
				Foreground(lipgloss.Color(charmtone.Cherry.Hex())).
				Background(lipgloss.Color("#3a3030")),
			Code: lipgloss.NewStyle().
				Foreground(lipgloss.Color(charmtone.Salt.Hex())).
				Background(lipgloss.Color("#3a3030")),
		},
	}
}
