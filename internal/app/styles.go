package app

import "github.com/charmbracelet/lipgloss"

var (
	teal     = lipgloss.Color("#5bc0be")
	tealDark = lipgloss.Color("#3b9b99")
	amber    = lipgloss.Color("#e4b95b")
)

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(teal).
			Padding(0, 1)

	ItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#c0d4e4")).
			Padding(0, 1)

	SelectedItemStyle = lipgloss.NewStyle().
				Background(tealDark).
				Foreground(lipgloss.Color("#ffffff")).
				Padding(0, 1)

	ResultStyle = lipgloss.NewStyle().Foreground(teal)

	FooterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#60788a"))

	LoadingStyle = lipgloss.NewStyle().
			Foreground(teal).
			Italic(true)

	LoadingCountStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(amber)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ef4444")).
			Bold(true)

	DetailTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(amber).
				Padding(0, 1)

	DetailLabelStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(teal).
				Padding(0, 1)

	DetailValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#c0d4e4")).
				Padding(0, 1)

	DetailSectionStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(amber).
				Padding(0, 1)

	DividerStyle = lipgloss.NewStyle().
			Foreground(teal)

	SearchBarStyle = lipgloss.NewStyle().
			Foreground(teal).
			Bold(true).
			Padding(0, 1)

	LinkStyle = lipgloss.NewStyle().
			Foreground(amber).Underline(true)

	SectionContentStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#e0e0e0")).
				Padding(0, 1)

	docStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(teal).
			Padding(1, 2)
)
