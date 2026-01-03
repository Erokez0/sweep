package gametui

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	config "sweep/config"
	gameengine "sweep/game-engine"
	actions "sweep/shared/consts/action"
	misc "sweep/shared/consts/misc"
	tilecontent "sweep/shared/consts/tile-content"
	tiles "sweep/shared/consts/tiles"
	types "sweep/shared/types"
	keyactionmap "sweep/shared/vars/key-action-map"
	endscreen "sweep/tui/end-screen"
	styles "sweep/tui/styles"
	tilerenderer "sweep/tui/tile-renderer"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	config         config.Config
	cursorPosition types.Position
	gameEngine     types.IGameEngine
	tiles          [][]tilecontent.TileContent
	startTime      time.Time
	moves          int16
	flags          int16
}

func CreateModel(config *config.Config) model {
	gameEngine := gameengine.GameEngine{}
	gameEngine.SetMineCount(config.Mines)
	gameEngine.SetFieldSize(config.Width, config.Height)

	tiles := make([][]tilecontent.TileContent, config.Width)
	for x := range config.Width {
		tiles[x] = make([]tilecontent.TileContent, config.Height)
		for y := range config.Height {
			tiles[x][y] = tilecontent.Empty
		}
	}

	return model{
		cursorPosition: types.Position{
			X: config.Width / 2,
			Y: config.Height / 2,
		},
		gameEngine: &gameEngine,
		tiles:      tiles,
		startTime:  time.Now(),
		moves:      0,
		config:     *config,
	}
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle(misc.APP_NAME)
}

var _ tea.Model = model{}

func (m model) openTile(position types.Position) {
	tileType := m.gameEngine.GetTile(position)
	switch tileType {
	case tiles.OutOfBounds:
		return
	case tiles.OpenSafe:
		m.openAroundOpenTile(position)
		return
	}

	x, y := position.GetCoords()

	m.gameEngine.OpenTile(position)

	if m.gameEngine.IsFinished() {
		defer tea.Quit()
		return
	}

	count := m.gameEngine.CountNeighbouringMines(position)
	tileContent, err := tilecontent.FromNumber(count)
	if err != nil {
		panic(err)
	}
	m.tiles[x][y] = tileContent

	if count == 0 {
		m.openSafeAroundTile(position)
	}
}

func (m model) openSafeAroundTile(position types.Position) {
	x, y := position.GetCoords()

	neighbours := []types.Position{
		{X: x - 1, Y: y - 1},
		{X: x - 1, Y: y},
		{X: x - 1, Y: y + 1},
		{X: x, Y: y - 1},
		{X: x, Y: y + 1},
		{X: x + 1, Y: y - 1},
		{X: x + 1, Y: y},
		{X: x + 1, Y: y + 1},
	}

	var wg sync.WaitGroup
	for _, neighbour := range neighbours {
		wg.Go(func() {
			switch m.gameEngine.GetTile(neighbour) {
			case tiles.ClosedSafe:
				m.openTile(neighbour)
			}
		})
	}
	wg.Wait()
}

func (m model) openAroundOpenTile(position types.Position) {
	x, y := position.GetCoords()

	neighbours := []types.Position{
		{X: x - 1, Y: y - 1},
		{X: x - 1, Y: y},
		{X: x - 1, Y: y + 1},
		{X: x, Y: y - 1},
		{X: x, Y: y + 1},
		{X: x + 1, Y: y - 1},
		{X: x + 1, Y: y},
		{X: x + 1, Y: y + 1},
	}

	var wg sync.WaitGroup
	for _, position := range neighbours {
		wg.Go(func() {
			switch m.gameEngine.GetTile(position) {
			case tiles.FlaggedSafe, tiles.FlaggedMine, tiles.OpenSafe, tiles.OutOfBounds:
				return
			}
			x, y := position.GetCoords()
			m.gameEngine.OpenTile(position)

			if m.gameEngine.IsFinished() {
				defer tea.Quit()
				return
			}

			count := m.gameEngine.CountNeighbouringMines(position)

			tileContent, err := tilecontent.FromNumber(count)
			if err != nil {
				panic(err)
			}
			m.tiles[x][y] = tileContent
			if count == 0 {
				m.openSafeAroundTile(position)
			}
		})
	}
	wg.Wait()
}

func (m model) MoveCursorUp() (model, tea.Cmd) {
	if m.cursorPosition.X > 0 {
		m.cursorPosition.X--
	}
	
	return m, nil
}
func (m model) MoveCursorDown() (model, tea.Cmd) {
	if m.cursorPosition.X < uint16(m.gameEngine.GetWidth())-1 {
		m.cursorPosition.X++
	}
	
	return m, nil
}
func (m model) MoveCursorRight() (model, tea.Cmd) {
	if m.cursorPosition.Y < uint16(m.gameEngine.GetHeight())-1 {
		m.cursorPosition.Y++
	}
	
	return m, nil
}
func (m model) MoveCursorLeft() (model, tea.Cmd) {
	if m.cursorPosition.Y > 0 {
		m.cursorPosition.Y--
	}
	return m, nil
}
func (m model) FlagTile() (model, tea.Cmd) {
	if m.moves == 0 {
		return m, nil
	}
	m.moves++

	x, y := m.cursorPosition.GetCoords()
	m.gameEngine.FlagToggleTile(m.cursorPosition)
	switch m.tiles[x][y] {
	case tilecontent.Empty:
		m.tiles[x][y] = tilecontent.Flag
		m.flags++
	case tilecontent.Flag:
		m.tiles[x][y] = tilecontent.Empty
		m.flags--
	}

	return m, nil
}
func (m model) OpenTile() (model, tea.Cmd) {
	if m.moves == 0 {
		m.gameEngine.SetMines(m.cursorPosition)
	}
	m.moves++
	m.openTile(m.cursorPosition)
	
	return m, nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.gameEngine.IsFinished() {
		switch msg.(tea.KeyMsg).String() {
		case "q", "ctrl+c":
			os.Exit(0)
		}
		return m, tea.Quit
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		msgString := msg.String()

		switch msgString {
		case "ctrl+c", "q":
			os.Exit(0)
		}
	
		action, ok := keyactionmap.KeyActionMap[msgString]

		if !ok {
			return m, nil
		}
		
		switch action {
		case actions.FlagTile:
			return m.FlagTile()
		case actions.MoveCursorDown:
			return m.MoveCursorDown()
		case actions.MoveCursorLeft:
			return m.MoveCursorLeft()
		case actions.MoveCursorRight:
			return m.MoveCursorRight()
		case actions.MoveCursorUp:
			return m.MoveCursorUp()
		case actions.OpenTile:
			return m.OpenTile()
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

	if m.gameEngine.IsFinished() {
		endscreen := endscreen.CreateModel(m.startTime, m.gameEngine)
		return endscreen.View()
	}

	var s strings.Builder
	s.WriteString(styles.HeaderStyle.Render(fmt.Sprintf("%v %v/%v", misc.APP_NAME, m.flags, m.config.Mines)))

	var lines strings.Builder
	for x, row := range m.tiles {
		var line string
		for y, tile := range row {
			isFocused := uint16(x) == m.cursorPosition.X && uint16(y) == m.cursorPosition.Y
			renderedTile := tilerenderer.RenderTileByContent(tile, isFocused)
			line += renderedTile
		}
		lines.WriteString("\n")
		if x == 0 {
			lines.WriteString(styles.BorderTop.Render(line))
		} else if x == len(m.tiles)-1 {
			lines.WriteString(styles.BorderBottom.Render(line))
		} else {
			lines.WriteString(line)
		}
	}

	s.WriteString(lines.String())
	s.WriteRune('\n')
	timeSinceStart := time.Since(m.startTime)

	beautifiedTime := beautifyTimeDuration(timeSinceStart)
	tea.SetWindowTitle(fmt.Sprintf("%v - %v", misc.APP_NAME, beautifiedTime))
	fmt.Fprintf(&s, "time - %v", beautifiedTime)

	return styles.TableStyle.Render(s.String())

}
