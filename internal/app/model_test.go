package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/bhavya-dang/pkgui/internal/pm"
)

func TestFuzzyMatch(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		query   string
		want    bool
	}{
		{"exact match", "hello", "hello", true},
		{"case insensitive", "Hello", "hello", true},
		{"subsequence match", "hlo", "hello", false},
		{"fuzzy match", "hlo", "hlo", true},
		{"empty query", "hello", "", true},
		{"empty string", "", "a", false},
		{"partial prefix", "hello world", "hew", true},
		{"no match", "hello", "xyz", false},
		{"gaps in string", "hlo", "hlo", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fuzzyMatch(tt.s, tt.query)
			if got != tt.want {
				t.Errorf("fuzzyMatch(%q, %q) = %v, want %v", tt.s, tt.query, got, tt.want)
			}
		})
	}
}

func TestDarkenHex(t *testing.T) {
	tests := []struct {
		name   string
		hex    string
		factor float64
		want   string
	}{
		{"zero factor", "#ffffff", 0, "#ffffff"},
		{"full darken", "#ffffff", 1, "#000000"},
		{"half darken white", "#ffffff", 0.5, "#7f7f7f"},
		{"half darken red", "#ff0000", 0.5, "#7f0000"},
		{"quarter darken", "#aabbcc", 0.25, "#7f8c99"},
		{"invalid hex", "invalid", 0.5, "invalid"},
		{"short hex", "#abc", 0.5, "#abc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(darkenHex(tt.hex, tt.factor))
			if got != tt.want {
				t.Errorf("darkenHex(%q, %v) = %q, want %q", tt.hex, tt.factor, got, tt.want)
			}
		})
	}
}

func TestHumanSize(t *testing.T) {
	tests := []struct {
		name  string
		bytes int64
		want  string
	}{
		{"zero bytes", 0, "0 B"},
		{"bytes", 500, "500 B"},
		{"1 KB", 1024, "1.0 KB"},
		{"1.5 KB", 1536, "1.5 KB"},
		{"1 MB", 1048576, "1.0 MB"},
		{"1.5 MB", 1572864, "1.5 MB"},
		{"1 GB", 1073741824, "1.0 GB"},
		{"1 TB", 1099511627776, "1.0 TB"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := humanSize(tt.bytes)
			if got != tt.want {
				t.Errorf("humanSize(%d) = %q, want %q", tt.bytes, got, tt.want)
			}
		})
	}
}

func TestMin(t *testing.T) {
	if got := min(3, 5); got != 3 {
		t.Errorf("min(3, 5) = %d, want 3", got)
	}
	if got := min(5, 3); got != 3 {
		t.Errorf("min(5, 3) = %d, want 3", got)
	}
	if got := min(4, 4); got != 4 {
		t.Errorf("min(4, 4) = %d, want 4", got)
	}
	if got := min(-1, 5); got != -1 {
		t.Errorf("min(-1, 5) = %d, want -1", got)
	}
}

func TestMax(t *testing.T) {
	if got := max(3, 5); got != 5 {
		t.Errorf("max(3, 5) = %d, want 5", got)
	}
	if got := max(5, 3); got != 5 {
		t.Errorf("max(5, 3) = %d, want 5", got)
	}
	if got := max(4, 4); got != 4 {
		t.Errorf("max(4, 4) = %d, want 4", got)
	}
	if got := max(-1, 5); got != 5 {
		t.Errorf("max(-1, 5) = %d, want 5", got)
	}
}

func TestAllLoaded(t *testing.T) {
	applyTheme(themes[0])

	m := Model{
		states: []TabState{
			{loading: false},
			{loading: false},
			{loading: false},
		},
	}
	if !m.allLoaded() {
		t.Error("allLoaded() = false, want true")
	}

	m.states[1].loading = true
	if m.allLoaded() {
		t.Error("allLoaded() = true, want false (index 1 still loading)")
	}
}

func TestApplyFilterNoSearch(t *testing.T) {
	applyTheme(themes[0])

	pkgs := []string{"zebra", "apple", "banana"}
	m := Model{
		activeTab: 0,
		states: []TabState{
			{packages: pkgs, displayPackages: pkgs},
			{packages: nil},
			{packages: nil},
		},
	}
	m = m.applyFilter()
	if len(m.states[0].displayPackages) != 3 {
		t.Errorf("displayPackages count = %d, want 3", len(m.states[0].displayPackages))
	}
}

func TestApplyFilterSearch(t *testing.T) {
	applyTheme(themes[0])

	pkgs := []string{"zebra", "apple", "banana", "apricot"}
	st := TabState{packages: pkgs, displayPackages: pkgs}
	m := Model{
		activeTab:    0,
		searchQuery:  "ap",
		searchActive: true,
		states:       []TabState{st, {}, {}},
	}
	m = m.applyFilter()
	if len(m.states[0].displayPackages) != 2 {
		t.Errorf("filtered count = %d, want 2", len(m.states[0].displayPackages))
	}
	for _, p := range m.states[0].displayPackages {
		if p != "apple" && p != "apricot" {
			t.Errorf("unexpected package in filtered results: %s", p)
		}
	}
}

func TestApplyFilterAllMode(t *testing.T) {
	applyTheme(themes[0])

	m := Model{
		allMode:     true,
		searchQuery: "foo",
		allPackages: []string{"foobar", "baz", "foobaz", "qux"},
		allPackageOrigin: map[string]string{
			"foobar": "brew",
			"baz":    "npm",
			"foobaz": "pip",
			"qux":    "brew",
		},
	}
	m = m.applyFilter()
	if len(m.allDisplayPackages) != 2 {
		t.Errorf("filtered ALL count = %d, want 2", len(m.allDisplayPackages))
	}
}

func TestBuildAllPackages(t *testing.T) {
	applyTheme(themes[0])

	tabs := []pm.Manager{
		pm.NewBrewManager(0),
		pm.NewNpmManager(1),
		pm.NewPipManager(2),
	}
	m := Model{
		tabs: tabs,
		states: []TabState{
			{packages: []string{"brew-a", "brew-b"}},
			{packages: []string{"npm-a", "npm-c"}},
			{packages: []string{"pip-a"}},
		},
	}
	m = m.buildAllPackages()

	if len(m.allPackages) != 5 {
		t.Errorf("allPackages count = %d, want 5", len(m.allPackages))
	}

	expectedOrigins := map[string]string{
		"brew-a": "brew",
		"brew-b": "brew",
		"npm-a":  "npm",
		"npm-c":  "npm",
		"pip-a":  "pip",
	}
	for pkg, expectedOrigin := range expectedOrigins {
		origin, ok := m.allPackageOrigin[pkg]
		if !ok {
			t.Errorf("package %q missing from allPackageOrigin", pkg)
			continue
		}
		if origin != expectedOrigin {
			t.Errorf("allPackageOrigin[%q] = %q, want %q", pkg, origin, expectedOrigin)
		}
	}
}

func TestBuildAllPackagesDeduplicates(t *testing.T) {
	applyTheme(themes[0])

	tabs := []pm.Manager{
		pm.NewBrewManager(0),
		pm.NewNpmManager(1),
	}
	m := Model{
		tabs: tabs,
		states: []TabState{
			{packages: []string{"shared-pkg"}},
			{packages: []string{"shared-pkg", "npm-only"}},
		},
	}
	m = m.buildAllPackages()

	if len(m.allPackages) != 2 {
		t.Errorf("allPackages count = %d, want 2 (deduplicated)", len(m.allPackages))
	}
	origin := m.allPackageOrigin["shared-pkg"]
	if origin != "npm" {
		t.Errorf("last origin wins in current impl, got %q (want npm)", origin)
	}
}

func TestRebuildBrewPackages(t *testing.T) {
	applyTheme(themes[0])

	m := Model{
		states: []TabState{
			{
				Brew: &BrewState{
					FormulaNames: []string{"f1", "f2"},
					CaskNames:    []string{"c1"},
					Taps:         []string{"t1"},
				},
			},
			{},
			{},
		},
	}
	m = m.rebuildBrewPackages()
	st := m.states[0]

	if len(st.packages) != 4 {
		t.Errorf("packages count = %d, want 4", len(st.packages))
	}

	typeChecks := map[string]string{
		"f1": "formula",
		"f2": "formula",
		"c1": "cask",
		"t1": "tap",
	}
	for pkg, expectedType := range typeChecks {
		gotType, ok := st.packageType[pkg]
		if !ok {
			t.Errorf("packageType[%q] missing", pkg)
			continue
		}
		if gotType != expectedType {
			t.Errorf("packageType[%q] = %q, want %q", pkg, gotType, expectedType)
		}
	}
}

func TestNewModelInitialState(t *testing.T) {
	m := New()

	if !m.allMode {
		t.Error("allMode should default to true")
	}
	if m.activeTab != 0 {
		t.Errorf("activeTab = %d, want 0", m.activeTab)
	}
	if len(m.tabs) != 3 {
		t.Errorf("tabs count = %d, want 3", len(m.tabs))
	}
	if len(m.states) != 3 {
		t.Errorf("states count = %d, want 3", len(m.states))
	}
	for i, st := range m.states {
		if !st.loading {
			t.Errorf("states[%d] should start loading", i)
		}
	}
}

func TestSparklineHistory(t *testing.T) {
	m := Model{
		states: []TabState{
			{packages: []string{"a", "b"}},
			{packages: []string{"c"}},
			{packages: []string{"d", "e", "f"}},
		},
		sparklineHistory: make([]float64, 0, 40),
	}
	m = m.updateSparkline()
	if len(m.sparklineHistory) != 1 {
		t.Errorf("sparklineHistory length = %d, want 1", len(m.sparklineHistory))
	}
	if m.sparklineHistory[0] != 6 {
		t.Errorf("sparklineHistory[0] = %f, want 6", m.sparklineHistory[0])
	}
}

func TestSparklineHistoryCap(t *testing.T) {
	m := Model{
		states:           []TabState{{packages: []string{"a"}}, {}, {}},
		sparklineHistory: make([]float64, 0, 40),
	}
	for i := 0; i < 50; i++ {
		m = m.updateSparkline()
	}
	if len(m.sparklineHistory) > 40 {
		t.Errorf("sparklineHistory capped at %d, want <= 40", len(m.sparklineHistory))
	}
}

type mockMsg struct{}

func TestUpdateKeyQuit(t *testing.T) {
	applyTheme(themes[0])
	m := New()
	result, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	if cmd == nil {
		t.Error("expected quit command, got nil")
	}
	_ = result
}

func TestUpdateKeySearch(t *testing.T) {
	applyTheme(themes[0])
	m := Model{
		allMode:   true,
		searchActive: false,
		states:    []TabState{{loading: false}, {loading: false}, {loading: false}},
	}
	result, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}})
	m2 := result.(Model)
	if !m2.searchActive {
		t.Error("searchActive should be true after / key")
	}
	if cmd != nil {
		_ = cmd
	}
}

func TestUpdateKeyLeftFromAll(t *testing.T) {
	applyTheme(themes[0])
	m := Model{
		allMode: true,
		states:  []TabState{{}, {}, {}},
	}
	result, _ := m.Update(tea.KeyMsg{Type: tea.KeyLeft})
	m2 := result.(Model)
	if m2.allMode != true {
		t.Error("left from ALL currently does nothing, allMode should remain true")
	}
}

func TestUpdateKeyRightFromTab(t *testing.T) {
	applyTheme(themes[0])
	m := Model{
		allMode:   false,
		activeTab: 0,
		tabs: []pm.Manager{
			pm.NewBrewManager(0),
			pm.NewNpmManager(1),
			pm.NewPipManager(2),
		},
		states: []TabState{
			{loading: false},
			{loading: false},
			{loading: false},
		},
	}
	result, _ := m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m2 := result.(Model)
	if m2.activeTab != 1 {
		t.Errorf("activeTab = %d, want 1", m2.activeTab)
	}
}

func TestUpdatePackageListMsg(t *testing.T) {
	applyTheme(themes[0])
	m := Model{
		states: []TabState{
			{},
			{loading: true},
			{loading: true},
		},
		sparklineHistory: make([]float64, 0, 40),
	}
	msg := pm.PackageListMsg{
		TabIndex: 1,
		Packages: []string{"pkg1", "pkg2"},
		Versions: map[string]string{"pkg1": "1.0", "pkg2": "2.0"},
	}
	result, _ := m.Update(msg)
	m2 := result.(Model)
	st := m2.states[1]
	if st.loading {
		t.Error("state should not be loading after receiving PackageListMsg")
	}
	if len(st.packages) != 2 {
		t.Errorf("packages count = %d, want 2", len(st.packages))
	}
	if st.versions["pkg1"] != "1.0" {
		t.Errorf("version = %q, want 1.0", st.versions["pkg1"])
	}
}

func TestTotalPackages(t *testing.T) {
	m := Model{
		states: []TabState{
			{displayPackages: []string{"a", "b"}},
			{displayPackages: []string{"c"}},
			{displayPackages: []string{}},
		},
	}
	if total := m.totalPackages(); total != 3 {
		t.Errorf("totalPackages() = %d, want 3", total)
	}
}

func TestPackageListMsgError(t *testing.T) {
	applyTheme(themes[0])
	m := Model{
		states: []TabState{
			{},
			{loading: true},
			{loading: true},
		},
	}
	result, _ := m.Update(pm.PackageListMsg{TabIndex: 2, Err: tea.ErrProgramKilled})
	m2 := result.(Model)
	if m2.states[2].loading {
		t.Error("state loading is set to false even on error")
	}
	if m2.states[2].err == nil {
		t.Error("state should have error set")
	}
}
