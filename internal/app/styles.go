package app

import "github.com/charmbracelet/lipgloss"

var (
	violet     = lipgloss.Color("#a78bfa")
	violetDark = lipgloss.Color("#7c3aed")
	violetDim  = lipgloss.Color("#8b5cf6")
)

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(violet).
			Padding(0, 1)

	ItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#e0e0e0")).
			Padding(0, 1)

	SelectedItemStyle = lipgloss.NewStyle().
				Background(violetDark).
				Foreground(lipgloss.Color("#ffffff")).
				Padding(0, 1)

	ResultStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(violet))

	FooterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6b7280"))

	LoadingStyle = lipgloss.NewStyle().
			Foreground(violet).
			Italic(true)

	LoadingCountStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#f472b6"))

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ef4444")).
			Bold(true)

	DetailTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(violet).
				Padding(0, 1)

	DetailLabelStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(violet).
				Padding(0, 1)

	DetailValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#e0e0e0")).
				Padding(0, 1)

	DetailSectionStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(violetDim).
				Padding(0, 1)

	DividerStyle = lipgloss.NewStyle().
			Foreground(violet)

	SearchBarStyle = lipgloss.NewStyle().
			Foreground(violet).
			Bold(true).
			Padding(0, 1)

	docStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(violet).
			Padding(1, 2)
)
