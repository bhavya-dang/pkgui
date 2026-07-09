package app

import "github.com/charmbracelet/lipgloss"

var (
	currentTheme *Theme

	violet     lipgloss.Color
	violetDark lipgloss.Color

	TitleStyle          lipgloss.Style
	ItemStyle           lipgloss.Style
	SelectedItemStyle   lipgloss.Style
	ResultStyle         lipgloss.Style
	FooterStyle         lipgloss.Style
	LoadingStyle        lipgloss.Style
	LoadingCountStyle   lipgloss.Style
	ErrorStyle          lipgloss.Style
	DetailTitleStyle    lipgloss.Style
	DetailLabelStyle    lipgloss.Style
	DetailValueStyle    lipgloss.Style
	DetailSectionStyle  lipgloss.Style
	SearchBarStyle      lipgloss.Style
	SearchPlaceholderStyle lipgloss.Style
	LinkStyle           lipgloss.Style
	SectionContentStyle lipgloss.Style
	docStyle            lipgloss.Style
)

func applyTheme(t *Theme) {
	currentTheme = t
	violet = t.Primary
	violetDark = t.Muted

	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(violet).
		Padding(0, 1)

	ItemStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Text)

	SelectedItemStyle = lipgloss.NewStyle().
		Bold(true).
		Background(t.SelectedBg).
		Foreground(t.SelectedFg)

	ResultStyle = lipgloss.NewStyle().Foreground(violet)

	FooterStyle = lipgloss.NewStyle().
		Foreground(t.DimText)

	LoadingStyle = lipgloss.NewStyle().
		Foreground(violet).
		Italic(true)

	LoadingCountStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(violet)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(t.Error).
		Bold(true)

	DetailTitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(violet)

	DetailLabelStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(violet)

	DetailValueStyle = lipgloss.NewStyle().
		Foreground(t.DetailText)

	DetailSectionStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(violet).
		Padding(0, 1)

	SearchBarStyle = lipgloss.NewStyle().
		Foreground(violet).
		Bold(true).
		Padding(0, 1)

	SearchPlaceholderStyle = lipgloss.NewStyle().
		Foreground(t.DimText).
		Italic(true)

	LinkStyle = lipgloss.NewStyle().
		Foreground(violet).
		Underline(true)

	SectionContentStyle = lipgloss.NewStyle().
		Foreground(t.Text).
		Padding(0, 1)

	docStyle = lipgloss.NewStyle().
		Padding(1, 3, 1, 3)
}
