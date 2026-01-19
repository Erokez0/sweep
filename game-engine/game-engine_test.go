package gameengine

import (
	"errors"
	"fmt"
	"sync"
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
			t.Errorf("[Assertion failed] %v != %v\nmines: %v\n", result, expected, g.mines)
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

func TestSetAllFieldParameters(t *testing.T) {
	var wg sync.WaitGroup
	for width := range 25 {
		for height := range 25 {
			for mines := range 25 {
				g := GameEngine{}
				var expected error
				if width == 0 {
					expected = &FieldParameterCannotBe0Error{"field width"}
				} else if height == 0 {
					expected = &FieldParameterCannotBe0Error{"field height"}
				}

				actual := g.SetFieldSize(uint16(width), uint16(height))
				if !errors.Is(actual, expected) {
					t.Errorf("[Assertion failed]\nExpected errors: %v\nActual error: %v\nmines - %v, height - %v, width - %v", expected, actual, mines, height, width)
					continue
				}
				if width == 0 || height == 0 {
					continue
				}

				if mines == 0 {
					expected = &FieldParameterCannotBe0Error{"mine count"}
				} else if width != 0 && height != 0 && mines >= height*width {
					expected = &TooManyMinesError{}
				}

				actual = g.SetMineCount(uint16(mines))

				if !errors.Is(actual, expected) {
					t.Errorf("[Assertion failed]\nExpected errors: %v\nActual error: %v\nmines - %v, height - %v, width - %v", expected, actual, mines, height, width)
					continue
				}
			}
		}
	}
	wg.Wait()
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
			for mineCount := range (height * width) - 1 {
				g := GameEngine{}
				g.SetFieldSize(uint16(width), uint16(height))
				g.SetMineCount(uint16(mineCount))
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
				if uint16(count) != g.mines {
					t.Errorf("[Assertion failed] %v != %v\nreal mine count != wanted mine count", count, g.mines)
				}
				if g.GetTile(position) != tiles.ClosedSafe {
					t.Errorf("[Assertion failed] %v - %v is not safe!", position.X, position.Y)
				}
				if g.mines != uint16(mineCount) {
					t.Errorf("[Assertion failed] %v != %v\ngameEngine.mines != set mine count", g.mines, mineCount)
				}
			}
		}
	}
}

func TestGetWidth(t *testing.T) {
	const height uint16 = 2
	for width := range uint16(255) {
		if width == 0 {
			continue
		}
		g := GameEngine{}

		expected := width

		err := g.SetFieldSize(width, expected)
		if err != nil {
			t.Fatal(err)
		}

		actual := g.GetWidth()
		if expected != actual {
			t.Errorf("[Assertion failed] %v != %v\ngameEngine.GetWidth() != width", actual, expected)
		}
	}
}
func TestGetHeight(t *testing.T) {
	const width uint16 = 2
	for height := range uint16(255) {
		if height == 0 {
			continue
		}

		g := GameEngine{}
		expected := height
		err := g.SetFieldSize(width, expected)
		if err != nil {
			t.Fatal(err)
		}

		actual := g.GetHeight()
		if expected != actual {
			t.Errorf("[Assertion failed] %v != %v\ngameEngine.GetHeight() != height", actual, expected)
		}
		if expected != g.height {
			t.Errorf("[Assertion failed] %v != %v\ngameEngine.height != height", g.height, expected)
		}
	}
}

func Test_areAllSafeTilesOpen(t *testing.T) {
	type TestCase struct {
		expected bool
		prepare  func(*GameEngine)
	}

	testCases := []TestCase{
		{
			expected: true,
			prepare: func(g *GameEngine) {
				g.field = [][]types.Tile{
					{
						tiles.ClosedSafe, tiles.ClosedMine,
					},
					{
						tiles.ClosedSafe, tiles.ClosedSafe,
					},
				}
				g.width = 2
				g.height = 2
				g.mines = 1

				g.OpenTile(types.Position{X: 0, Y: 0})
				g.OpenTile(types.Position{X: 0, Y: 1})
				g.OpenTile(types.Position{X: 1, Y: 1})
			},
		},
		{
			expected: false,
			prepare: func(g *GameEngine) {
				g.field = [][]types.Tile{
					{
						tiles.ClosedSafe, tiles.ClosedMine,
					},
					{
						tiles.ClosedSafe, tiles.ClosedSafe,
					},
				}
				g.width = 2
				g.height = 2
				g.mines = 1

				g.OpenTile(types.Position{X: 0, Y: 0})
				g.OpenTile(types.Position{X: 0, Y: 1})
			},
		},
		{
			expected: true,
			prepare: func(g *GameEngine) {
				g.field = [][]types.Tile{
					{
						tiles.ClosedSafe,
						tiles.ClosedSafe,
						tiles.ClosedSafe,
					},

					{
						tiles.ClosedSafe,
						tiles.ClosedMine,
						tiles.ClosedSafe,
					},

					{
						tiles.ClosedSafe,
						tiles.ClosedSafe,
						tiles.ClosedSafe,
					},
				}
				g.mines = 1
				g.width = 3
				g.height = 3

				g.OpenTile(types.Position{X: 0, Y: 0})
				g.OpenTile(types.Position{X: 0, Y: 1})
				g.OpenTile(types.Position{X: 0, Y: 2})

				g.OpenTile(types.Position{X: 1, Y: 0})
				g.FlagToggleTile(types.Position{X: 1, Y: 1})
				g.OpenTile(types.Position{X: 1, Y: 2})

				g.OpenTile(types.Position{X: 2, Y: 0})
				g.OpenTile(types.Position{X: 2, Y: 1})
				g.OpenTile(types.Position{X: 2, Y: 2})
			},
		},

		{
			expected: true,
			prepare: func(g *GameEngine) {
				g.field = [][]types.Tile{
					{
						tiles.ClosedSafe, tiles.ClosedMine,
					},
					{
						tiles.ClosedSafe, tiles.ClosedSafe,
					},
				}
				g.width = 2
				g.height = 2
				g.mines = 1

				g.OpenTile(types.Position{X: 0, Y: 0})
				g.OpenTile(types.Position{X: 1, Y: 1})
				g.OpenTile(types.Position{X: 0, Y: 1})
			},
		},
	}

	for _, testCase := range testCases {
		g := new(GameEngine)
		testCase.prepare(g)

		expected := testCase.expected
		actual := g.areAllSafeTilesOpen()

		if actual != expected {
			stats := fmt.Sprintf("open count: %v\nflagged mine count: %v\nflagged count: %v\ntile count: %v\n", g.openCount, g.flaggedMineCount, g.flaggedCount, g.width*g.height)
			t.Errorf("[Assertion failed] g.areAllSafeTilesOpen\nExpected: %v\nActual: %v\n\n%v", expected, actual, stats)
		}
	}
}

func TestWinCondition(t *testing.T) {
	type TestCase struct {
		isFinished bool
		win        bool
		prepare    func(*GameEngine)
	}

	testCases := []TestCase{
		{
			isFinished: true,
			prepare: func(g *GameEngine) {
				g.field = [][]types.Tile{
					{
						tiles.ClosedSafe,
						tiles.ClosedSafe,
						tiles.ClosedSafe,
					},

					{
						tiles.ClosedSafe,
						tiles.ClosedMine,
						tiles.ClosedSafe,
					},

					{
						tiles.ClosedSafe,
						tiles.ClosedSafe,
						tiles.ClosedSafe,
					},
				}
				g.mines = 1
				g.width = 3
				g.height = 3

				g.OpenTile(types.Position{X: 0, Y: 0})
				g.OpenTile(types.Position{X: 0, Y: 1})
				g.OpenTile(types.Position{X: 0, Y: 2})

				g.OpenTile(types.Position{X: 1, Y: 0})
				g.FlagToggleTile(types.Position{X: 1, Y: 1})
				g.OpenTile(types.Position{X: 1, Y: 2})

				g.OpenTile(types.Position{X: 2, Y: 0})
				g.OpenTile(types.Position{X: 2, Y: 1})
				g.OpenTile(types.Position{X: 2, Y: 2})
			},
		},

		{
			isFinished: false,
			prepare: func(g *GameEngine) {
				g.field = [][]types.Tile{
					{
						tiles.ClosedSafe,
						tiles.ClosedSafe,
						tiles.ClosedSafe,
					},

					{
						tiles.ClosedSafe,
						tiles.ClosedMine,
						tiles.ClosedSafe,
					},

					{
						tiles.ClosedSafe,
						tiles.ClosedSafe,
						tiles.ClosedSafe,
					},
				}
				g.mines = 1
				g.width = 3
				g.height = 3

				g.OpenTile(types.Position{X: 0, Y: 0})
				g.OpenTile(types.Position{X: 0, Y: 1})
				g.OpenTile(types.Position{X: 0, Y: 2})

				g.OpenTile(types.Position{X: 1, Y: 0})
				g.OpenTile(types.Position{X: 1, Y: 2})

				g.OpenTile(types.Position{X: 2, Y: 0})
				g.OpenTile(types.Position{X: 2, Y: 1})
				g.OpenTile(types.Position{X: 2, Y: 2})
			},
		},
	}

	for n, testCase := range testCases {
		g := new(GameEngine)
		testCase.prepare(g)

		expected := testCase.isFinished
		actual := g.isFinished

		if expected != actual {
			t.Errorf("[Assertion failed] #%v g.isFinished\nExpected: %v\nActual: %v", n+1, expected, actual)
		}
	}
}
