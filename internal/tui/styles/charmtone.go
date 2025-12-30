package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/exp/charmtone"
)

func NewCharmtoneTheme() *Theme {
	t := &Theme{
		Name:   "charmtone",
		IsDark: true,

		Primary:   charmtone.Charple,
		Secondary: charmtone.Dolly,
		Tertiary:  charmtone.Bok,
		Accent:    charmtone.Zest,

		// Backgrounds
		BgBase:        charmtone.Pepper,
		BgBaseLighter: charmtone.BBQ,
		BgSubtle:      charmtone.Charcoal,
		BgOverlay:     charmtone.Iron,

		// Foregrounds
		FgBase:      charmtone.Ash,
		FgMuted:     charmtone.Squid,
		FgHalfMuted: charmtone.Smoke,
		FgSubtle:    charmtone.Oyster,
		FgSelected:  charmtone.Salt,

		// Borders
		Border:      charmtone.Charcoal,
		BorderFocus: charmtone.Charple,

		// Status
		Success: charmtone.Guac,
		Error:   charmtone.Sriracha,
		Warning: charmtone.Zest,
		Info:    charmtone.Malibu,

		// Colors
		White: charmtone.Butter,

		BlueLight: charmtone.Sardine,
		BlueDark:  charmtone.Damson,
		Blue:      charmtone.Malibu,

		Yellow: charmtone.Mustard,
		Citron: charmtone.Citron,

		Green:      charmtone.Julep,
		GreenDark:  charmtone.Guac,
		GreenLight: charmtone.Bok,

		Red:      charmtone.Coral,
		RedDark:  charmtone.Sriracha,
		RedLight: charmtone.Salmon,
		Cherry:   charmtone.Cherry,
	}

	// Text selection.
	t.TextSelection = lipgloss.NewStyle().Foreground(lipgloss.Color(charmtone.Salt.Hex())).Background(lipgloss.Color(charmtone.Charple.Hex()))

	// LSP and MCP status.
	t.ItemOfflineIcon = lipgloss.NewStyle().Foreground(lipgloss.Color(charmtone.Squid.Hex())).SetString("●")
	t.ItemBusyIcon = t.ItemOfflineIcon.Foreground(lipgloss.Color(charmtone.Citron.Hex()))
	t.ItemErrorIcon = t.ItemOfflineIcon.Foreground(lipgloss.Color(charmtone.Coral.Hex()))
	t.ItemOnlineIcon = t.ItemOfflineIcon.Foreground(lipgloss.Color(charmtone.Guac.Hex()))

	// Editor: Yolo Mode.
	t.YoloIconFocused = lipgloss.NewStyle().Foreground(lipgloss.Color(charmtone.Oyster.Hex())).Background(lipgloss.Color(charmtone.Citron.Hex())).Bold(true).SetString(" ! ")
	t.YoloIconBlurred = t.YoloIconFocused.Foreground(lipgloss.Color(charmtone.Pepper.Hex())).Background(lipgloss.Color(charmtone.Squid.Hex()))
	t.YoloDotsFocused = lipgloss.NewStyle().Foreground(lipgloss.Color(charmtone.Zest.Hex())).SetString(":::")
	t.YoloDotsBlurred = t.YoloDotsFocused.Foreground(lipgloss.Color(charmtone.Squid.Hex()))

	// oAuth Chooser.
	t.AuthBorderSelected = lipgloss.NewStyle().BorderForeground(lipgloss.Color(charmtone.Guac.Hex()))
	t.AuthTextSelected = lipgloss.NewStyle().Foreground(lipgloss.Color(charmtone.Julep.Hex()))
	t.AuthBorderUnselected = lipgloss.NewStyle().BorderForeground(lipgloss.Color(charmtone.Iron.Hex()))
	t.AuthTextUnselected = lipgloss.NewStyle().Foreground(lipgloss.Color(charmtone.Squid.Hex()))

	return t
}
