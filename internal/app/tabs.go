package app

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderTabBar(width int) string {
	activeStyle := lipgloss.NewStyle().
		Bold(true).
		Background(amber).
		Foreground(lipgloss.Color("#000")).
		Padding(0, 2)

	inactiveStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#60788a")).
		Padding(0, 1)

	separator := lipgloss.NewStyle().
		Foreground(tealDark).
		Render(" ")

	var cells []string
	for i, tab := range m.tabs {
		label := strings.ToUpper(tab.TabLabel())
		if i == m.activeTab {
			cells = append(cells, activeStyle.Render(label))
		} else {
			cells = append(cells, inactiveStyle.Render(label))

		}
	}

	tabLine := strings.Join(cells, separator)

	return lipgloss.NewStyle().
		Width(width).
		PaddingTop(1).
		Render(tabLine)
}
