package tilecontent

import (
	"fmt"

	glyphs "sweep/shared/vars/glyphs"
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

	zeroString  = "0"
	oneString   = "1"
	twoString   = "2"
	threeString = "3"
	fourString  = "4"
	fiveString  = "5"
	sixString   = "6"
	sevenString = "7"
	eightString = "8"

	mineString      = "mine"
	flagString      = "flag"
	emptyString     = "empty"
	wrongFlagString = "wrong flag"
)

type InvalidTileContentOptionError struct {
	option string
}

func (e *InvalidTileContentOptionError) Error() string {
	return fmt.Sprintf("\n%v\n is not a valid tile content option", e.option)
}
func (e *InvalidTileContentOptionError) Is(target error) bool {
	return e.Error() == target.Error()
}

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
	case zeroString, Zero.String():
		return Zero, nil
	case oneString, One.String():
		return One, nil
	case twoString, Two.String():
		return Two, nil
	case threeString, Three.String():
		return Three, nil
	case fourString, Four.String():
		return Four, nil
	case fiveString, Five.String():
		return Five, nil
	case sixString, Six.String():
		return Six, nil
	case sevenString, Seven.String():
		return Seven, nil
	case eightString, Eight.String():
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
		return *new(TileContent), &InvalidTileContentOptionError{str}
	}
}

type InvalidTileContentByteOptionError struct {
	value any
}

func (e *InvalidTileContentByteOptionError) Error() string {
	return fmt.Sprintf("\"%v\" is not a valid option for byte conversion", e.value)
}
func (e *InvalidTileContentByteOptionError) Is(target error) bool {
	return e.Error() == target.Error()
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
		return *new(TileContent), &InvalidTileContentByteOptionError{n}
	}
}

func SetGlyph(tileContent TileContent, glyph string) {
	switch tileContent {
	case Zero:
		glyphs.Zero = glyph
	case One:
		glyphs.One = glyph
	case Two:
		glyphs.Two = glyph
	case Three:
		glyphs.Three = glyph
	case Four:
		glyphs.Four = glyph
	case Five:
		glyphs.Five = glyph
	case Six:
		glyphs.Six = glyph
	case Seven:
		glyphs.Seven = glyph
	case Eight:
		glyphs.Eight = glyph
	case Mine:
		glyphs.Mine = glyph
	case Flag:
		glyphs.Flag = glyph
	case WrongFlag:
		glyphs.WrongFlag = glyph
	case Empty:
		glyphs.Empty = glyph
	}
}
func (tileContent TileContent) ToNumber() (byte, error) {
	switch tileContent {
	case Zero:
		return 0, nil
	case One:
		return 1, nil
	case Two:
		return 2, nil
	case Three:
		return 3, nil
	case Four:
		return 4, nil
	case Five:
		return 5, nil
	case Six:
		return 6, nil
	case Seven:
		return 7, nil
	case Eight:
		return 8, nil
	default:
		return *new(byte), &InvalidTileContentByteOptionError{tileContent}
	}
}
