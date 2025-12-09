package gameengine

import (
	"math/rand"
	"slices"
	"sync"
	"sync/atomic"

	types "sweep/types"
)

var _ types.IGameEngine = (*GameEngine)(nil)

type GameEngine struct {
	types.IGameEngine
	isFinished       bool
	bombCount        uint16
	width            uint16
	height           uint16
	flaggedBombCount uint16
	flaggedCount     uint16
	field            [][]types.Tile
}

func (g *GameEngine) GetField() [][]types.Tile {
	return g.field
}
func (g *GameEngine) SetBombCount(count uint16) {
	if g.height != 0 && g.width != 0 && count >= g.width*g.height {
		panic("bomb count must be less than field size squared")
	}
	g.bombCount = count
}

func (g *GameEngine) CountNeighbouringBombs(position types.Position) byte {
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
			case types.ClosedBomb, types.FlaggedBomb, types.OpenBomb:
				atomic.AddUint32(&counter, 1)
			}
		}()
	}
	wg.Wait()
	return byte(counter)

}

// First return value is the content of the tile to show in the UI
// Second return value is whether the tile is a bomb
func (g *GameEngine) OpenTile(position types.Position) {
	switch g.GetTile(position) {
	case types.ClosedBomb, types.FlaggedBomb, types.OpenBomb, types.OutOfBounds:
		g.isFinished = true
		g.setTile(position, types.OpenBomb)
	case types.ClosedSafe, types.FlaggedSafe:
		g.setTile(position, types.OpenSafe)
	}
}

func (g *GameEngine) GetTile(position types.Position) types.Tile {
	x, y := position.GetCoords()
	if x >= g.width || y >= g.height {
		return types.OutOfBounds
	}
	return g.field[x][y]
}

func (g *GameEngine) setTile(position types.Position, tile types.Tile) {
	x, y := position.GetCoords()
	if x > g.width {
	}
	if y > g.height {
	}
	g.field[x][y] = tile
}
func (g *GameEngine) FlagToggleTile(position types.Position) {
	tile := g.GetTile(position)
	switch tile {
	case types.ClosedBomb:
		g.flaggedBombCount++
		g.flaggedCount++
		g.setTile(position, types.FlaggedBomb)
	case types.FlaggedBomb:
		g.flaggedBombCount--
		g.flaggedCount--
		g.setTile(position, types.ClosedBomb)
	case types.ClosedSafe:
		g.flaggedCount++
		g.setTile(position, types.FlaggedSafe)
	case types.FlaggedSafe:
		g.flaggedCount--
		g.setTile(position, types.ClosedSafe)
	}
	g.isFinished = g.flaggedBombCount == g.bombCount && g.flaggedCount == g.flaggedBombCount
}

// gameEngine.SetBombs() sets the bombs on the field
// Argument safeTile where no bomb can be generated
// To removed a safeTile simply set it out of bounds (less then 0 or more then fieldSize)
func (g *GameEngine) SetBombs(safeTile types.Position) {
	bombPositions := []types.Position{}
	maxX, maxY := g.width, g.height
	minValue := uint16(0)

	var wg sync.WaitGroup

	for len(bombPositions) < int(g.bombCount) {
		x, y := uint16(rand.Intn(int(maxX+minValue))), uint16(rand.Intn(int(maxY+minValue)))
		currentPosition := types.Position{X: x, Y: y}
		if currentPosition != safeTile && !slices.Contains(bombPositions, currentPosition) {
			bombPositions = append(bombPositions, currentPosition)
		}
	}

	for ix := range bombPositions {
		wg.Go(func() {
			g.setTile(bombPositions[ix], types.ClosedBomb)
		})
	}
	wg.Wait()

	g.bombCount = uint16(len(bombPositions))
}

func (g *GameEngine) IsFinished() bool {
	return g.isFinished
}

func (g *GameEngine) SetFieldSize(width uint16, height uint16) {
	if width == 0 {
		panic("field width cannot be 0")
	}
	if height == 0 {
		panic("field height cannot be 0")
	}

	if g.bombCount >= width*height {
		panic("width multiplied by height must be more than bomb count")
	}

	g.width = width
	g.height = height

	field := make([][]types.Tile, g.width)
	for x := range g.width {
		field[x] = make([]types.Tile, g.height)
		for y := range g.height {
			field[x][y] = types.ClosedSafe
		}
	}
	g.field = field
}

func (g *GameEngine) GetWidth() uint16 {
	return g.width
}

func (g *GameEngine) GetHeight() uint16 {
	return g.height
}
