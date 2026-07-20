package app

import (
	"encoding/json"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

type Config struct {
	Theme string `json:"theme"`
}

var ConfigDir string

var configPath = func() string {
	dir := ConfigDir
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		dir = filepath.Join(home, "pkgui")
	}
	return filepath.Join(dir, "config.json")
}

func loadConfig() Config {
	cfg := Config{Theme: themes[0].Name}
	path := configPath()
	if path == "" {
		return cfg
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg
	}
	var loaded Config
	if err := json.Unmarshal(data, &loaded); err != nil {
		return cfg
	}
	if loaded.Theme != "" {
		cfg.Theme = loaded.Theme
	}
	return cfg
}

func saveConfigCmd() tea.Cmd {
	return func() tea.Msg {
		path := configPath()
		if path == "" {
			return nil
		}
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil
		}
		cfg := Config{Theme: currentTheme.Name}
		data, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return nil
		}
		os.WriteFile(path, data, 0644)
		return nil
	}
}
