package main

import (
	"fmt"
	"os"

	app "github.com/andreiyard/go-happynewyear-ssh-app/internal/bubbletea"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := app.NewModel([]string{"*", "+", "."})

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}
}
