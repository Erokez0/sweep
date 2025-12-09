package endscreen

import (
	"fmt"
	"os"
	"time"

	"sweep/shared/consts"
	styles "sweep/tui/styles"
	tilerenderer "sweep/tui/tile-renderer"
	types "sweep/types"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	gameEngine types.IGameEngine
	startTime  time.Time
}

func CreateModel(startTime time.Time, gameEngine types.IGameEngine) model {
	return model{
		startTime:  startTime,
		gameEngine: gameEngine,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tea.SetWindowTitle(consts.APP_NAME), tea.ClearScreen)
}

var _ tea.Model = model{}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "ctrl+c", "q":
			os.Exit(0)
		default:
			return m, tea.Quit
		}
	}
	return m, nil
}

func beautifyTimeDuration(duration time.Duration) string {
	milliseconds := int(duration.Milliseconds()) % 100
	seconds := int(duration.Seconds()) % 60
	minutes := int(duration.Minutes()) % 60

	return fmt.Sprintf("%d:%d,%d", minutes, seconds, milliseconds)
}

func (m model) View() string {
	var lines string

	win := true

	field := m.gameEngine.GetField()

	for x := range field {
		var line string
		for y := range field[x] {
			tile := field[x][y]
			if tile == types.OpenBomb {
				win = false
			}
			content := fmt.Sprint(m.gameEngine.CountNeighbouringBombs(types.Position{
				X: uint16(x),
				Y: uint16(y),
			}))
			line += tilerenderer.RenderTileByType(tile, content)
		}
		lines += "\n"
		if x == 0 {
			lines += styles.BorderTop.Render(line)
		} else if x == len(field)-1 {
			lines += styles.BorderBottom.Render(line)
		} else {
			lines += line
		}
	}
	var s string
	if win {
		s = "You won!"
	} else {
		s = "You lost!"
	}
	s += lines + "\n"
	timeSinceStart := time.Since(m.startTime)

	beautifiedTime := beautifyTimeDuration(timeSinceStart)
	tea.SetWindowTitle(fmt.Sprintf("%v - %v", consts.APP_NAME, beautifiedTime))
	s += fmt.Sprintf("time - %v", beautifiedTime)
	return styles.TableStyle.Render(s)

}
