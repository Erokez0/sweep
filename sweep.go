package main

import (
	"fmt"

	config "sweep/config"
	gametui "sweep/tui/game-tui"
	startscreen "sweep/tui/start-screen"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	conf := config.GetConfig()
	for {
		if conf.Height == 0 || conf.Mines == 0 || conf.Width == 0 {
			startScreen := startscreen.CreateModel(conf)

			if _, err := tea.NewProgram(startScreen, tea.WithAltScreen()).Run(); err != nil {
				panic(fmt.Sprintf("could not start startscreen TUI: %v", err))
			}
		}

		gameModel := gametui.CreateModel(conf)

		if _, err := tea.NewProgram(gameModel, tea.WithAltScreen()).Run(); err != nil {
			panic(fmt.Sprintf("could not start game TUI: %v", err))
		}
	}
}
