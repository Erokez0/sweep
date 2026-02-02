package gametui

import (
	"fmt"
	"os"
	"strings"
	"time"

	config "sweep/config"
	gameengine "sweep/game-engine"
	actions "sweep/shared/consts/actions"
	misc "sweep/shared/consts/misc"
	tilecontent "sweep/shared/consts/tile-content"
	tiles "sweep/shared/consts/tiles"
	types "sweep/shared/types"
	utils "sweep/shared/utils"
	endscreen "sweep/tui/end-screen"
	styles "sweep/tui/styles"
	tilerenderer "sweep/tui/tile-renderer"

	tea "github.com/charmbracelet/bubbletea"
)

type Tiles [][]tilecontent.TileContent

func (t *Tiles) SetTile(position types.Position, tile tilecontent.TileContent) {
	x, y := position.GetCoords()
	(*t)[y][x] = tile
}

type TileOutOfBoundsError struct {
	position types.Position
}

func (e *TileOutOfBoundsError) Error() string {
	return fmt.Sprintf("tile is out of bounds\nx: %v,y: %v", e.position.X, e.position.Y)
}
func (e *TileOutOfBoundsError) Is(target error) bool {
	return e.Error() == target.Error()
}

func (t Tiles) GetTile(position types.Position) (tilecontent.TileContent, error) {
	x, y := position.GetCoords()
	if int(y) >= len(t) || int(x) >= len(t[y]) {
		return *new(tilecontent.TileContent), &TileOutOfBoundsError{}
	}
	return t[y][x], nil
}

func CreateTiles(width, height uint16) *Tiles {
	tiles := make(Tiles, height)
	for y := range height {
		tiles[y] = make([]tilecontent.TileContent, width)
		for x := range width {
			tiles[y][x] = tilecontent.Empty
		}
	}
	return &tiles
}

type model struct {
	screenWidth            int
	keyPressBuffer         string
	previousKeyPressBuffer string
	config                 config.Config
	cursorPosition         types.Position
	gameEngine             types.IGameEngine
	tiles                  Tiles
	startTime              time.Time
	openedATile            bool
	flags                  int16
}

func CreateModel(config *config.Config) model {
	gameEngine := gameengine.GameEngine{}
	err := gameEngine.SetFieldSize(config.Width, config.Height)
	if err != nil {
		fmt.Println(err)
	}
	err = gameEngine.SetMineCount(config.Mines)
	if err != nil {
		fmt.Println(err)
	}

	return model{
		cursorPosition: types.Position{
			X: config.Width / 2,
			Y: config.Height / 2,
		},
		gameEngine:     &gameEngine,
		tiles:          *CreateTiles(config.Width, config.Height),
		startTime:      time.Now(),
		openedATile:    false,
		config:         *config,
		keyPressBuffer: "",
	}
}

func (m model) Init() tea.Cmd {
	return tea.SetWindowTitle(misc.AppName)
}

var _ tea.Model = model{}

func (m *model) openTile(position types.Position) {
	tileType := m.gameEngine.GetTile(position)
	switch tileType {
	case tiles.OutOfBounds:
		return
	case tiles.OpenSafe:
		m.openAroundOpenTile(position)
		return
	case tiles.FlaggedSafe, tiles.FlaggedMine:
		m.flags--
	}

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

	m.tiles.SetTile(position, tileContent)

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

	for _, neighbour := range neighbours {
		switch m.gameEngine.GetTile(neighbour) {
		case tiles.ClosedSafe:
			m.openTile(neighbour)
		}
	}
}

func (m *model) openAroundOpenTile(position types.Position) {
	tileContent, err := m.tiles.GetTile(position)
	if err != nil {
		return
	}
	tileCount, _ := tileContent.ToNumber()
	var flagCount byte

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

	for _, position := range neighbours {
		if tc, err := m.tiles.GetTile(position); err == nil && tc == tilecontent.Flag {
			flagCount++
		}
	}

	if tileCount != flagCount {
		return
	}

	for _, position := range neighbours {
		switch m.gameEngine.GetTile(position) {
		case tiles.FlaggedMine, tiles.FlaggedSafe, tiles.OutOfBounds:
			continue
		}
		m.gameEngine.OpenTile(position)

		if m.gameEngine.IsFinished() {
			defer tea.Quit()
			continue
		}

		switch m.gameEngine.GetTile(position) {
		case tiles.OpenSafe:
			count := m.gameEngine.CountNeighbouringMines(position)

			tileContent, err := tilecontent.FromNumber(count)
			if err != nil {
				panic(err)
			}

			m.tiles.SetTile(position, tileContent)
			if count == 0 {
				m.openSafeAroundTile(position)
			}
		}
	}
}

func (m *model) MoveCursorUp(quantifier uint16) {
	if m.cursorPosition.Y >= uint16(m.gameEngine.GetHeight())-1 {
		return
	}

	isOutOfBounds := m.cursorPosition.Y+quantifier >= m.config.Height
	if isOutOfBounds {
		m.MoveCursorToTopRow(1)
		return
	}
	m.cursorPosition.Y += quantifier

}
func (m *model) MoveCursorDown(quantifier uint16) {
	if m.cursorPosition.Y <= 0 {
		return
	}

	isOutOfBounds := int32(m.cursorPosition.Y)-int32(quantifier) < 0
	if isOutOfBounds {
		m.MoveCursorToBottomRow(1)
		return
	}
	m.cursorPosition.Y -= quantifier

}
func (m *model) MoveCursorRight(quantifier uint16) {
	if m.cursorPosition.X >= uint16(m.gameEngine.GetWidth())-1 {
		return
	}

	isOutOfBounds := int32(m.cursorPosition.X)-int32(quantifier) >= int32(m.config.Width)
	if isOutOfBounds {
		m.MoveCursorToLastColumn(1)
		return
	}
	m.cursorPosition.X += quantifier
}

func (m *model) MoveCursorLeft(quantifier uint16) {
	if m.cursorPosition.X == 0 {
		return
	}
	isOutOfBounds := int32(m.cursorPosition.X)-int32(quantifier) < 0
	if isOutOfBounds {
		m.MoveCursorToFirstColumn(1)
		return
	}
	m.cursorPosition.X -= quantifier
}

func (m *model) FlagTile(_ uint16) {
	if !m.openedATile {
		return
	}

	m.gameEngine.FlagToggleTile(m.cursorPosition)

	tile, err := m.tiles.GetTile(m.cursorPosition)
	if err != nil {
		return
	}

	switch tile {
	case tilecontent.Empty:
		m.tiles.SetTile(m.cursorPosition, tilecontent.Flag)
		m.flags++
	case tilecontent.Flag:
		m.tiles.SetTile(m.cursorPosition, tilecontent.Empty)
		m.flags--
	}
}
func (m *model) OpenTile(_ uint16) {
	if !m.openedATile {
		m.gameEngine.SetMines(m.cursorPosition)
		m.openedATile = true
	}
	m.openTile(m.cursorPosition)
}

func (m *model) MoveCursorToTopRow(quantifier uint16) {
	if quantifier == 1 {
		m.cursorPosition.Y = m.config.Height - 1
		return
	}
	isOutOfBounds := int32(m.config.Height)-int32(quantifier) < 0
	if isOutOfBounds {
		m.cursorPosition.Y = m.config.Height - 1
		return
	}
	m.cursorPosition.Y = m.config.Height - quantifier
}

func (m *model) MoveCursorToBottomRow(quantifier uint16) {
	if quantifier == 1 {
		m.cursorPosition.Y = 0
		return
	}
	isOutOfBounds := int32(m.config.Height)-int32(quantifier) < 0
	if isOutOfBounds {
		m.cursorPosition.Y = m.config.Height - 1
		return
	}
	m.cursorPosition.Y = m.config.Height - quantifier
}

func (m *model) MoveCursorToFirstColumn(_ uint16) {
	m.cursorPosition.X = 0
}

func (m *model) MoveCursorToLastColumn(quantifier uint16) {
	m.cursorPosition.X = m.config.Width - 1
	if quantifier == 1 {
		return
	}
	isOutOfBounds := int32(m.cursorPosition.Y)-int32(quantifier) < 0
	if isOutOfBounds {
		m.MoveCursorToBottomRow(1)
		return
	}
	m.cursorPosition.Y -= quantifier
}

func (m *model) doAction(action *actions.Action) {
	quantifier := action.Quantifier

	var actionHandler func(uint16)
	switch action.Kind {
	case actions.FlagTile:
		actionHandler = m.FlagTile
	case actions.MoveCursorDown:
		actionHandler = m.MoveCursorDown
	case actions.MoveCursorLeft:
		actionHandler = m.MoveCursorLeft
	case actions.MoveCursorRight:
		actionHandler = m.MoveCursorRight
	case actions.MoveCursorUp:
		actionHandler = m.MoveCursorUp
	case actions.OpenTile:
		actionHandler = m.OpenTile
	case actions.MoveCursorToBottomRow:
		actionHandler = m.MoveCursorToBottomRow
	case actions.MoveCursorToTopRow:
		actionHandler = m.MoveCursorToTopRow
	case actions.MoveCursorToFirstColumn:
		actionHandler = m.MoveCursorToFirstColumn
	case actions.MoveCursorToLastColumn:
		actionHandler = m.MoveCursorToLastColumn
	}
	actionHandler(quantifier)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.gameEngine.IsFinished() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "q", "ctrl+c":
				os.Exit(0)
			}
		}
		return m, tea.Quit
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.screenWidth = msg.Width
	case tea.KeyMsg:
		msgString := msg.String()

		switch msgString {
		case "ctrl+c", "q":
			os.Exit(0)
		case "esc":
			m.previousKeyPressBuffer = m.keyPressBuffer
			m.keyPressBuffer = ""
			return m, nil
		}
		m.keyPressBuffer += msgString

		if !actions.AnyBindingStartWith(m.keyPressBuffer) {
			m.previousKeyPressBuffer = m.keyPressBuffer
			m.keyPressBuffer = ""
			return m, nil
		}

		action, err := actions.GetAction(m.keyPressBuffer)
		if err != nil {
			return m, nil
		}

		m.previousKeyPressBuffer = m.keyPressBuffer
		m.keyPressBuffer = ""

		m.doAction(action)
	}

	return m, nil
}

func (m model) renderHeader(s *strings.Builder) {
	header := styles.HeaderStyle.Render(fmt.Sprintf("%v %v/%v", misc.AppName, m.flags, m.config.Mines))
	s.WriteString(header)
}

func (m model) renderTiles(s *strings.Builder) {
	var lines strings.Builder
	for row := range m.config.Height {
		y := (m.config.Height - 1 - row)
		var line string
		for col := range m.config.Width {
			x := col
			isFocused := uint16(x) == m.cursorPosition.X && uint16(y) == m.cursorPosition.Y
			tile, err := m.tiles.GetTile(types.Position{X: x, Y: y})
			if err != nil {
				panic(err)
			}
			renderedTile := tilerenderer.RenderTileByContent(tile, isFocused)
			line += renderedTile
		}
		lines.WriteString("\n")
		if row == 0 {
			lines.WriteString(styles.BorderTop.Render(line))
		} else if row == uint16(len(m.tiles)-1) {
			lines.WriteString(styles.BorderBottom.Render(line))
		} else {
			lines.WriteString(line)
		}
	}

	s.WriteString(lines.String())
	s.WriteRune('\n')
}

func (m model) renderFooter(s *strings.Builder) {
	formattedDuration := utils.FormatTime(time.Since(m.startTime))
	timeStr := fmt.Sprintf("%v", formattedDuration)

	var keysStr string
	if m.keyPressBuffer != "" {
		keysStr = m.keyPressBuffer
	} else {
		keysStr = m.previousKeyPressBuffer
	}

	margin := int(m.gameEngine.GetWidth())*3 - len(timeStr)

	s.WriteString(timeStr + styles.MarginLeft(margin, keysStr))
}

func (m model) View() string {
	if m.gameEngine.IsFinished() {
		endscreen := endscreen.CreateModel(m.startTime, m.gameEngine)
		return endscreen.View()
	}

	var s strings.Builder
	m.renderHeader(&s)

	m.renderTiles(&s)

	m.renderFooter(&s)

	return styles.TableStyle.Render(s.String())
}
