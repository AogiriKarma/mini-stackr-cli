package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	accentColor     = lipgloss.Color("12")
	successColor    = lipgloss.Color("10")
	errorColor      = lipgloss.Color("9")
	warningColor    = lipgloss.Color("11")
	mutedColor      = lipgloss.Color("8")
	textColor       = lipgloss.Color("15")
	bgSelectedColor = lipgloss.Color("236")

	// Title
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor).
			Padding(0, 1)

	// List styles
	selectedStyle = lipgloss.NewStyle().
			Foreground(accentColor).
			Bold(true)

	runningStyle = lipgloss.NewStyle().
			Foreground(successColor)

	stoppedStyle = lipgloss.NewStyle().
			Foreground(errorColor)

	// Status indicators
	runningDot = lipgloss.NewStyle().Foreground(successColor).Render("●")
	stoppedDot = lipgloss.NewStyle().Foreground(errorColor).Render("○")
	pausedDot  = lipgloss.NewStyle().Foreground(warningColor).Render("◐")

	// Text styles
	nameStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(textColor)

	imageStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	labelStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	valueStyle = lipgloss.NewStyle().
			Foreground(textColor)

	// Boxes for detail view
	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(mutedColor).
			Padding(0, 1)

	boxTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(accentColor)

	// Header bar
	headerBarStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("0")).
			Background(accentColor).
			Padding(0, 2).
			MarginBottom(1)

	// Help
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(1)

	// Progress bar
	progressFull  = lipgloss.NewStyle().Foreground(accentColor).Render("█")
	progressEmpty = lipgloss.NewStyle().Foreground(mutedColor).Render("░")
)

func renderProgressBar(percent float64, width int) string {
	filled := int(percent / 100 * float64(width))
	empty := width - filled
	if filled < 0 {
		filled = 0
	}
	if empty < 0 {
		empty = 0
	}

	bar := ""
	for i := 0; i < filled; i++ {
		bar += progressFull
	}
	for i := 0; i < empty; i++ {
		bar += progressEmpty
	}
	return bar
}

func statusDot(state string) string {
	switch state {
	case "running":
		return runningDot
	case "paused":
		return pausedDot
	default:
		return stoppedDot
	}
}