package pm

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type BrewManager struct {
	tabIndex int
}

func NewBrewManager(tabIndex int) *BrewManager {
	return &BrewManager{tabIndex: tabIndex}
}

func (b *BrewManager) Name() string {
	return "brew"
}

func (b *BrewManager) TabLabel() string {
	return "Brew"
}

func (b *BrewManager) ListInstalled() tea.Cmd {
	return tea.Batch(
		b.fetchBrewList(),
		b.fetchFormulae(),
	)
}

type BrewListMsg struct {
	Names             []string
	Paths             map[string]string
	InstalledVersions map[string]string
	Sizes             map[string]int64
}

type BrewErrMsg error

type BrewFormulaeMsg map[string]FormulaData

type BrewFormulaeErrMsg error

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

func (b *BrewManager) fetchBrewList() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("brew", "list", "--formula", "--versions")
		out, err := cmd.Output()
		if err != nil {
			return BrewErrMsg(err)
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

		sizes := make(map[string]int64, len(names))
		for _, name := range names {
			if path, ok := paths[name]; ok {
				if out, err := exec.Command("du", "-skL", path).Output(); err == nil {
					fields := strings.Fields(string(out))
					if len(fields) > 0 {
						if kb, err := strconv.ParseInt(fields[0], 10, 64); err == nil {
							sizes[name] = kb * 1024
						}
					}
				}
			}
		}

		return BrewListMsg{names, paths, installedVersions, sizes}
	}
}

func (b *BrewManager) fetchFormulae() tea.Cmd {
	return func() tea.Msg {
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Get("https://formulae.brew.sh/api/formula.json")
		if err != nil {
			return BrewFormulaeErrMsg(err)
		}
		defer resp.Body.Close()

		dec := json.NewDecoder(resp.Body)

		_, err = dec.Token()
		if err != nil {
			return BrewFormulaeErrMsg(err)
		}

		m := make(map[string]FormulaData)

		for dec.More() {
			var f FormulaData
			if err := dec.Decode(&f); err != nil {
				return BrewFormulaeErrMsg(err)
			}
			m[f.Name] = f
		}

		return BrewFormulaeMsg(m)
	}
}
