package pm

import (
	"context"
	"encoding/json"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func brewOutput(name string, args ...string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, name, args...)
	return cmd.Output()
}

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
		b.fetchCaskList(),
		b.fetchTapList(),
	)
}

func (b *BrewManager) FetchCaskData() tea.Cmd {
	return b.fetchCaskData()
}

func (b *BrewManager) FetchFormulae() tea.Cmd {
	return b.fetchFormulae()
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

type BrewCaskListMsg struct {
	Names             []string
	Paths             map[string]string
	InstalledVersions map[string]string
	Sizes             map[string]int64
}

type BrewCaskErrMsg error

type BrewTapListMsg struct {
	Names []string
}

type BrewTapErrMsg error

type CaskData struct {
	Token    string   `json:"token"`
	Name     []string `json:"name"`
	Desc     string   `json:"desc"`
	Homepage string   `json:"homepage"`
	Version  string   `json:"version"`
}

type BrewCaskDataMsg map[string]*CaskData

type BrewCaskDataErrMsg error

type BrewTapFormulaeMsg struct {
	TapFormulae map[string][]string
}

type BrewTapFormulaeErrMsg error

func (b *BrewManager) fetchBrewList() tea.Cmd {
	return func() tea.Msg {
		out, err := brewOutput("brew", "list", "--formula", "--versions")
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
			prefixOut, perr := brewOutput("brew", "--prefix")
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
				sizeOut, serr := brewOutput("du", "-skL", path)
				if serr == nil {
					fields := strings.Fields(string(sizeOut))
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

func (b *BrewManager) fetchCaskList() tea.Cmd {
	return func() tea.Msg {
		out, err := brewOutput("brew", "list", "--casks", "--versions")
		if err != nil {
			return BrewCaskListMsg{nil, nil, nil, nil}
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
			prefixOut, perr := brewOutput("brew", "--prefix")
			if perr == nil {
				prefix := strings.TrimSpace(string(prefixOut))
				for _, name := range names {
					ver, ok := installedVersions[name]
					if ok {
						paths[name] = prefix + "/Caskroom/" + name + "/" + ver
					}
				}
			}
		}

		sizes := make(map[string]int64, len(names))
		for _, name := range names {
			if path, ok := paths[name]; ok {
				sizeOut, serr := brewOutput("du", "-skL", path)
				if serr == nil {
					fields := strings.Fields(string(sizeOut))
					if len(fields) > 0 {
						if kb, err := strconv.ParseInt(fields[0], 10, 64); err == nil {
							sizes[name] = kb * 1024
						}
					}
				}
			}
		}

		return BrewCaskListMsg{names, paths, installedVersions, sizes}
	}
}

func (b *BrewManager) fetchCaskData() tea.Cmd {
	return func() tea.Msg {
		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Get("https://formulae.brew.sh/api/cask.json")
		if err != nil {
			return BrewCaskDataErrMsg(err)
		}
		defer resp.Body.Close()

		dec := json.NewDecoder(resp.Body)

		_, err = dec.Token()
		if err != nil {
			return BrewCaskDataErrMsg(err)
		}

		m := make(map[string]*CaskData)

		for dec.More() {
			var c CaskData
			if err := dec.Decode(&c); err != nil {
				return BrewCaskDataErrMsg(err)
			}
			m[c.Token] = &c
		}

		return BrewCaskDataMsg(m)
	}
}

func (b *BrewManager) fetchTapList() tea.Cmd {
	return func() tea.Msg {
		out, err := brewOutput("brew", "tap")
		if err != nil {
			return BrewTapListMsg{nil}
		}
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		var names []string
		for _, line := range lines {
			if line != "" {
				names = append(names, line)
			}
		}
		return BrewTapListMsg{names}
	}
}

func (b *BrewManager) FetchTapFormulae(taps []string) tea.Cmd {
	return func() tea.Msg {
		result := make(map[string][]string)
		for _, tap := range taps {
			out, err := brewOutput("brew", "list", tap, "--versions")
			if err != nil {
				continue
			}
			lines := strings.Split(strings.TrimSpace(string(out)), "\n")
			var formulae []string
			for _, line := range lines {
				parts := strings.Fields(line)
				if len(parts) >= 1 {
					formulae = append(formulae, parts[0])
				}
			}
			result[tap] = formulae
		}
		return BrewTapFormulaeMsg{result}
	}
}
