package tilerenderer

import (
	"fmt"

	tilecontent "sweep/shared/consts/tile-content"
	tiles "sweep/shared/consts/tiles"
	types "sweep/shared/types"
	glyphs "sweep/shared/vars/glyphs"
	styles "sweep/tui/styles"
)

func RenderTileByContent(tileContent tilecontent.TileContent, isFocused bool) string {
	stringTileContent := tileContent.String()

	style := styles.TileStyles[tileContent]
	template := style.Render("%v%v%v")

	switch stringTileContent {
	case "0":
		stringTileContent = "x"
	}
	stringTileContent = style.Render(stringTileContent)

	leftCursosHalf := style.Render(" ")
	rightCursorHalf := leftCursosHalf
	if isFocused {
		leftCursosHalf = styles.RenderCursor(style, glyphs.CursorLeftHalf)
		rightCursorHalf = styles.RenderCursor(style, glyphs.CursorRightHalf)
	}

	return fmt.Sprintf(template, leftCursosHalf, stringTileContent, rightCursorHalf)
}

func RenderTileByType(tile types.Tile, tileContent tilecontent.TileContent) string {
	switch tile {
	case tiles.ClosedMine, tiles.OpenMine:
		tileContent = tilecontent.Mine
	case tiles.FlaggedMine:
		tileContent = tilecontent.Flag
	case tiles.FlaggedSafe:
		tileContent = tilecontent.WrongFlag
	}
	return RenderTileByContent(tileContent, false)
}
