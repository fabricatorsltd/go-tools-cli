package wizard

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FF87")).
		MarginBottom(1)

	subtitleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true)

	activeStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00CFFF"))

	inactiveStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#555555"))

	selectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FF87"))

	unselectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#555555"))

	progressStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#444444"))

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF5F87"))

	successStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00FF87"))

	boxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#333333")).
		Padding(0, 1)

	labelStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#AAAAAA"))

	valueStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))
)
