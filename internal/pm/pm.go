// common interface for multiple package managers.
package pm

import tea "github.com/charmbracelet/bubbletea"

// the interface that each package manager backend implements.
type Manager interface {
	Name() string
	TabLabel() string
	ListInstalled() tea.Cmd
}

// carries the list of installed packages from a manager.
type PackageListMsg struct {
	Packages []string
	Versions map[string]string
	Err      error
	TabIndex int
}
