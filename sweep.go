package main

import (
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

			if _, err := tea.NewProgram(startScreen).Run(); err != nil {
				panic(err)
			}
		}

		gameModel := gametui.CreateModel(conf)

		if _, err := tea.NewProgram(gameModel).Run(); err != nil {
			panic(err)
		}

	}

}
