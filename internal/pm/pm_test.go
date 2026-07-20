package pm

import (
	"errors"
	"testing"
)

func TestPackageListMsg(t *testing.T) {
	msg := PackageListMsg{
		Packages: []string{"a", "b"},
		Versions: map[string]string{"a": "1.0"},
		Err:      nil,
		TabIndex: 2,
	}
	if len(msg.Packages) != 2 {
		t.Errorf("Packages count = %d, want 2", len(msg.Packages))
	}
	if msg.Versions["a"] != "1.0" {
		t.Errorf("Version = %q, want 1.0", msg.Versions["a"])
	}
	if msg.TabIndex != 2 {
		t.Errorf("TabIndex = %d, want 2", msg.TabIndex)
	}
	if msg.Err != nil {
		t.Errorf("Err = %v, want nil", msg.Err)
	}
}

func TestPackageListMsgWithError(t *testing.T) {
	msg := PackageListMsg{
		Err: errors.New("test error"),
	}
	if msg.Err == nil {
		t.Error("Err should be set")
	}
}

func TestBrewManagerInterface(t *testing.T) {
	var m Manager = NewBrewManager(0)
	if m.Name() != "brew" {
		t.Errorf("Name() = %q, want brew", m.Name())
	}
	if m.TabLabel() != "Brew" {
		t.Errorf("TabLabel() = %q, want Brew", m.TabLabel())
	}
}

func TestNpmManagerInterface(t *testing.T) {
	var m Manager = NewNpmManager(1)
	if m.Name() != "npm" {
		t.Errorf("Name() = %q, want npm", m.Name())
	}
	if m.TabLabel() != "npm" {
		t.Errorf("TabLabel() = %q, want npm", m.TabLabel())
	}
}

func TestPipManagerInterface(t *testing.T) {
	var m Manager = NewPipManager(2)
	if m.Name() != "pip" {
		t.Errorf("Name() = %q, want pip", m.Name())
	}
	if m.TabLabel() != "pip" {
		t.Errorf("TabLabel() = %q, want pip", m.TabLabel())
	}
}

func TestBrewListMsg(t *testing.T) {
	msg := BrewListMsg{
		Names:             []string{"formula1", "formula2"},
		Paths:             map[string]string{"formula1": "/opt/formula1"},
		InstalledVersions: map[string]string{"formula1": "1.0"},
		Sizes:             map[string]int64{"formula1": 1024},
	}
	if len(msg.Names) != 2 {
		t.Errorf("Names count = %d, want 2", len(msg.Names))
	}
	if msg.Paths["formula1"] != "/opt/formula1" {
		t.Errorf("Path = %q, want /opt/formula1", msg.Paths["formula1"])
	}
}

func TestCaskData(t *testing.T) {
	c := CaskData{
		Token:    "test-cask",
		Name:     []string{"Test Cask"},
		Desc:     "A test cask",
		Homepage: "https://example.com",
		Version:  "1.0",
	}
	if c.Token != "test-cask" {
		t.Errorf("Token = %q, want test-cask", c.Token)
	}
	if c.Version != "1.0" {
		t.Errorf("Version = %q, want 1.0", c.Version)
	}
	if len(c.Name) != 1 || c.Name[0] != "Test Cask" {
		t.Errorf("Name = %v, want [Test Cask]", c.Name)
	}
}

func TestFormulaData(t *testing.T) {
	f := FormulaData{
		Name:     "test-formula",
		Desc:     "A test formula",
		Homepage: "https://example.com",
		License:  "MIT",
	}
	f.Versions.Stable = "2.0"
	f.Dependencies = []string{"dep1", "dep2"}
	f.BuildDependencies = []string{"build-dep1"}

	if f.Name != "test-formula" {
		t.Errorf("Name = %q, want test-formula", f.Name)
	}
	if f.Versions.Stable != "2.0" {
		t.Errorf("Version = %q, want 2.0", f.Versions.Stable)
	}
	if len(f.Dependencies) != 2 {
		t.Errorf("Dependencies count = %d, want 2", len(f.Dependencies))
	}
}

func TestNpmDetailData(t *testing.T) {
	d := NpmDetailData{
		Name:        "test-pkg",
		Version:     "1.0.0",
		Description: "A test package",
		License:     "MIT",
		Homepage:    "https://example.com",
		Dist:        &NpmDist{UnpackedSize: 1024},
	}
	if d.Name != "test-pkg" {
		t.Errorf("Name = %q, want test-pkg", d.Name)
	}
	if d.Dist.UnpackedSize != 1024 {
		t.Errorf("UnpackedSize = %d, want 1024", d.Dist.UnpackedSize)
	}
}

func TestPipDetailData(t *testing.T) {
	d := PipDetailData{
		Name:        "test-pkg",
		Version:     "1.0.0",
		Summary:     "A test package",
		License:     "MIT",
		HomePage:    "https://example.com",
		Author:      "Test Author",
		AuthorEmail: "author@example.com",
	}
	if d.Name != "test-pkg" {
		t.Errorf("Name = %q, want test-pkg", d.Name)
	}
	if d.Author != "Test Author" {
		t.Errorf("Author = %q, want Test Author", d.Author)
	}
}

func TestPipManagerResolveCmd(t *testing.T) {
	p := NewPipManager(0)
	cmd, args := p.resolveCmd()
	if cmd == "" {
		t.Log("no pip/python found on PATH (expected in CI), cmd is empty")
	} else {
		t.Logf("resolved pip command: %s %v", cmd, args)
	}
}
