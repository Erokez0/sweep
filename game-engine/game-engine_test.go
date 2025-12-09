package gameengine

import (
	"fmt"
	"sweep/types"
	"testing"
)

func TestFlagToggleFile(t *testing.T) {
	g := GameEngine{}
	field := [][]types.Tile{
		{
			types.ClosedBomb, types.OpenBomb,
		},
		{
			types.FlaggedSafe, types.FlaggedBomb,
		},
		{
			types.OpenSafe, types.ClosedSafe,
		},
	}

	expectedField := [][]types.Tile{
		{
			types.FlaggedBomb, types.FlaggedBomb,
		},
		{
			types.ClosedSafe, types.ClosedBomb,
		},
		{
			types.FlaggedSafe, types.FlaggedSafe,
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

func TestCountNeighbouringBombs(t *testing.T) {
	for bombCount := range 8 {
		g := GameEngine{}
		center := types.Position{
			X: uint16(1),
			Y: uint16(1),
		}
		g.SetFieldSize(uint16(3), uint16(3))

		g.SetBombCount(uint16(bombCount))
		g.SetBombs(center)
		expected := bombCount

		result := g.CountNeighbouringBombs(center)
		if fmt.Sprint(result) != fmt.Sprint(expected) {
			t.Errorf("[Assertion failed] %v != %v\nbombCount: %v\n", result, expected, g.bombCount)
			t.Error(g.field)
		}
	}
}

func TestFinishCondition(t *testing.T) {
	g := GameEngine{}
	g.field = [][]types.Tile{
		{
			types.ClosedBomb, types.OpenSafe,
		},
		{
			types.OpenSafe, types.OpenSafe,
		},
	}

	g.FlagToggleTile(types.Position{
		X: 0, Y: 0,
	})

	if g.IsFinished() != true {
		panic("flagged the last bomb - should have won")
	}
}

func TestOpenTile(t *testing.T) {
	g := GameEngine{}

	g.SetFieldSize(2, 2)
	g.SetBombCount(3)

	safePosition := types.Position{X: 0, Y: 0}
	g.SetBombs(safePosition)

	field := g.GetField()
	for x := range field {
		for y := range field[x] {
			tile := field[x][y]
			currentPosition := types.Position{X: uint16(x), Y: uint16(y)}
			var wanted types.Tile

			if currentPosition == safePosition {
				wanted = types.ClosedSafe
			} else {
				wanted = types.ClosedBomb
			}

			if wanted != tile {
				t.Errorf("tile at [%v, %v] should have been [%v] received [%v]\n%v", x, y, wanted, tile, g.GetField())
			}

			g.OpenTile(currentPosition)

			if wanted == types.ClosedSafe {
				wanted = types.OpenSafe
			} else {
				wanted = types.OpenBomb
			}

			tile = field[x][y]

			if wanted != tile {
				t.Errorf("tile at [%v, %v] should have been [%v] received [%v]", x, y, wanted, tile)
			}
		}
	}
}

func TestSetTile(t *testing.T) {
	tiles := []types.Tile{
		types.ClosedBomb, types.ClosedSafe, types.FlaggedBomb, types.OpenSafe, types.OpenBomb, types.FlaggedSafe, types.OpenSafe, types.OutOfBounds,
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

func TestAllFieldSizesAndBombCounts(t *testing.T) {
	g := GameEngine{}
	for width := range 10_000_000 {
		for height := range 10_000_000 {
			for bombs := range 10_000_000 {

				defer func() {
					r := recover()
					if width == 0 || height == 0 || bombs >= width*height && r != nil {
						return
					}
					if width == 0 && r == nil {
						t.Errorf("Width == 0 but has not panicked %v", r)
					}
					if height == 0 && r == nil {
						t.Error("height == 0 but has not panicked")
					}
					if bombs >= width*height && r == nil {
						t.Error("bomb count >= width * height but has not panicked")
					}
					if r != nil {
						t.Errorf("Was not supposed to panic\n%v", r)
					}
				}()

				g.SetFieldSize(uint16(width), uint16(height))
				g.SetBombCount(uint16(bombs))
				g.SetBombs(types.Position{})

			}
		}
	}
}

func TestSetBombs(t *testing.T) {
	for width := range 25 {
		if width < 2 {
			continue
		}
		for height := range 25 {
			if height < 2 {
				continue
			}
			for bombCount := range (height * width) - 1 {
				g := GameEngine{}
				g.SetFieldSize(uint16(width), uint16(height))
				g.SetBombCount(uint16(bombCount))
				position := types.Position{
					X: uint16(width / 2),
					Y: uint16(height / 2),
				}
				g.SetBombs(position)

				count := 0
				for x := range g.field {
					for y := range g.field[x] {
						if g.field[x][y] == types.ClosedBomb {
							count++
						}
					}
				}
				if uint16(count) != g.bombCount {
					t.Errorf("[Assertion failed] %v != %v\nreal bomb count != wanted bomb count", count, g.bombCount)
				}
				if g.GetTile(position) != types.ClosedSafe {
					t.Errorf("[Assertion failed] %v - %v is not safe!", position.X, position.Y)
				}
				if g.bombCount != uint16(bombCount) {
					t.Errorf("[Assertion failed] %v != %v\ngameEngine.bombCount != set bombCount", g.bombCount, bombCount)
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
