package styles

import (
	tilecontent "sweep/shared/consts/tile-content"

	lipgloss "github.com/charmbracelet/lipgloss"
)

func CreateTileStyle(color string) lipgloss.Style {
	lpColor := lipgloss.Color(color)
	if IsFillSet {
		return tileStyle.Background(lpColor).Foreground(reverseAdaptiveColor)
	}
	return tileStyle.Foreground(lpColor).Background(adaptiveColor)
}

var (
	tileStyle   = noStyle.Bold(true).Foreground(adaptiveColor).Background(lipgloss.NoColor{})
	noStyle     = lipgloss.NewStyle()
	HeaderStyle = noStyle.Bold(true)
	TableStyle  = noStyle.BorderStyle(lipgloss.RoundedBorder())

	adaptiveColor = lipgloss.AdaptiveColor{
		Dark:  "FF",
		Light: "00",
	}
	reverseAdaptiveColor = lipgloss.AdaptiveColor{
		Dark:  "00",
		Light: "FF",
	}

	IsFillSet        = false
	isCursorStyleSet = false

	DimText    = zeroStyle
	BrightText = sevenStyle

	BorderTop    = noStyle.Border(lipgloss.RoundedBorder(), true, false, false, false)
	BorderBottom = noStyle.Border(lipgloss.RoundedBorder(), false, false, true, false)

	zeroColor      string = "8"
	oneColor       string = "12"
	twoColor       string = "10"
	threeColor     string = "3"
	fourColor      string = "9"
	fiveColor      string = "13"
	sixColor       string = "5"
	sevenColor     string = "1"
	eightColor     string = "14"
	flagColor      string = "15"
	wrongFlagColor string = "15"
	mineColor      string = "15"
	emptyColor     string = "15"
	cursorColor    string

	zeroStyle      TileStyle = CreateTileStyle(zeroColor)
	oneStyle       TileStyle = CreateTileStyle(oneColor)
	twoStyle       TileStyle = CreateTileStyle(twoColor)
	threeStyle     TileStyle = CreateTileStyle(threeColor)
	fourStyle      TileStyle = CreateTileStyle(fourColor)
	fiveStyle      TileStyle = CreateTileStyle(fiveColor)
	sixStyle       TileStyle = CreateTileStyle(sixColor)
	sevenStyle     TileStyle = CreateTileStyle(sevenColor)
	eightStyle     TileStyle = CreateTileStyle(eightColor)
	flagStyle      TileStyle = CreateTileStyle(flagColor)
	wrongFlagStyle TileStyle = CreateTileStyle(wrongFlagColor)
	mineStyle      TileStyle = CreateTileStyle(mineColor)
	emptyStyle     TileStyle = CreateTileStyle(emptyColor)
	cursorStyle    TileStyle = tileStyle
)

func SetFill(fill bool) {
	if fill {
		tileStyle = tileStyle.Background(adaptiveColor).Foreground(lipgloss.NoColor{})
	}
	IsFillSet = fill
}

func SetCursorColor(color string) {
	isCursorStyleSet = true
	cursorColor = color
	cursorStyle = tileStyle.Foreground(lipgloss.Color(cursorColor))
}

func RenderCursor(tileStyle *TileStyle, cursor string) string {
	if !isCursorStyleSet {
		return tileStyle.Render(cursor)
	}
	return cursorStyle.
		Background(tileStyle.GetBackground()).
		Render(cursor)
}

func GetTileStyle(tileContent tilecontent.TileContent) *TileStyle {
	return tileStyles[tileContent]
}

func SetTileStyle(tileContent tilecontent.TileContent, tileStyle *TileStyle) {
	tileStyles[tileContent] = tileStyle
}

func SetTileColor(tileContent tilecontent.TileContent, color string) {
	newStyle := CreateTileStyle(color)
	SetTileStyle(tileContent, &newStyle)
}

type TileStyle = lipgloss.Style

var tileStyles = map[tilecontent.TileContent]*lipgloss.Style{
	tilecontent.Zero:      &zeroStyle,
	tilecontent.One:       &oneStyle,
	tilecontent.Two:       &twoStyle,
	tilecontent.Three:     &threeStyle,
	tilecontent.Four:      &fourStyle,
	tilecontent.Five:      &fiveStyle,
	tilecontent.Six:       &sixStyle,
	tilecontent.Seven:     &sevenStyle,
	tilecontent.Eight:     &eightStyle,
	tilecontent.Flag:      &flagStyle,
	tilecontent.WrongFlag: &wrongFlagStyle,
	tilecontent.Mine:      &mineStyle,
	tilecontent.Empty:     &emptyStyle,
}
