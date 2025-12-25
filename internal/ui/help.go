package ui

const (
	// HelpTextNavigation contains navigation help text
	HelpTextNavigation = "Navigation: j/k (up/down), g/G (top/bottom), / (search)"
	// HelpTextActions contains action help text
	HelpTextActions = "Actions: Enter (select), Tab (next), Esc (back)"
	// HelpTextGeneral contains general help text
	HelpTextGeneral = "General: ? (help), q/ctrl+c (quit)"
)

// HelpText returns all help text
func HelpText() string {
	return HelpTextNavigation + "\n" + HelpTextActions + "\n" + HelpTextGeneral
}
