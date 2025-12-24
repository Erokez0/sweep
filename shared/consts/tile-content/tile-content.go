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
	

	mineString = "mine"
	flagString = "flag"
	emptyString = "empty"
	wrongFlagString = "wrong flag"
)

func (tc TileContent) String() string {
	switch tc {
	case Zero:
		return "0"
	case One:
		return "1"
	case Two:
		return "2"
	case Three:
		return "3"
	case Four:
		return "4"
	case Five:
		return "5"
	case Six:
		return "6"
	case Seven:
		return "7"
	case Eight:
		return "8"
	case Mine:
		return glyphs.MINE
	case Flag:
		return glyphs.FLAG
	case WrongFlag:
		return glyphs.WRONG_FLAG
	case Empty:
		return glyphs.EMPTY
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