package gameengine

import (
	"fmt"
	"math/rand"
	"slices"
	"sync"
	"sync/atomic"

	tiles "sweep/shared/consts/tiles"
	types "sweep/shared/types"
)

var _ types.IGameEngine = (*GameEngine)(nil)

type GameEngine struct {
	isFinished       bool
	mines            uint16
	width            uint16
	height           uint16
	flaggedMineCount uint16
	flaggedCount     uint16
	openCount        uint16
	field            [][]types.Tile
}

func (g *GameEngine) GetField() [][]types.Tile {
	return g.field
}

type TooManyMinesError struct{}

func (e *TooManyMinesError) Error() string {
	return "mine count must be less than field width multiplied by field height"
}
func (e *TooManyMinesError) Is(target error) bool {
	return target.Error() == e.Error()
}

type FieldParameterCannotBe0Error struct {
	param string
}

func (e *FieldParameterCannotBe0Error) Error() string {
	return fmt.Sprintf("%v can not be 0", e.param)
}
func (e *FieldParameterCannotBe0Error) Is(target error) bool {
	return target.Error() == e.Error()
}

func (g *GameEngine) SetMineCount(count uint16) error {
	if count == 0 {
		return &FieldParameterCannotBe0Error{"mine count"}
	}
	if count >= g.width*g.height {
		return &TooManyMinesError{}
	}
	g.mines = count

	return nil
}

func (g *GameEngine) CountNeighbouringMines(position types.Position) byte {
	x, y := position.GetCoords()
	if x > g.width {
		return 0
	}
	if y > g.height {
		return 0
	}
	var counter uint32 = 0

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
	wg.Add(len(neighbours))
	for _, neighbour := range neighbours {
		go func() {
			defer wg.Done()
			tile := g.GetTile(neighbour)
			switch tile {
			case tiles.ClosedMine, tiles.FlaggedMine, tiles.OpenMine:
				atomic.AddUint32(&counter, 1)
			}
		}()
	}
	wg.Wait()
	return byte(counter)

}

func (g *GameEngine) OpenTile(position types.Position) {
	switch g.GetTile(position) {
	case tiles.ClosedMine, tiles.FlaggedMine, tiles.OpenMine:
		g.openCount++
		g.isFinished = true
		g.setTile(position, tiles.OpenMine)
		return
	case tiles.FlaggedSafe:
		g.openCount++
		g.flaggedCount--
	case tiles.ClosedSafe:
		g.setTile(position, tiles.OpenSafe)
		g.openCount++
	}
	g.checkWinCondition()
}

func (g *GameEngine) GetTile(position types.Position) types.Tile {
	x, y := position.GetCoords()
	if x >= g.width || y >= g.height {
		return tiles.OutOfBounds
	}
	return g.field[y][x]
}

func (g *GameEngine) setTile(position types.Position, tile types.Tile) {
	x, y := position.GetCoords()
	if x > g.width {
	}
	if y > g.height {
	}
	g.field[y][x] = tile
}

func (g *GameEngine) areAllMinesFlagged() bool {
	return g.flaggedCount == g.mines && g.mines == g.flaggedMineCount
}

func (g *GameEngine) areAllSafeTilesOpen() bool {
	tileCount := g.width * g.height

	return tileCount-g.mines <= g.openCount
}

func (g *GameEngine) checkWinCondition() {
	if g.isFinished {
		return
	}
	g.isFinished = (g.areAllMinesFlagged() && g.areAllSafeTilesOpen())
}

// Second return value is whether the tile is a Mine
func (g *GameEngine) FlagToggleTile(position types.Position) {
	tile := g.GetTile(position)
	switch tile {
	case tiles.ClosedMine:
		g.flaggedMineCount++
		g.flaggedCount++
		g.setTile(position, tiles.FlaggedMine)
	case tiles.FlaggedMine:
		g.flaggedMineCount--
		g.flaggedCount--
		g.setTile(position, tiles.ClosedMine)
	case tiles.ClosedSafe:
		g.flaggedCount++
		g.setTile(position, tiles.FlaggedSafe)
	case tiles.FlaggedSafe:
		g.flaggedCount--
		g.setTile(position, tiles.ClosedSafe)
	}
	g.checkWinCondition()
}

// gameEngine.SetMines() sets the Mines on the field
// Argument safeTile where no Mine can be generated
// To removed a safeTile simply set it out of bounds (less then 0 or more then fieldSize)
func (g *GameEngine) SetMines(safeTile types.Position) {
	MinePositions := []types.Position{}
	maxX, maxY := g.width, g.height
	minValue := uint16(0)

	var wg sync.WaitGroup

	for len(MinePositions) < int(g.mines) {
		x, y := uint16(rand.Intn(int(maxX+minValue))), uint16(rand.Intn(int(maxY+minValue)))
		currentPosition := types.Position{X: x, Y: y}
		if currentPosition != safeTile && !slices.Contains(MinePositions, currentPosition) {
			MinePositions = append(MinePositions, currentPosition)
		}
	}

	for ix := range MinePositions {
		wg.Go(func() {
			g.setTile(MinePositions[ix], tiles.ClosedMine)
		})
	}
	wg.Wait()

	g.mines = uint16(len(MinePositions))
}

func (g *GameEngine) IsFinished() bool {
	return g.isFinished
}

func (g *GameEngine) SetFieldSize(width uint16, height uint16) error {
	if width == 0 {
		return &FieldParameterCannotBe0Error{"field width"}
	}
	if height == 0 {
		return &FieldParameterCannotBe0Error{"field height"}
	}

	if g.mines != 0 && g.mines >= width*height {
		return &TooManyMinesError{}
	}

	g.width = width
	g.height = height

	field := make([][]types.Tile, g.height)
	for y := range g.height {
		field[y] = make([]types.Tile, g.width)
		for x := range g.width {
			field[y][x] = tiles.ClosedSafe
		}
	}
	g.field = field

	return nil
}

func (g *GameEngine) GetWidth() uint16 {
	return g.width
}

func (g *GameEngine) GetHeight() uint16 {
	return g.height
}
