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
	packages        []string // all installed formulae
	displayPackages []string // filtered subset for display

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

	installPaths      map[string]string
	installedVersions map[string]string
}

// Formula schema from the Homebrew JSON API.
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

type brewListMsg struct {
	names             []string
	paths             map[string]string
	installedVersions map[string]string
}

type brewErrMsg error
type brewFormulaeMsg map[string]FormulaData
type brewFormulaeErrMsg error
type formulaeProgressMsg struct {
	Count int
}

// Runs the command and parses output to get formulae list, installed versions, and paths
func fetchBrewList() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("brew", "list", "--formula", "--versions")
		out, err := cmd.Output()
		if err != nil {
			return brewErrMsg(err)
		}
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")

		var names []string
		paths := make(map[string]string)
		installedVersions := make(map[string]string)

		for _, line := range lines {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				name, ver := parts[0], parts[1]
				names = append(names, name)
				installedVersions[name] = ver
			}
		}

		if len(names) > 0 {
			prefixOut, perr := exec.Command("brew", "--prefix").Output()
			if perr == nil {
				prefix := strings.TrimSpace(string(prefixOut))
				for _, name := range names {
					paths[name] = prefix + "/opt/" + name
				}
			}
		}

		return brewListMsg{names, paths, installedVersions}
	}
}

// Draws a bordered box with a title for detail sections
func renderSection(maxWidth int, title string, lines ...string) string {
	amberStyle := lipgloss.NewStyle().Bold(true).Foreground(amber)
	border := lipgloss.NewStyle().Foreground(teal)

	maxContent := 0
	for _, line := range lines {
		w := lipgloss.Width(line)
		if w > maxContent {
			maxContent = w
		}
	}
	boxWidth := max(maxContent+4, lipgloss.Width(title)+6)
	boxWidth = min(boxWidth, maxWidth)

	inner := boxWidth - 4

	top := border.Render("╭─ ") +
		amberStyle.Render(title) +
		border.Render(" "+strings.Repeat("─", max(0, boxWidth-5-lipgloss.Width(title)))+"╮")

	var body []string
	for _, line := range lines {
		padded := lipgloss.NewStyle().Width(inner).Render(line)
		body = append(body, border.Render("│ ")+padded+border.Render(" │"))
	}

	bottom := border.Render("╰" + strings.Repeat("─", boxWidth-2) + "╯")

	return strings.Join(append([]string{top}, append(body, bottom)...), "\n")
}

// Wraps content in a bordered box for the left/right panels.
func renderPaneBox(width int, title string, content string) string {
	amberStyle := lipgloss.NewStyle().Bold(true).Foreground(amber)
	border := lipgloss.NewStyle().Foreground(teal)

	top := border.Render("╭─ ") +
		amberStyle.Render(title) +
		border.Render(" "+strings.Repeat("─", max(0, width-5-lipgloss.Width(title)))+"╮")

	inner := width - 4
	lines := strings.Split(content, "\n")
	var body []string
	for _, line := range lines {
		padded := lipgloss.NewStyle().Width(inner).Render(line)
		body = append(body, border.Render("│ ")+padded+border.Render(" │"))
	}

	bottom := border.Render("╰" + strings.Repeat("─", width-2) + "╯")

	return strings.Join(append([]string{top}, append(body, bottom)...), "\n")
}

// Streams the formulae JSON from the Homebrew API, sending progress updates.
func startDownload(ch chan<- tea.Msg) {
	go func() {
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Get("https://formulae.brew.sh/api/formula.json")
		if err != nil {
			ch <- brewFormulaeErrMsg(err)
			return
		}
		defer resp.Body.Close() // closes the http response after the rest of the function is done

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

// Receives the next message from the download channel.
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

// Returns a model with loading state and a buffered download channel.
func New() Model {
	return Model{
		loading:    true,
		downloadCh: make(chan tea.Msg, 200),
	}
}

// Starts fetching the brew list and formulae data concurrently.
func (m Model) Init() tea.Cmd {
	startDownload(m.downloadCh)
	return tea.Batch(fetchBrewList(), m.recvDownload())
}

// Sequential character-level fuzzy matching.
func fuzzyMatch(s, query string) bool {
	// s = the name of the package
	// q/query = the search query being typed

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

// Filters displayPackages by search query and updates the info panel.
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

// Sets m.info to the formula data at the current cursor position.
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

// Handles all events: window resize, data from background goroutines, keyboard.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case formulaeProgressMsg:
		m.formulaeCount = msg.Count
		return m, m.recvDownload()

	case brewListMsg:
		m.packages = msg.names
		m.displayPackages = m.packages
		m.installPaths = msg.paths
		m.installedVersions = msg.installedVersions
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
				return m, nil
			case "enter":
				m.searchActive = false
				return m, nil
			case "backspace":
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
					m = m.applyFilter()
				}
				return m, nil
			case "up":
				if m.cursor > 0 {
					m.cursor--
					m = m.updateInfo()
				}
				return m, nil
			case "down":
				if m.cursor < len(m.displayPackages)-1 {
					m.cursor++
					m = m.updateInfo()
				}
				return m, nil
			default:
				if len(msg.String()) == 1 && msg.String()[0] >= 32 {
					m.searchQuery += msg.String()
					m = m.applyFilter()
				}
				return m, nil
			}
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

// Renders the complete terminal UI.
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
	leftWidth := int(float64(contentWidth) * 0.35)
	rightWidth := contentWidth - leftWidth

	searchLine := m.renderSearchBar(contentWidth, m.searchActive)
	searchOffset := strings.Count(searchLine, "\n") + 1

	boxHeight := m.height - 4 - searchOffset

	leftPanel := m.renderLeftPanel(leftWidth, boxHeight)
	rightPanel := m.renderRightPanel(rightWidth)

	leftStyled := lipgloss.NewStyle().Width(leftWidth).Height(boxHeight).Render(leftPanel)
	rightStyled := lipgloss.NewStyle().Width(rightWidth).Height(boxHeight).Render(rightPanel)

	top := lipgloss.JoinHorizontal(lipgloss.Top, leftStyled, rightStyled)

	var bodyParts []string
	bodyParts = append(bodyParts, searchLine)
	bodyParts = append(bodyParts, "")
	bodyParts = append(bodyParts, top)
	bodyParts = append(bodyParts, m.renderFooter())

	body := lipgloss.JoinVertical(lipgloss.Left, bodyParts...)
	return docStyle.Render(body)
}

// Loading screen shown while formulae are being fetched.
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

// Draws the search input with a bordered frame and cursor.
func (m Model) renderSearchBar(width int, focused bool) string {
	borderColor := teal
	if !focused {
		borderColor = tealDark
	}
	border := lipgloss.NewStyle().Foreground(borderColor)
	amberBold := lipgloss.NewStyle().Bold(true).Foreground(amber)

	badge := "Search"
	top := border.Render("╭─ ") +
		amberBold.Render(badge) +
		border.Render(" "+strings.Repeat("─", max(0, width-5-lipgloss.Width(badge)))+"╮")

	inner := width - 4

	var inputLine string
	if focused {
		cursor := lipgloss.NewStyle().Foreground(amber).Render("█")
		if m.searchQuery == "" {
			inputLine = cursor + " " + SearchPlaceholderStyle.Render("Search packages...")
		} else {
			inputLine = DetailValueStyle.Render(m.searchQuery) + cursor
		}
	} else {
		if m.searchQuery == "" {
			inputLine = SearchPlaceholderStyle.Render("Search packages...")
		} else {
			inputLine = DetailValueStyle.Render(m.searchQuery)
		}
	}

	padded := lipgloss.NewStyle().Width(inner).Render(inputLine)
	body := border.Render("│ ") + padded + border.Render(" │")
	bottom := border.Render("╰" + strings.Repeat("─", width-2) + "╯")

	return strings.Join([]string{top, body, bottom}, "\n")
}

// Scrollable package list on the left side.
func (m Model) renderLeftPanel(width int, boxHeight int) string {
	visibleHeight := boxHeight - 2

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

	boxTitle := fmt.Sprintf("Installed Packages (%d)", len(m.displayPackages))
	return renderPaneBox(width, boxTitle, strings.Join(listItems, "\n"))
}

// Detail panel on the right with formula details
func (m Model) renderRightPanel(width int) string {
	if len(m.displayPackages) == 0 {
		return renderPaneBox(width, "Details",
			lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#c0d4e4")).Render("No packages match your query"))
	}

	pkgName := m.displayPackages[m.cursor]
	title := DetailTitleStyle.Render("📦 " + pkgName)

	var contentLines []string
	contentLines = append(contentLines, title)

	if m.apiErr != nil {
		contentLines = append(contentLines,
			DetailValueStyle.Render("  Formula data unavailable"))
	} else if m.info != nil {
		info := m.info

		// Description
		if info.Desc != "" {
			descStyle := lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("#c0d4e4"))
			contentLines = append(contentLines, descStyle.Render(info.Desc))
		}

		type sectionData struct {
			title string
			lines []string
		}
		var sections []sectionData
		var allWidths []int

		// Package section
		var pkgPairs [][2]string
		if ver, ok := m.installedVersions[pkgName]; ok {
			pkgPairs = append(pkgPairs, [2]string{"Installed", ver})
		}
		if info.Versions.Stable != "" {
			pkgPairs = append(pkgPairs, [2]string{"Latest", info.Versions.Stable})
		}
		if path, ok := m.installPaths[pkgName]; ok && path != "" {
			pkgPairs = append(pkgPairs, [2]string{"Path", path})
		}
		if len(pkgPairs) > 0 {
			maxLabel := 0
			for _, p := range pkgPairs {
				w := lipgloss.Width(p[0])
				if w > maxLabel {
					maxLabel = w
				}
			}
			var lines []string
			for _, p := range pkgPairs {
				label := lipgloss.NewStyle().Width(maxLabel).Bold(true).Foreground(teal).Render(p[0])
				value := DetailValueStyle.Render(p[1])
				line := label + "  " + value
				allWidths = append(allWidths, lipgloss.Width(line))
				lines = append(lines, line)
			}
			sections = append(sections, sectionData{"Package", lines})
		}

		// Metadata section
		var metaPairs [][2]string
		if info.License != "" {
			metaPairs = append(metaPairs, [2]string{"License", info.License})
		}

		if info.Homepage != "" {
			metaPairs = append(metaPairs, [2]string{"Homepage", info.Homepage})
		}

		if len(metaPairs) > 0 {
			maxLabel := 0
			for _, p := range metaPairs {
				w := lipgloss.Width(p[0])
				if w > maxLabel {
					maxLabel = w
				}
			}

			var lines []string
			for _, p := range metaPairs {
				label := lipgloss.NewStyle().Width(maxLabel).Bold(true).Foreground(teal).Render(p[0])
				var value string
				if p[0] == "Homepage" {
					value = LinkStyle.Render(p[1])
				} else {
					value = DetailValueStyle.Render(p[1])
				}
				line := label + "  " + value
				allWidths = append(allWidths, lipgloss.Width(line))
				lines = append(lines, line)
			}
			sections = append(sections, sectionData{"Metadata", lines})
		}

		// Dependencies section
		if len(info.Dependencies) > 0 {
			line := DetailValueStyle.Render(strings.Join(info.Dependencies, ", "))
			allWidths = append(allWidths, lipgloss.Width(line))
			sections = append(sections, sectionData{"Dependencies", []string{line}})
		}

		// Build Dependencies section
		if len(info.BuildDependencies) > 0 {
			line := DetailValueStyle.Render(strings.Join(info.BuildDependencies, ", "))
			allWidths = append(allWidths, lipgloss.Width(line))
			sections = append(sections, sectionData{"Build Dependencies", []string{line}})
		}

		sectionWidth := width
		if len(allWidths) > 0 {
			maxW := 0
			for _, w := range allWidths {
				if w > maxW {
					maxW = w
				}
			}

			sectionWidth = min(width, max(maxW+4, 6))
		}

		// Divider after description
		hasContent := len(sections) > 0 || info.Desc != ""
		if hasContent {
			contentLines = append(contentLines, "")
			contentLines = append(contentLines,
				lipgloss.NewStyle().Foreground(teal).
					Render(strings.Repeat("─", max(0, width-4))))
			contentLines = append(contentLines, "")
		}

		// Render sections
		for _, s := range sections {
			contentLines = append(contentLines, renderSection(sectionWidth, s.title, s.lines...))
		}

	} else {
		contentLines = append(contentLines,
			DetailValueStyle.Render("  No formula data available"))
	}

	return renderPaneBox(width, "Details", strings.Join(contentLines, "\n"))
}

// Status bar with keybindings and error state.
func (m Model) renderFooter() string {
	apiErrMsg := ""
	if m.apiErr != nil {
		apiErrMsg += "  " + ErrorStyle.Render("API unavailable")
	}
	help := FooterStyle.Render("/ search  •  ↑↓ navigate  •  q quit")
	return apiErrMsg + "  " + help
}

// Single-column list for terminals narrower than 60 columns.
func (m Model) listViewFallback() string {
	title := TitleStyle.Render(fmt.Sprintf("pkgui  (%d)", len(m.packages)))

	sep := lipgloss.NewStyle().
		Foreground(teal).
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

	body := lipgloss.JoinVertical(lipgloss.Left, title, sep, list)
	return docStyle.Render(body)
}
