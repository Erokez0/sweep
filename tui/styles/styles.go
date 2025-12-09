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
	}
	isFill = fill
}

func SetColor(key, value string) {
	switch key {
	case "0":
		ZeroColor = value
		Zero = createTileStyle(value)
	case "1":
		OneColor = value
		One = createTileStyle(value)
	case "2":
		TwoColor = value
		Two = createTileStyle(value)
	case "3":
		ThreeColor = value
		Three = createTileStyle(value)
	case "4":
		FourColor = value
		Four = createTileStyle(value)
	case "5":
		FiveColor = value
		Five = createTileStyle(value)
	case "6":
		SixColor = value
		Six = createTileStyle(value)
	case "7":
		SevenColor = value
		Seven = createTileStyle(value)
	case "8":
		EightColor = value
		Eight = createTileStyle(value)
	case "mine":
		MineColor = value
		Mine = createTileStyle(value)
	case "flag":
		FlagColor = value
		Flag = createTileStyle(value)
	case "wrong flag":
		WrongFlagColor = value
		WrongFlag = createTileStyle(value)
	case "empty":
		EmptyColor = value
		Empty = createTileStyle(value)
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
