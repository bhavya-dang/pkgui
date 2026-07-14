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

type PipManager struct {
	tabIndex int
}

func NewPipManager(tabIndex int) *PipManager {
	return &PipManager{tabIndex: tabIndex}
}

func (p *PipManager) Name() string {
	return "pip"
}

func (p *PipManager) TabLabel() string {
	return "pip"
}

func (p *PipManager) ListInstalled() tea.Cmd {
	return func() tea.Msg {
		cmd, args := p.resolveCmd()
		if cmd == "" {
			return PackageListMsg{TabIndex: p.tabIndex}
		}
		args = append(args, "list", "--format=json")
		out, err := exec.Command(cmd, args...).Output()
		if err != nil {
			return PackageListMsg{TabIndex: p.tabIndex}
		}

		var result []struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		}
		if err := json.Unmarshal(out, &result); err != nil {
			return PackageListMsg{Err: err, TabIndex: p.tabIndex}
		}

		names := make([]string, 0, len(result))
		versions := make(map[string]string, len(result))
		for _, pkg := range result {
			names = append(names, pkg.Name)
			versions[pkg.Name] = pkg.Version
		}

		return PackageListMsg{Packages: names, Versions: versions, TabIndex: p.tabIndex}
	}
}

func (p *PipManager) resolveCmd() (string, []string) {
	if _, err := exec.LookPath("pip"); err == nil {
		return "pip", nil
	}
	if _, err := exec.LookPath("pip3"); err == nil {
		return "pip3", nil
	}
	if _, err := exec.LookPath("python3"); err == nil {
		return "python3", []string{"-m", "pip"}
	}
	if _, err := exec.LookPath("python"); err == nil {
		return "python", []string{"-m", "pip"}
	}
	return "", nil
}

type PipDetailData struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Summary     string `json:"summary"`
	License     string `json:"license"`
	HomePage    string `json:"home_page"`
	Author      string `json:"author"`
	AuthorEmail string `json:"author_email"`
}

type PipAllDetailsMsg map[string]*PipDetailData

func FetchAllPipDetails(names []string) tea.Cmd {
	return func() tea.Msg {
		results := make(map[string]*PipDetailData, len(names))
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

				url := fmt.Sprintf("https://pypi.org/pypi/%s/json", n)
				resp, err := client.Get(url)
				if err != nil {
					return
				}
				defer resp.Body.Close()

				var pkgResp struct {
					Info PipDetailData `json:"info"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&pkgResp); err != nil {
					return
				}
				mu.Lock()
				results[n] = &pkgResp.Info
				mu.Unlock()
			}(name)
		}
		wg.Wait()
		return PipAllDetailsMsg(results)
	}
}
