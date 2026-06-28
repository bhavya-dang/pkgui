package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	cursor          int
	packages        []string
	displayPackages []string

	loading bool
	err     error

	width  int
	height int

	formulaeMap   map[string]FormulaData
	formulaeReady bool
	apiErr        error

	downloadCh    chan tea.Msg
	formulaeCount int

	searchActive bool
	searchQuery  string

	info *FormulaData
}

type FormulaData struct {
	Name     string `json:"name"`
	Desc     string `json:"desc"`
	Homepage string `json:"homepage"`
	License  string `json:"license"`
	Versions struct {
		Stable string `json:"stable"`
	} `json:"versions"`
	Dependencies      []string `json:"dependencies"`
	BuildDependencies []string `json:"build_dependencies"`
}

type brewListMsg []string
type brewErrMsg error
type brewFormulaeMsg map[string]FormulaData
type brewFormulaeErrMsg error
type formulaeProgressMsg struct {
	Count int
}

func fetchBrewList() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("brew", "list", "--formula")
		out, err := cmd.Output()
		if err != nil {
			return brewErrMsg(err)
		}
		return brewListMsg(strings.Fields(string(out)))
	}
}

func startDownload(ch chan<- tea.Msg) {
	go func() {
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Get("https://formulae.brew.sh/api/formula.json")
		if err != nil {
			ch <- brewFormulaeErrMsg(err)
			return
		}
		defer resp.Body.Close()

		dec := json.NewDecoder(resp.Body)

		_, err = dec.Token()
		if err != nil {
			ch <- brewFormulaeErrMsg(err)
			return
		}

		m := make(map[string]FormulaData)
		count := 0

		for dec.More() {
			var f FormulaData
			if err := dec.Decode(&f); err != nil {
				ch <- brewFormulaeErrMsg(err)
				return
			}
			m[f.Name] = f
			count++
			if count%50 == 0 {
				ch <- formulaeProgressMsg{Count: count}
			}
		}

		ch <- formulaeProgressMsg{Count: count}
		ch <- brewFormulaeMsg(m)
	}()
}

func (m Model) recvDownload() tea.Cmd {
	if m.downloadCh == nil {
		return nil
	}
	return func() tea.Msg {
		msg, ok := <-m.downloadCh
		if !ok {
			return nil
		}
		return msg
	}
}

func New() Model {
	return Model{
		loading:    true,
		downloadCh: make(chan tea.Msg, 200),
	}
}

func (m Model) Init() tea.Cmd {
	startDownload(m.downloadCh)
	return tea.Batch(fetchBrewList(), m.recvDownload())
}

func fuzzyMatch(s, query string) bool {
	s = strings.ToLower(s)
	q := strings.ToLower(query)
	qi := 0
	for i := 0; i < len(s) && qi < len(q); i++ {
		if s[i] == q[qi] {
			qi++
		}
	}
	return qi == len(q)
}

func (m Model) applyFilter() Model {
	if m.searchQuery == "" {
		m.displayPackages = m.packages
	} else {
		query := strings.ToLower(m.searchQuery)
		var filtered []string
		for _, pkg := range m.packages {
			if fuzzyMatch(pkg, query) {
				filtered = append(filtered, pkg)
			}
		}
		m.displayPackages = filtered
	}
	if m.cursor >= len(m.displayPackages) {
		m.cursor = max(0, len(m.displayPackages)-1)
	}
	if len(m.displayPackages) > 0 {
		return m.updateInfo()
	}
	m.info = nil
	return m
}

func (m Model) updateInfo() Model {
	if len(m.displayPackages) > 0 && m.cursor < len(m.displayPackages) {
		name := m.displayPackages[m.cursor]
		if f, ok := m.formulaeMap[name]; ok {
			m.info = &f
		} else {
			m.info = nil
		}
	}
	return m
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case formulaeProgressMsg:
		m.formulaeCount = msg.Count
		return m, m.recvDownload()

	case brewListMsg:
		m.packages = []string(msg)
		m.displayPackages = m.packages
		if m.formulaeReady {
			m.loading = false
			m = m.updateInfo()
		}

	case brewErrMsg:
		m.err = error(msg)
		m.loading = false

	case brewFormulaeMsg:
		m.formulaeMap = map[string]FormulaData(msg)
		m.formulaeReady = true
		m.downloadCh = nil
		if m.packages != nil {
			m.loading = false
			m = m.updateInfo()
		}

	case brewFormulaeErrMsg:
		m.apiErr = error(msg)
		m.formulaeReady = true
		m.downloadCh = nil
		if m.packages != nil {
			m.loading = false
		}

	case tea.KeyMsg:
		if m.searchActive {
			switch msg.String() {
			case "esc":
				m.searchActive = false
				m.searchQuery = ""
				m = m.applyFilter()
			case "enter":
				m.searchActive = false
			case "backspace":
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
					m = m.applyFilter()
				}
			default:
				if len(msg.String()) == 1 && msg.String()[0] >= 32 {
					m.searchQuery += msg.String()
					m = m.applyFilter()
				}
			}
			return m, nil
		}

		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "/":
			m.searchActive = true
			m.searchQuery = ""

		case "up":
			if m.cursor > 0 {
				m.cursor--
				m = m.updateInfo()
			}

		case "down":
			if m.cursor < len(m.displayPackages)-1 {
				m.cursor++
				m = m.updateInfo()
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	switch {
	case m.loading:
		return m.renderLoading()
	case m.err != nil:
		return ErrorStyle.Render(fmt.Sprintf("Error: %v\n\nEnsure Homebrew is installed.", m.err))
	}

	if m.width < 60 {
		return m.listViewFallback()
	}

	contentWidth := m.width - 6
	dividerWidth := 1
	leftWidth := int(float64(contentWidth-dividerWidth) * 0.35)
	rightWidth := contentWidth - dividerWidth - leftWidth

	searchOffset := 0
	var searchLine string
	if m.searchActive {
		searchLine = m.renderSearchBar(contentWidth)
		searchOffset = 1
	}

	panelHeight := m.height - 5 - searchOffset

	leftPanel := m.renderLeftPanel(leftWidth)
	rightPanel := m.renderRightPanel(rightWidth)
	divider := m.renderDivider(panelHeight)

	leftStyled := lipgloss.NewStyle().Width(leftWidth).Height(panelHeight).Render(leftPanel)
	rightStyled := lipgloss.NewStyle().Width(rightWidth).Height(panelHeight).Render(rightPanel)

	top := lipgloss.JoinHorizontal(lipgloss.Top, leftStyled, divider, rightStyled)

	var bodyParts []string
	if searchLine != "" {
		bodyParts = append(bodyParts, searchLine)
	}
	bodyParts = append(bodyParts, top)
	bodyParts = append(bodyParts, m.renderFooter())

	body := lipgloss.JoinVertical(lipgloss.Left, bodyParts...)
	return docStyle.Render(body)
}

func (m Model) renderLoading() string {
	countLine := "  Loading formula data..."
	if m.formulaeCount > 0 {
		countLine = fmt.Sprintf("  %s formulae loaded",
			LoadingCountStyle.Render(fmt.Sprintf("%d", m.formulaeCount)))
	}
	return LoadingStyle.Render(strings.Join([]string{
		"",
		"  Loading pkgui...",
		"",
		countLine,
	}, "\n"))
}

func (m Model) renderSearchBar(width int) string {
	text := m.searchQuery
	if text == "" {
		text = "█"
	}
	return SearchBarStyle.Render("  / " + text)
}

func (m Model) renderLeftPanel(width int) string {
	titleText := fmt.Sprintf("pkgui  (%d)", len(m.packages))
	if m.searchQuery != "" {
		titleText = fmt.Sprintf("pkgui  (%d/%d)", len(m.displayPackages), len(m.packages))
	}
	title := TitleStyle.Render(titleText)

	sep := lipgloss.NewStyle().
		Foreground(violet).
		Padding(0, 1).
		Render(strings.Repeat("─", max(0, width-4)))

	panelHeight := m.height - 5
	visibleHeight := panelHeight - 3

	start := 0
	if m.cursor >= visibleHeight {
		start = m.cursor - visibleHeight + 1
	}
	end := start + visibleHeight
	if end > len(m.displayPackages) {
		end = len(m.displayPackages)
	}

	var listItems []string
	for i := start; i < end; i++ {
		pkg := m.displayPackages[i]
		if i == m.cursor {
			listItems = append(listItems, SelectedItemStyle.Render("▸ "+pkg))
		} else {
			listItems = append(listItems, ItemStyle.Render("  "+pkg))
		}
	}

	// result := ResultStyle.Render(fmt.Sprintf("%d formulae", len(m.displayPackages)))

	var lines []string
	lines = append(lines, title)
	lines = append(lines, sep)
	lines = append(lines, listItems...)
	// lines = append(lines, result)
	return strings.Join(lines, "\n")
}

func (m Model) renderRightPanel(width int) string {
	if len(m.displayPackages) == 0 {
		return ""
	}

	pkgName := m.displayPackages[m.cursor]
	title := lipgloss.NewStyle().
		Width(width).
		Render(DetailTitleStyle.Render("▸ " + pkgName))

	sep := lipgloss.NewStyle().
		Foreground(violet).
		Padding(0, 1).
		Render(strings.Repeat("─", max(0, width-4)))

	var contentLines []string
	contentLines = append(contentLines, title)
	contentLines = append(contentLines, sep)

	if m.apiErr != nil {
		contentLines = append(contentLines,
			DetailValueStyle.Render("  Formula data unavailable"))
	} else if m.info != nil {
		info := m.info

		contentLines = append(contentLines,
			DetailLabelStyle.Render("Version:")+"  "+DetailValueStyle.Render(info.Versions.Stable))
		contentLines = append(contentLines, "")

		if info.Desc != "" {
			contentLines = append(contentLines,
				DetailLabelStyle.Render("Description:"))
			contentLines = append(contentLines,
				DetailValueStyle.Render("  "+info.Desc))
			contentLines = append(contentLines, "")
		}

		if info.Homepage != "" {
			contentLines = append(contentLines,
				DetailLabelStyle.Render("Homepage:")+"  "+DetailValueStyle.Render(info.Homepage))
			contentLines = append(contentLines, "")
		}

		if info.License != "" {
			contentLines = append(contentLines,
				DetailLabelStyle.Render("License:")+"  "+DetailValueStyle.Render(info.License))
			contentLines = append(contentLines, "")
		}

		if len(info.Dependencies) > 0 {
			contentLines = append(contentLines,
				DetailSectionStyle.Render("Dependencies"))
			contentLines = append(contentLines,
				DetailValueStyle.Render("  "+strings.Join(info.Dependencies, ", ")))
			contentLines = append(contentLines, "")
		}

		if len(info.BuildDependencies) > 0 {
			contentLines = append(contentLines,
				DetailSectionStyle.Render("Build Dependencies"))
			contentLines = append(contentLines,
				DetailValueStyle.Render("  "+strings.Join(info.BuildDependencies, ", ")))
		}
	} else {
		contentLines = append(contentLines,
			DetailValueStyle.Render("  No formula data available"))
	}

	return strings.Join(contentLines, "\n")
}

func (m Model) renderDivider(height int) string {
	lines := make([]string, height)
	for i := range lines {
		lines[i] = DividerStyle.Render("│")
	}
	return strings.Join(lines, "\n")
}

func (m Model) renderFooter() string {
	left := ResultStyle.Render(fmt.Sprintf("%d formulae", len(m.packages)))
	if m.apiErr != nil {
		left += "  " + ErrorStyle.Render("API unavailable")
	}
	help := FooterStyle.Render("•  / search  •  ↑↓ navigate  •  q quit")
	return left + "  " + help
}

func (m Model) listViewFallback() string {
	title := TitleStyle.Render(fmt.Sprintf("pkgui  (%d)", len(m.packages)))

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
	if end > len(m.displayPackages) {
		end = len(m.displayPackages)
	}

	for i := start; i < end; i++ {
		pkg := m.displayPackages[i]
		if i == m.cursor {
			list += SelectedItemStyle.Render("▸ "+pkg) + "\n"
		} else {
			list += ItemStyle.Render("  "+pkg) + "\n"
		}
	}

	result := ResultStyle.Render(fmt.Sprintf("%d formulae", len(m.displayPackages)))
	footer := result + FooterStyle.Render("•  / search  •  ↑↓ navigate  •  q quit")
	body := lipgloss.JoinVertical(lipgloss.Left, title, sep, list, footer)
	return docStyle.Render(body)
}
