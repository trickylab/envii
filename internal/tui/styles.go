package tui

import "github.com/charmbracelet/lipgloss"

var (
	colorPrimary = lipgloss.Color("#7D56F4")
	colorAccent  = lipgloss.Color("#43BF6D")
	colorMuted   = lipgloss.Color("#626262")
	colorDanger  = lipgloss.Color("#E06C75")
	colorBg      = lipgloss.Color("#1A1A1A")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(colorPrimary).
			Padding(0, 1)

	breadcrumbStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DDDDDD"))

	secretStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Italic(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			Padding(0, 1)

	statusStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Padding(0, 1)

	errorStyle = lipgloss.NewStyle().
			Foreground(colorDanger).
			Padding(0, 1)

	inputBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorPrimary).
			Padding(0, 1)
)
