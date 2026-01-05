package gameengine

import (
	"fmt"
	"testing"

	tiles "sweep/shared/consts/tiles"
	types "sweep/shared/types"
)

func TestFlagToggleFile(t *testing.T) {
	g := GameEngine{}
	field := [][]types.Tile{
		{
			tiles.ClosedMine, tiles.OpenMine,
		},
		{
			tiles.FlaggedSafe, tiles.FlaggedMine,
		},
		{
			tiles.OpenSafe, tiles.ClosedSafe,
		},
	}

	expectedField := [][]types.Tile{
		{
			tiles.FlaggedMine, tiles.FlaggedMine,
		},
		{
			tiles.ClosedSafe, tiles.ClosedMine,
		},
		{
			tiles.FlaggedSafe, tiles.FlaggedSafe,
		},
	}

	for x := range g.field {
		for y := range g.field[x] {

			position := types.Position{
				X: uint16(x),
				Y: uint16(y),
			}

			g.setTile(position, field[x][y])

			g.FlagToggleTile(position)

			expected := expectedField[x][y]
			result := g.GetTile(position)
			if result != expected {
				t.Errorf("[Assertion failed] %v != %v [x: %v, y: %v]", result, expected, x, y)
			}
		}
	}
}

func TestCountNeighbouringMines(t *testing.T) {
	for MineCount := range 8 {
		g := GameEngine{}
		center := types.Position{
			X: uint16(1),
			Y: uint16(1),
		}
		g.SetFieldSize(uint16(3), uint16(3))

		g.SetMineCount(uint16(MineCount))
		g.SetMines(center)
		expected := MineCount

		result := g.CountNeighbouringMines(center)
		if fmt.Sprint(result) != fmt.Sprint(expected) {
			t.Errorf("[Assertion failed] %v != %v\nMineCount: %v\n", result, expected, g.MineCount)
			t.Error(g.field)
		}
	}
}

func TestFinishCondition(t *testing.T) {
	g := GameEngine{}
	g.field = [][]types.Tile{
		{
			tiles.ClosedMine, tiles.OpenSafe,
		},
		{
			tiles.OpenSafe, tiles.OpenSafe,
		},
	}

	g.FlagToggleTile(types.Position{
		X: 0, Y: 0,
	})

	if g.IsFinished() != true {
		panic("flagged the last Mine - should have won")
	}
}

func TestOpenTile(t *testing.T) {
	g := GameEngine{}

	g.SetFieldSize(2, 2)
	g.SetMineCount(3)

	safePosition := types.Position{X: 0, Y: 0}
	g.SetMines(safePosition)

	width, height := g.width, g.height

	for y := range height {
		for x := range width {
			currentPosition := types.Position{X: uint16(x), Y: uint16(y)}
			tile := g.GetTile(currentPosition)
			var wanted types.Tile

			if currentPosition == safePosition {
				wanted = tiles.ClosedSafe
			} else {
				wanted = tiles.ClosedMine
			}

			if wanted != tile {
				t.Errorf("tile at [%v, %v] should have been [%v] received [%v]\n%v", x, y, wanted, tile, g.GetField())
			}

			g.OpenTile(currentPosition)

			if wanted == tiles.ClosedSafe {
				wanted = tiles.OpenSafe
			} else {
				wanted = tiles.OpenMine
			}

			tile = g.GetTile(currentPosition)

			if wanted != tile {
				t.Errorf("tile at [%v, %v] should have been [%v] received [%v]", x, y, wanted, tile)
			}
		}
	}
}

func TestSetTile(t *testing.T) {
	tiles := []types.Tile{
		tiles.ClosedMine, tiles.ClosedSafe, tiles.FlaggedMine, tiles.OpenSafe, tiles.OpenMine, tiles.FlaggedSafe, tiles.OpenSafe, tiles.OutOfBounds,
	}
	for _, tile := range tiles {
		g := GameEngine{}
		g.SetFieldSize(3, 3)
		position := types.Position{
			X: 1,
			Y: 1,
		}
		g.setTile(position, tile)

		actual := g.GetTile(position)
		expected := tile

		if actual != expected {
			t.Errorf("[Assertion failed] %v != %v\ntile != gameEngine.GetTile()", expected, actual)
		}
	}
}

func TestAllFieldSizesAndMineCounts(t *testing.T) {
	g := GameEngine{}
	for width := range 10_000_000 {
		for height := range 10_000_000 {
			for Mines := range 10_000_000 {

				defer func() {
					r := recover()
					if width == 0 || height == 0 || Mines >= width*height && r != nil {
						return
					}
					if width == 0 && r == nil {
						t.Errorf("Width == 0 but has not panicked %v", r)
					}
					if height == 0 && r == nil {
						t.Error("height == 0 but has not panicked")
					}
					if Mines >= width*height && r == nil {
						t.Error("Mine count >= width * height but has not panicked")
					}
					if r != nil {
						t.Errorf("Was not supposed to panic\n%v", r)
					}
				}()

				g.SetFieldSize(uint16(width), uint16(height))
				g.SetMineCount(uint16(Mines))
				g.SetMines(types.Position{})

			}
		}
	}
}

func TestSetMines(t *testing.T) {
	for width := range 25 {
		if width < 2 {
			continue
		}
		for height := range 25 {
			if height < 2 {
				continue
			}
			for MineCount := range (height * width) - 1 {
				g := GameEngine{}
				g.SetFieldSize(uint16(width), uint16(height))
				g.SetMineCount(uint16(MineCount))
				position := types.Position{
					X: uint16(width / 2),
					Y: uint16(height / 2),
				}
				g.SetMines(position)

				count := 0
				for x := range g.field {
					for y := range g.field[x] {
						if g.field[x][y] == tiles.ClosedMine {
							count++
						}
					}
				}
				if uint16(count) != g.MineCount {
					t.Errorf("[Assertion failed] %v != %v\nreal Mine count != wanted Mine count", count, g.MineCount)
				}
				if g.GetTile(position) != tiles.ClosedSafe {
					t.Errorf("[Assertion failed] %v - %v is not safe!", position.X, position.Y)
				}
				if g.MineCount != uint16(MineCount) {
					t.Errorf("[Assertion failed] %v != %v\ngameEngine.MineCount != set MineCount", g.MineCount, MineCount)
				}
			}
		}
	}
}

func TestGetWidth(t *testing.T) {
	const height uint16 = 2
	for width := range 10_000_000 {
		if width < 1 {
			defer func() {
				r := recover()
				if r == nil {
					t.Errorf("[Error expected] width == %v, must not be 0", width)
				}
			}()
		}
		g := GameEngine{}

		expected := uint16(width)
		g.SetFieldSize(expected, height)

		actual := g.GetWidth()
		if expected != actual {
			t.Errorf("[Assertion failed] %v != %v\ngameEngine.GetWidth() != width", actual, expected)
		}
		if expected != g.width {
			t.Errorf("[Assertion failed] %v != %v\ngameEngine.width != width", g.width, expected)
		}
	}
}
func TestGetHeight(t *testing.T) {
	const width uint16 = 2
	for height := range 10_000_000 {
		if height < 1 {
			defer func() {
				r := recover()
				if r == nil {
					t.Errorf("[Error expected] height == %v, must not be 0", height)
				}
			}()
		}
		g := GameEngine{}
		expected := uint16(height)
		g.SetFieldSize(width, expected)

		actual := g.GetHeight()
		if expected != actual {
			t.Errorf("[Assertion failed] %v != %v\ngameEngine.GetHeight() != height", actual, expected)
		}
		if expected != g.height {
			t.Errorf("[Assertion failed] %v != %v\ngameEngine.height != height", g.height, expected)
		}
	}
}
