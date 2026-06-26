package app

import (
	"fmt"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	cursor   int
	packages []string

	loading bool
	err     error

	width  int
	height int
}

type brewListMsg []string
type brewErrMsg error

func fetchBrewList() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("brew", "list", "--formula")
		out, err := cmd.Output()
		if err != nil {
			return brewErrMsg(err)
		}
		packages := strings.Fields(string(out))
		return brewListMsg(packages)
	}
}

func New() Model {
	return Model{
		loading: true,
	}
}

func (m Model) Init() tea.Cmd {
	return fetchBrewList()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case brewListMsg:
		m.packages = []string(msg)
		m.loading = false
		return m, nil

	case brewErrMsg:
		m.err = error(msg)
		m.loading = false
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c":
			return m, tea.Quit

		case "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down":
			if m.cursor < len(m.packages)-1 {
				m.cursor++
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	switch {
	case m.loading:
		return LoadingStyle.Render("Fetching installed packages...")
	case m.err != nil:
		return ErrorStyle.Render(fmt.Sprintf("Error: %v\n\nEnsure Homebrew is installed.", m.err))
	}

	title := TitleStyle.Render(fmt.Sprintf("🍺 brew-tui  (%d)", len(m.packages)))

	sep := lipgloss.NewStyle().
		Foreground(violet).
		Padding(0, 1).
		Render(strings.Repeat("─", m.width-8))

	var list string

	visibleHeight := m.height - 8

	start := 0

	if m.cursor >= visibleHeight {
		start = m.cursor - visibleHeight + 1
	}

	end := start + visibleHeight

	if end > len(m.packages) {
		end = len(m.packages)
	}

	for i := start; i < end; i++ {
		pkg := m.packages[i]

		if i == m.cursor {
			list += SelectedItemStyle.Render("▸ "+pkg) + "\n"
		} else {
			list += ItemStyle.Render("  "+pkg) + "\n"
		}
	}

	footer := FooterStyle.Render(
		fmt.Sprintf("%d results  •  ↑↓ navigate  •  q quit", len(m.packages)),
	)

	body := lipgloss.JoinVertical(lipgloss.Left, title, sep, list, footer)
	return docStyle.Render(body)
}
