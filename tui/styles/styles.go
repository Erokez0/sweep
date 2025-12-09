package styles

import (
	"sweep/shared/vars/glyphs"

	lipgloss "github.com/charmbracelet/lipgloss"
)

func createTileStyle(color string) lipgloss.Style {
	lpColor := lipgloss.Color(color)
	if isFill {
		return tileStyle.Background(lpColor).Foreground(reverseAdaptiveColor)
	}
	return tileStyle.Foreground(lpColor).Background(adaptiveColor)
}

var (
	tileStyle   = NoStyle.Bold(true).Foreground(adaptiveColor).Background(lipgloss.NoColor{})
	NoStyle     = lipgloss.NewStyle()
	HeaderStyle = NoStyle.Bold(true)
	TableStyle  = NoStyle.BorderStyle(lipgloss.RoundedBorder())

	adaptiveColor = lipgloss.AdaptiveColor{
		Dark:  "FF",
		Light: "00",
	}
	reverseAdaptiveColor = lipgloss.AdaptiveColor{
		Dark:  "00",
		Light: "FF",
	}

	isFill = false

	DimText    = Zero
	BrightText = Seven

	BorderTop    = NoStyle.Border(lipgloss.RoundedBorder(), true, false, false, false)
	BorderBottom = NoStyle.Border(lipgloss.RoundedBorder(), false, false, true, false)

	ZeroColor  = "8"
	OneColor   = "12"
	TwoColor   = "10"
	ThreeColor = "3"
	FourColor  = "9"
	FiveColor  = "13"
	SixColor   = "5"
	SevenColor = "1"
	EightColor = "14"
	FlagColor = "15"
	WrongFlagColor = "15"
	MineColor = "15"
	EmptyColor = "15"

	Zero  TileStyle = tileStyle.Foreground(lipgloss.Color(ZeroColor))
	One   TileStyle = tileStyle.Foreground(lipgloss.Color(OneColor))
	Two   TileStyle = tileStyle.Foreground(lipgloss.Color(TwoColor))
	Three TileStyle = tileStyle.Foreground(lipgloss.Color(ThreeColor))
	Four  TileStyle = tileStyle.Foreground(lipgloss.Color(FourColor))
	Five  TileStyle = tileStyle.Foreground(lipgloss.Color(FiveColor))
	Six   TileStyle = tileStyle.Foreground(lipgloss.Color(SixColor))
	Seven TileStyle = tileStyle.Foreground(lipgloss.Color(SevenColor))
	Eight TileStyle = tileStyle.Foreground(lipgloss.Color(EightColor))
	Flag TileStyle = tileStyle.Foreground(lipgloss.Color(EightColor))
	WrongFlag TileStyle = tileStyle.Foreground(lipgloss.Color(EightColor))
	Mine TileStyle = tileStyle.Foreground(lipgloss.Color(EightColor))
	Empty TileStyle = tileStyle.Foreground(lipgloss.Color(EightColor))

)

func SetFill(fill bool) {
	if fill {
		tileStyle = tileStyle.Background(adaptiveColor).Foreground(lipgloss.NoColor{})
		Flag = tileStyle
		Mine = tileStyle
		WrongFlag = tileStyle
		Empty = tileStyle
	}
	isFill = fill
}

func SetColor(colorName, newColor string) {
	switch colorName {
	case "0":
		ZeroColor = newColor
		Zero = createTileStyle(newColor)
	case "1":
		OneColor = newColor
		One = createTileStyle(newColor)
	case "2":
		TwoColor = newColor
		Two = createTileStyle(newColor)
	case "3":
		ThreeColor = newColor
		Three = createTileStyle(newColor)
	case "4":
		FourColor = newColor
		Four = createTileStyle(newColor)
	case "5":
		FiveColor = newColor
		Five = createTileStyle(newColor)
	case "6":
		SixColor = newColor
		Six = createTileStyle(newColor)
	case "7":
		SevenColor = newColor
		Seven = createTileStyle(newColor)
	case "8":
		EightColor = newColor
		Eight = createTileStyle(newColor)
	default:
		return
	}
}

type TileStyle = lipgloss.Style

var TileStyles = map[string]*lipgloss.Style{
	"0": &Zero,
	"1": &One,
	"2": &Two,
	"3": &Three,
	"4": &Four,
	"5": &Five,
	"6": &Six,
	"7": &Seven,
	"8": &Eight,

	glyphs.FLAG: &Flag,
	glyphs.WRONG_FLAG: &WrongFlag,
	glyphs.MINE: &Mine,
	glyphs.EMPTY: &Empty,
}
