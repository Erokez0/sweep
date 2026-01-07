package tilecontent

import (
	"fmt"
	"sweep/shared/vars/glyphs"
)

type TileContent uint16

const (
	Zero TileContent = iota
	One
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight

	Mine
	Flag
	WrongFlag
	Empty

	mineString      = "mine"
	flagString      = "flag"
	emptyString     = "empty"
	wrongFlagString = "wrong flag"
)

func (tc TileContent) String() string {
	switch tc {
	case Zero:
		return glyphs.Zero
	case One:
		return glyphs.One
	case Two:
		return glyphs.Two
	case Three:
		return glyphs.Three
	case Four:
		return glyphs.Four
	case Five:
		return glyphs.Five
	case Six:
		return glyphs.Six
	case Seven:
		return glyphs.Seven
	case Eight:
		return glyphs.Eight
	case Mine:
		return glyphs.Mine
	case Flag:
		return glyphs.Flag
	case WrongFlag:
		return glyphs.WrongFlag
	case Empty:
		return glyphs.Empty
	default:
		panic("unknown tile content key")
	}
}

func FromString(str string) (TileContent, error) {
	switch str {
	case Zero.String():
		return Zero, nil
	case One.String():
		return One, nil
	case Two.String():
		return Two, nil
	case Three.String():
		return Three, nil
	case Four.String():
		return Four, nil
	case Five.String():
		return Five, nil
	case Six.String():
		return Six, nil
	case Seven.String():
		return Seven, nil
	case Eight.String():
		return Eight, nil
	case mineString, Mine.String():
		return Mine, nil
	case flagString, Flag.String():
		return Flag, nil
	case wrongFlagString, WrongFlag.String():
		return WrongFlag, nil
	case emptyString, Empty.String():
		return Empty, nil
	default:
		return *new(TileContent), fmt.Errorf("\"%v\" is not a valid tile content", str)
	}
}

func FromNumber(n byte) (TileContent, error) {
	switch n {
	case 0:
		return Zero, nil
	case 1:
		return One, nil
	case 2:
		return Two, nil
	case 3:
		return Three, nil
	case 4:
		return Four, nil
	case 5:
		return Five, nil
	case 6:
		return Six, nil
	case 7:
		return Seven, nil
	case 8:
		return Eight, nil
	default:
		return *new(TileContent), fmt.Errorf("\"%v\" is not a valid tile content number", n)
	}
}

