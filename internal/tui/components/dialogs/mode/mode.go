package mode

import (
	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/trankhanh040147/prepf/internal/tui/components/core"
	"github.com/trankhanh040147/prepf/internal/tui/components/dialogs"
	"github.com/trankhanh040147/prepf/internal/tui/styles"
	"github.com/trankhanh040147/prepf/internal/tui/util"
)

const ModeDialogID dialogs.DialogID = "mode"

type ModeDialog interface {
	dialogs.DialogModel
}

type ModeOption struct {
	ID   string
	Name string
	Desc string
}

var (
	ModeMock = ModeOption{
		ID:   "mock",
		Name: "Mock (The Gauntlet)",
		Desc: "Real-world interview simulation with harsh feedback",
	}
	ModeGym = ModeOption{
		ID:   "gym",
		Name: "Gym (Training)",
		Desc: "Targeted practice questions with immediate feedback",
	}
)

type modeDialogCmp struct {
	selectedIndex int
	wWidth        int
	wHeight       int
	width         int
	options       []ModeOption
	keyMap        KeyMap
	help          help.Model
}

func NewModeDialogCmp() ModeDialog {
	t := styles.CurrentTheme()
	keyMap := DefaultKeyMap()
	help := help.New()
	help.Styles = t.S().Help

	options := []ModeOption{ModeMock, ModeGym}

	s := &modeDialogCmp{
		selectedIndex: 0,
		options:       options,
		keyMap:        keyMap,
		help:          help,
	}

	return s
}

func (s *modeDialogCmp) Init() tea.Cmd {
	return nil
}

func (s *modeDialogCmp) Update(msg tea.Msg) (util.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.wWidth = msg.Width
		s.wHeight = msg.Height
		s.width = min(80, s.wWidth-8)
		return s, nil
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, s.keyMap.Select):
			if s.selectedIndex < len(s.options) {
				selected := s.options[s.selectedIndex]
				return s, tea.Sequence(
					util.CmdHandler(dialogs.CloseDialogMsg{}),
					util.CmdHandler(ModeSelectedMsg{Mode: selected.ID}),
				)
			}
		case key.Matches(msg, s.keyMap.Close):
			return s, util.CmdHandler(dialogs.CloseDialogMsg{})
		case key.Matches(msg, s.keyMap.Next):
			s.selectedIndex = (s.selectedIndex + 1) % len(s.options)
			return s, nil
		case key.Matches(msg, s.keyMap.Previous):
			s.selectedIndex = (s.selectedIndex - 1 + len(s.options)) % len(s.options)
			return s, nil
		}
	}
	return s, nil
}

func (s *modeDialogCmp) View() string {
	t := styles.CurrentTheme()
	var items []string

	for i, opt := range s.options {
		prefix := "  "
		style := t.S().Base
		if i == s.selectedIndex {
			prefix = "▶ "
			style = style.Foreground(t.Primary).Bold(true)
		}
		item := prefix + opt.Name
		if opt.Desc != "" {
			item += "\n    " + t.S().Base.Foreground(t.Secondary).Render(opt.Desc)
		}
		items = append(items, style.Render(item))
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		t.S().Base.Padding(0, 1, 1, 1).Render(core.Title("Select Mode", s.width-4)),
		lipgloss.JoinVertical(lipgloss.Left, items...),
		"",
		t.S().Base.Width(s.width-2).PaddingLeft(1).AlignHorizontal(lipgloss.Left).Render(s.help.View(s.keyMap)),
	)

	return s.style().Render(content)
}

func (s *modeDialogCmp) Cursor() *tea.Cursor {
	return nil
}

func (s *modeDialogCmp) style() lipgloss.Style {
	t := styles.CurrentTheme()
	return t.S().Base.
		Width(s.width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.BorderFocus)
}

func (s *modeDialogCmp) Position() (int, int) {
	row := s.wHeight/4 - 2
	col := s.wWidth / 2
	col -= s.width / 2
	return row, col
}

func (s *modeDialogCmp) ID() dialogs.DialogID {
	return ModeDialogID
}

type ModeSelectedMsg struct {
	Mode string
}
