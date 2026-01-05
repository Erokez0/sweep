package endscreen

import (
	"fmt"
	"os"
	"strings"
	"time"

	misc "sweep/shared/consts/misc"
	tilecontent "sweep/shared/consts/tile-content"
	tiles "sweep/shared/consts/tiles"
	types "sweep/shared/types"
	"sweep/shared/utils"
	styles "sweep/tui/styles"
	tilerenderer "sweep/tui/tile-renderer"

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
	return tea.Batch(tea.SetWindowTitle(misc.AppName), tea.ClearScreen)
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

func (m model) View() string {
	var lines strings.Builder

	win := true

	field := m.gameEngine.GetField()

	width, height := m.gameEngine.GetWidth(), m.gameEngine.GetHeight()

	for row := range height {
		var line string
		y := height - 1 - row
		for col := range width {
			x := col

			position := types.Position{
				X: x,
				Y: y,
			}

			tile := m.gameEngine.GetTile(position)

			if tile == tiles.OpenMine {
				win = false
			}

			count := m.gameEngine.CountNeighbouringMines(types.Position{
				X: uint16(x),
				Y: uint16(y),
			})
			tileContent, err := tilecontent.FromNumber(count)
			if err != nil {
				panic(err)
			}
			line += tilerenderer.RenderTileByType(tile, tileContent)
		}

		lines.WriteRune('\n')
		if row == 0 {
			lines.WriteString(styles.BorderTop.Render(line))
		} else if row == uint16(len(field)-1) {
			lines.WriteString(styles.BorderBottom.Render(line))
		} else {
			lines.WriteString(line)
		}
	}

	var s strings.Builder
	if win {
		s.WriteString("You won!")
	} else {
		s.WriteString("You lost!")
	}

	s.WriteString(lines.String())
	s.WriteRune('\n')

	formattedDuration := utils.FormatTime(time.Since(m.startTime))

	fmt.Fprintf(&s, "time - %v", formattedDuration)
	return styles.TableStyle.Render(s.String())
}
