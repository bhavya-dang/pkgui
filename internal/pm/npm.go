package pm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type NpmManager struct {
	tabIndex int
}

func NewNpmManager(tabIndex int) *NpmManager {
	return &NpmManager{tabIndex: tabIndex}
}

func (n *NpmManager) Name() string {
	return "npm"
}

func (n *NpmManager) TabLabel() string {
	return "npm"
}

func (n *NpmManager) ListInstalled() tea.Cmd {
	return func() tea.Msg {
		if _, err := exec.LookPath("npm"); err != nil {
			return PackageListMsg{TabIndex: n.tabIndex}
		}
		cmd := exec.Command("npm", "ls", "-g", "--depth=0", "--json")
		out, err := cmd.Output()
		if err != nil {
			return PackageListMsg{TabIndex: n.tabIndex}
		}

		var result struct {
			Dependencies map[string]struct {
				Version string `json:"version"`
			} `json:"dependencies"`
		}
		if err := json.Unmarshal(out, &result); err != nil {
			return PackageListMsg{Err: err, TabIndex: n.tabIndex}
		}

		names := make([]string, 0, len(result.Dependencies))
		versions := make(map[string]string, len(result.Dependencies))
		for name, dep := range result.Dependencies {
			names = append(names, name)
			versions[name] = dep.Version
		}

		return PackageListMsg{Packages: names, Versions: versions, TabIndex: n.tabIndex}
	}
}

type NpmDist struct {
	UnpackedSize int64 `json:"unpackedSize"`
}

type NpmDetailData struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	License     string   `json:"license"`
	Homepage    string   `json:"homepage"`
	Dist        *NpmDist `json:"dist,omitempty"`
}

type NpmAllDetailsMsg map[string]*NpmDetailData

func FetchAllNpmDetails(names []string) tea.Cmd {
	return func() tea.Msg {
		results := make(map[string]*NpmDetailData, len(names))
		var mu sync.Mutex
		sem := make(chan struct{}, 5)
		var wg sync.WaitGroup

		client := http.Client{Timeout: 10 * time.Second}
		for _, name := range names {
			wg.Add(1)
			go func(n string) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()

				url := fmt.Sprintf("https://registry.npmjs.org/%s/latest", n)
				resp, err := client.Get(url)
				if err != nil {
					return
				}
				defer resp.Body.Close()

				var data NpmDetailData
				if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
					return
				}
				mu.Lock()
				results[n] = &data
				mu.Unlock()
			}(name)
		}
		wg.Wait()
		return NpmAllDetailsMsg(results)
	}
}
