package app

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Name        string
	Description string
	URL         string

	Primary    lipgloss.Color
	Muted      lipgloss.Color
	Text       lipgloss.Color
	DimText    lipgloss.Color
	DetailText lipgloss.Color
	SelectedBg lipgloss.Color
	SelectedFg lipgloss.Color
	Success    lipgloss.Color
	Error      lipgloss.Color
}

var themes = []*Theme{
	{
		Name:        "solace",
		Description: "calm violet-pastel palette",
		Primary:     lipgloss.Color("#C9A8FF"),
		Muted:       lipgloss.Color("#9880E8"),
		Text:        lipgloss.Color("#F1EEF8"),
		DimText:     lipgloss.Color("#A09AB8"),
		DetailText:  lipgloss.Color("#D8D3F0"),
		SelectedBg:  lipgloss.Color("#C9A8FF"),
		SelectedFg:  lipgloss.Color("#16141F"),
		Success:     lipgloss.Color("#A8FFEC"),
		Error:       lipgloss.Color("#FF80C8"),
	},
	{
		Name:        "gruvbox",
		Description: "warm earthy retro tones",
		Primary:     lipgloss.Color("#d79921"),
		Muted:       lipgloss.Color("#a89984"),
		Text:        lipgloss.Color("#ebdbb2"),
		DimText:     lipgloss.Color("#928374"),
		DetailText:  lipgloss.Color("#d5c4a1"),
		SelectedBg:  lipgloss.Color("#d79921"),
		SelectedFg:  lipgloss.Color("#282828"),
		Success:     lipgloss.Color("#98971a"),
		Error:       lipgloss.Color("#cc241d"),
	},
	{
		Name:        "nord",
		Description: "cold arctic blue palette",
		Primary:     lipgloss.Color("#88c0d0"),
		Muted:       lipgloss.Color("#4c566a"),
		Text:        lipgloss.Color("#eceff4"),
		DimText:     lipgloss.Color("#616e88"),
		DetailText:  lipgloss.Color("#d8dee9"),
		SelectedBg:  lipgloss.Color("#88c0d0"),
		SelectedFg:  lipgloss.Color("#2e3440"),
		Success:     lipgloss.Color("#a3be8c"),
		Error:       lipgloss.Color("#bf616a"),
	},
	{
		Name:        "vesper",
		Description: "bright amber glow",
		Primary:     lipgloss.Color("#f5b342"),
		Muted:       lipgloss.Color("#c0841a"),
		Text:        lipgloss.Color("#fdf6e3"),
		DimText:     lipgloss.Color("#b8a87a"),
		DetailText:  lipgloss.Color("#ede0c8"),
		SelectedBg:  lipgloss.Color("#f5b342"),
		SelectedFg:  lipgloss.Color("#1a1408"),
		Success:     lipgloss.Color("#8bc34a"),
		Error:       lipgloss.Color("#ef5350"),
	},
	{
		Name:        "wave",
		Description: "bright cyan and teal surf",
		Primary:     lipgloss.Color("#00d4c8"),
		Muted:       lipgloss.Color("#1a7a72"),
		Text:        lipgloss.Color("#e0f0ef"),
		DimText:     lipgloss.Color("#5a8a86"),
		DetailText:  lipgloss.Color("#b0d4d0"),
		SelectedBg:  lipgloss.Color("#00d4c8"),
		SelectedFg:  lipgloss.Color("#001a18"),
		Success:     lipgloss.Color("#4ecdc4"),
		Error:       lipgloss.Color("#ff6b6b"),
	},
}
