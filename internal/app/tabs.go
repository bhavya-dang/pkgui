package app

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) renderTabBar(width int) string {
	activeStyle := lipgloss.NewStyle().
		Bold(true).
		Background(currentTheme.Primary).
		Foreground(currentTheme.SelectedFg).
		Padding(0, 2)

	inactiveStyle := lipgloss.NewStyle().
		Foreground(currentTheme.DimText).
		Padding(0, 1)

	separator := lipgloss.NewStyle().
		Foreground(currentTheme.Muted).
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
