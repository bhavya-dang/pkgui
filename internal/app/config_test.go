package app

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigDefault(t *testing.T) {
	cfg := loadConfig()
	if cfg.Theme == "" {
		t.Error("loadConfig() returned empty theme")
	}
}

func TestConfigRoundTrip(t *testing.T) {
	dir := t.TempDir()

	origConfigPath := configPath
	configPath = func() string {
		return filepath.Join(dir, "pkgui", "config.json")
	}
	defer func() { configPath = origConfigPath }()

	applyTheme(themes[0])

	cmd := saveConfigCmd()
	msg := cmd()
	if msg != nil {
		t.Errorf("saveConfigCmd returned non-nil msg: %v", msg)
	}

	data, err := os.ReadFile(filepath.Join(dir, "pkgui", "config.json"))
	if err != nil {
		t.Fatal("config file was not written:", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatal("config file is not valid JSON:", err)
	}

	if cfg.Theme != themes[0].Name {
		t.Errorf("saved theme = %q, want %q", cfg.Theme, themes[0].Name)
	}
}

func TestConfigLoadAfterSave(t *testing.T) {
	dir := t.TempDir()

	origConfigPath := configPath
	configPath = func() string {
		return filepath.Join(dir, "pkgui", "config.json")
	}
	defer func() { configPath = origConfigPath }()

	applyTheme(themes[1])
	saveConfigCmd()()

	loadedCfg := loadConfig()
	if loadedCfg.Theme != themes[1].Name {
		t.Errorf("loaded theme = %q, want %q", loadedCfg.Theme, themes[1].Name)
	}
}

func TestConfigLoadInvalidFile(t *testing.T) {
	dir := t.TempDir()

	cfgPath := filepath.Join(dir, "pkgui", "config.json")
	os.MkdirAll(filepath.Dir(cfgPath), 0755)
	os.WriteFile(cfgPath, []byte("invalid json"), 0644)

	origConfigPath := configPath
	configPath = func() string { return cfgPath }
	defer func() { configPath = origConfigPath }()

	cfg := loadConfig()
	if cfg.Theme == "" {
		t.Error("loadConfig() should return default theme on invalid file")
	}
}
