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
	MineCount        uint16
	width            uint16
	height           uint16
	flaggedMineCount uint16
	flaggedCount     uint16
	field            [][]types.Tile
}

func (g *GameEngine) GetField() [][]types.Tile {
	return g.field
}
func (g *GameEngine) SetMineCount(count uint16) {
	if g.height != 0 && g.width != 0 && count >= g.width*g.height {
		panic("Mine count must be less than field size squared")
	}
	g.MineCount = count
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
			case types.ClosedMine, types.FlaggedMine, types.OpenMine:
				atomic.AddUint32(&counter, 1)
			}
		}()
	}
	wg.Wait()
	return byte(counter)

}

// First return value is the content of the tile to show in the UI
// Second return value is whether the tile is a Mine
func (g *GameEngine) OpenTile(position types.Position) {
	switch g.GetTile(position) {
	case types.ClosedMine, types.FlaggedMine, types.OpenMine, types.OutOfBounds:
		g.isFinished = true
		g.setTile(position, types.OpenMine)
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
	case types.ClosedMine:
		g.flaggedMineCount++
		g.flaggedCount++
		g.setTile(position, types.FlaggedMine)
	case types.FlaggedMine:
		g.flaggedMineCount--
		g.flaggedCount--
		g.setTile(position, types.ClosedMine)
	case types.ClosedSafe:
		g.flaggedCount++
		g.setTile(position, types.FlaggedSafe)
	case types.FlaggedSafe:
		g.flaggedCount--
		g.setTile(position, types.ClosedSafe)
	}
	g.isFinished = g.flaggedMineCount == g.MineCount && g.flaggedCount == g.flaggedMineCount
}

// gameEngine.SetMines() sets the Mines on the field
// Argument safeTile where no Mine can be generated
// To removed a safeTile simply set it out of bounds (less then 0 or more then fieldSize)
func (g *GameEngine) SetMines(safeTile types.Position) {
	MinePositions := []types.Position{}
	maxX, maxY := g.width, g.height
	minValue := uint16(0)

	var wg sync.WaitGroup

	for len(MinePositions) < int(g.MineCount) {
		x, y := uint16(rand.Intn(int(maxX+minValue))), uint16(rand.Intn(int(maxY+minValue)))
		currentPosition := types.Position{X: x, Y: y}
		if currentPosition != safeTile && !slices.Contains(MinePositions, currentPosition) {
			MinePositions = append(MinePositions, currentPosition)
		}
	}

	for ix := range MinePositions {
		wg.Go(func() {
			g.setTile(MinePositions[ix], types.ClosedMine)
		})
	}
	wg.Wait()

	g.MineCount = uint16(len(MinePositions))
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

	if g.MineCount >= width*height {
		panic("width multiplied by height must be more than Mine count")
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
