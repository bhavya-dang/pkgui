package cli

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/bhavya-dang/pkgui/internal/app"
	"github.com/spf13/cobra"
)

var configDir string

var rootCmd = &cobra.Command{
	Use:   "pkgui",
	Short: "terminal dashboard for installed packages",
	RunE:  run,
}

func SetVersion(v string) {
	rootCmd.Version = v
}

func init() {
	rootCmd.Flags().StringVar(&configDir, "config", "", "config directory path (default $HOME/pkgui)")
}

func run(cmd *cobra.Command, args []string) error {
	if configDir != "" {
		app.ConfigDir = configDir
	}
	p := tea.NewProgram(app.New(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Print("\033[2J\033[H")
		log.Fatal(err)
	}

	fmt.Print("\033[2J\033[H")
	return nil
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
