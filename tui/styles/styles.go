package styles

import (
	tilecontent "sweep/shared/consts/tile-content"

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
	isCursorStyleSet = false

	DimText    = zeroStyle
	BrightText = sevenStyle

	BorderTop    = NoStyle.Border(lipgloss.RoundedBorder(), true, false, false, false)
	BorderBottom = NoStyle.Border(lipgloss.RoundedBorder(), false, false, true, false)

	ZeroColor      string = "8"
	OneColor       string = "12"
	TwoColor       string = "10"
	ThreeColor     string = "3"
	FourColor      string = "9"
	FiveColor      string = "13"
	SixColor       string = "5"
	SevenColor     string = "1"
	EightColor     string = "14"
	FlagColor      string = "15"
	WrongFlagColor string = "15"
	MineColor      string = "15"
	EmptyColor     string = "15"
	CursorColor    string

	zeroStyle      TileStyle = tileStyle.Foreground(lipgloss.Color(ZeroColor))
	oneStyle       TileStyle = tileStyle.Foreground(lipgloss.Color(OneColor))
	twoStyle       TileStyle = tileStyle.Foreground(lipgloss.Color(TwoColor))
	threeStyle     TileStyle = tileStyle.Foreground(lipgloss.Color(ThreeColor))
	fourStyle      TileStyle = tileStyle.Foreground(lipgloss.Color(FourColor))
	fiveStyle      TileStyle = tileStyle.Foreground(lipgloss.Color(FiveColor))
	sixStyle       TileStyle = tileStyle.Foreground(lipgloss.Color(SixColor))
	sevenStyle     TileStyle = tileStyle.Foreground(lipgloss.Color(SevenColor))
	eightStyle     TileStyle = tileStyle.Foreground(lipgloss.Color(EightColor))
	flagStyle      TileStyle = tileStyle.Foreground(lipgloss.Color(EightColor))
	wrongFlagStyle TileStyle = tileStyle.Foreground(lipgloss.Color(EightColor))
	mineStyle      TileStyle = tileStyle.Foreground(lipgloss.Color(EightColor))
	emptyStyle     TileStyle = tileStyle.Foreground(lipgloss.Color(EightColor))
	cursorStyle    TileStyle = tileStyle
)

func SetFill(fill bool) {
	if fill {
		tileStyle = tileStyle.Background(adaptiveColor).Foreground(lipgloss.NoColor{})
	}
	isFill = fill
}

func SetCursorColor(color string) {
	isCursorStyleSet = true
	CursorColor = color
	cursorStyle = tileStyle.Foreground(lipgloss.Color(CursorColor))
}

func RenderCursor(tileStyle *TileStyle, cursor string) string {
	if !isCursorStyleSet {
		return tileStyle.Render(cursor)
	}
	return cursorStyle.
		Background(tileStyle.GetBackground()).
		Render(cursor)
}

func SetTileColor(key tilecontent.TileContent, color string) {
	newStyle := createTileStyle(color)
	TileStyles[key] = &newStyle
}

type TileStyle = lipgloss.Style

var TileStyles = map[tilecontent.TileContent]*lipgloss.Style{
	tilecontent.Zero:  &zeroStyle,
	tilecontent.One:   &oneStyle,
	tilecontent.Two:   &twoStyle,
	tilecontent.Three: &threeStyle,
	tilecontent.Four:  &fourStyle,
	tilecontent.Five:  &fiveStyle,
	tilecontent.Six:   &sixStyle,
	tilecontent.Seven: &sevenStyle,
	tilecontent.Eight: &eightStyle,

	tilecontent.Flag:      &flagStyle,
	tilecontent.WrongFlag: &wrongFlagStyle,
	tilecontent.Mine:      &mineStyle,
	tilecontent.Empty:     &emptyStyle,
}
