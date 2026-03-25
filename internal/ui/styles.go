package ui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF6600")).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF6600")).
			Padding(0, 2)

	ISSMarkerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF00FF")).
			Background(lipgloss.Color("#FFFFFF"))

	MapStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#444444")).
			Padding(0, 1)

	InfoStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#00AAFF")).
			Padding(0, 1).
			MarginTop(1)

	LabelStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00AAFF"))

	ValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	AstronautStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF88"))

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			MarginTop(1)

	UserMarkerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#0088FF"))

	LandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#55DD55"))

	OceanStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#112233"))

	WarnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFAA00")).
			MarginTop(1)
)
