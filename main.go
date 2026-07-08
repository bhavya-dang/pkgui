package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/bhavya-dang/pkgui/internal/app"
)

func main() {
	p := tea.NewProgram(app.New(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Print("\033[2J\033[H")
		log.Fatal(err)
	}

	fmt.Print("\033[2J\033[H")
}
