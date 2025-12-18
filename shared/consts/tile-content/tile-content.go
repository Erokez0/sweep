package tilecontent

import "sweep/shared/vars/glyphs"

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

func FromNumber(n byte) TileContent {
	switch n {
	case 0:
		return Zero
	case 1:
		return One
	case 2:
		return Two
	case 3:
		return Three
	case 4:
		return Four
	case 5:
		return Five
	case 6:
		return Six
	case 7:
		return Seven
	case 8:
		return Eight
	default: 
	panic("byte not in byte range")
	}
}