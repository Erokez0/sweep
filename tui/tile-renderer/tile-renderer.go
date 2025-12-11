package tilerenderer

import (
	"fmt"

	tiles "sweep/shared/consts/tiles"
	types "sweep/shared/types"
	glyphs "sweep/shared/vars/glyphs"
	styles "sweep/tui/styles"
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
	}
	if ok {
		template = style.Render(template)
	}
	return fmt.Sprintf(template, tileContent)

}

func RenderTileByType(tile types.Tile, content string) string {
	switch tile {
	case tiles.ClosedMine, tiles.OpenMine:
		content = glyphs.MINE
	case tiles.FlaggedMine:
		content = glyphs.FLAG
	case tiles.FlaggedSafe:
		content = glyphs.WRONG_FLAG
	}

	return RenderTileByContent(content, false)
}
