package main

import (
	"fmt"
	"os"

	app "github.com/andreiyard/go-happynewyear-ssh-app/internal/bubbletea"
	"github.com/andreiyard/go-happynewyear-ssh-app/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	cfg := config.Load()
	m := app.NewModel(cfg.Fps, cfg.SnowflakeRate, cfg.SnowflakeLimitPercent, cfg.SnowflakeChars)

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Println("Uh oh:", err)
		os.Exit(1)
	}
}
