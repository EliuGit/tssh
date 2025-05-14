package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	tableStyle   = lipgloss.NewStyle().Border(lipgloss.NormalBorder())
	footStyle    = lipgloss.NewStyle().Margin(0, 1)
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	noStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	titleStyle   = lipgloss.NewStyle().Bold(true)
	errorStyle   = lipgloss.NewStyle().MaxWidth(80).Inline(true).Foreground(lipgloss.Color("196"))

	dialogBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(1, 2).
			BorderTop(true).
			BorderLeft(true).
			BorderRight(true).
			BorderBottom(true)

	buttonStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("231")).
			Background(lipgloss.Color("57")).
			Padding(0, 3).
			MarginTop(1)

	focusedButtonStyle = buttonStyle.
				Foreground(lipgloss.Color("231")).
				Background(lipgloss.Color("99"))

	questionStyle = lipgloss.NewStyle().Bold(true).MarginBottom(1)

	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)
