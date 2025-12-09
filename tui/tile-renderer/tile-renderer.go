package tilerenderer

import (
	"fmt"

	glyphs "sweep/shared/vars/glyphs"
	styles "sweep/tui/styles"
	types "sweep/types"
)

func RenderTileByContent(tileContent string, isFocused bool) string {
	template := " %v "

	style, ok := styles.TileStyles[tileContent]
	if isFocused {
		template = "[%v]"
	}
	switch tileContent {
	case "0":
		tileContent = "x"
	case glyphs.BOMB:
		template = styles.Bomb.Render(template)
	case glyphs.WRONG_FLAG:
		template = styles.WrongFlag.Render(template)
	case glyphs.FLAG:
		template = styles.Flag.Render(template)
	case glyphs.EMPTY:
		template = styles.Empty.Render(template)
	}
	if ok {
		template = style.Render(template)
	}
	return fmt.Sprintf(template, tileContent)

}

func RenderTileByType(tile types.Tile, content string) string {
	switch tile {
	case types.ClosedBomb, types.OpenBomb:
		content = glyphs.BOMB
	case types.FlaggedBomb:
		content = glyphs.FLAG
	case types.FlaggedSafe:
		content = glyphs.WRONG_FLAG
	}

	return RenderTileByContent(content, false)
}
