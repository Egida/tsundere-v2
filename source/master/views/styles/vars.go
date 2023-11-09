package styles

import "github.com/charmbracelet/lipgloss"

var (
	Keypad = []int{
		1, 2, 3,
		4, 5, 6,
		7, 8, 9,
	}

	BaseButton = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#9d74db")).Align(lipgloss.Center)

	BigButton   = BaseButton.Copy().Width(32)
	SmallButton = BaseButton.Copy().Width(14)

	Text = lipgloss.NewStyle()

	ErrorText   = Text.Copy().Foreground(lipgloss.Color("#e64c4c"))
	SuccessText = Text.Copy().Foreground(lipgloss.Color("#67e66b"))
	MainText    = Text.Copy().Foreground(lipgloss.Color("#9d74db"))
)
