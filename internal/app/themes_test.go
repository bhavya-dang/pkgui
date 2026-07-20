package app

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestThemesNotEmpty(t *testing.T) {
	if len(themes) == 0 {
		t.Fatal("themes slice should not be empty")
	}
}

func TestThemesHaveRequiredFields(t *testing.T) {
	for i, th := range themes {
		if th.Name == "" {
			t.Errorf("themes[%d] has empty Name", i)
		}
		if th.Description == "" {
			t.Errorf("themes[%d] has empty Description", i)
		}
		if string(th.Primary) == "" {
			t.Errorf("themes[%d] (%q) has empty Primary", i, th.Name)
		}
		if string(th.Text) == "" {
			t.Errorf("themes[%d] (%q) has empty Text", i, th.Name)
		}
		if string(th.SelectedBg) == "" {
			t.Errorf("themes[%d] (%q) has empty SelectedBg", i, th.Name)
		}
		if string(th.SelectedFg) == "" {
			t.Errorf("themes[%d] (%q) has empty SelectedFg", i, th.Name)
		}
	}
}

func TestNoDuplicateThemeNames(t *testing.T) {
	seen := make(map[string]int)
	for i, th := range themes {
		if idx, ok := seen[th.Name]; ok {
			t.Errorf("duplicate theme name %q at indices %d and %d", th.Name, idx, i)
		}
		seen[th.Name] = i
	}
}

func TestThemeCount(t *testing.T) {
	if len(themes) < 3 {
		t.Error("expected at least 3 themes")
	}
}

func TestApplyThemeSetsCurrent(t *testing.T) {
	for _, th := range themes {
		applyTheme(th)
		if currentTheme != th {
			t.Errorf("applyTheme(%q) did not set currentTheme", th.Name)
		}
	}
}

func TestApplyThemeSetsPrimaryColor(t *testing.T) {
	applyTheme(themes[0])
	if string(violet) == "" {
		t.Error("violet color not set after applyTheme")
	}
	expected := string(themes[0].Primary)
	if string(violet) != expected {
		t.Errorf("violet = %q, want %q", string(violet), expected)
	}
}

func TestApplyThemeTogglesBetweenThemes(t *testing.T) {
	if len(themes) < 2 {
		t.Skip("need at least 2 themes")
	}
	applyTheme(themes[0])
	firstPrimary := string(violet)

	applyTheme(themes[1])
	secondPrimary := string(violet)

	if firstPrimary == secondPrimary {
		t.Error("theme colors should change after applyTheme with different theme")
	}
}

func TestStylesRender(t *testing.T) {
	applyTheme(themes[0])

	tests := []struct {
		name  string
		style lipgloss.Style
	}{
		{"TitleStyle", TitleStyle},
		{"ItemStyle", ItemStyle},
		{"SelectedItemStyle", SelectedItemStyle},
		{"ResultStyle", ResultStyle},
		{"FooterStyle", FooterStyle},
		{"LoadingStyle", LoadingStyle},
		{"ErrorStyle", ErrorStyle},
		{"DetailTitleStyle", DetailTitleStyle},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rendered := tt.style.Render("test")
			if rendered == "" {
				t.Errorf("%s rendered empty string", tt.name)
			}
		})
	}
}


