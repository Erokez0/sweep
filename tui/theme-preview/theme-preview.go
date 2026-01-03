package themepreview

import (
	"strings"

	tilecontent "sweep/shared/consts/tile-content"
	tilerenderer "sweep/tui/tile-renderer"
)

func RenderThemePreview() string {
	allTiles := [][]tilecontent.TileContent{
		{
			tilecontent.Zero, tilecontent.One,
		},
		{
			tilecontent.Two, tilecontent.Three,
		},
		{
			tilecontent.Four, tilecontent.Five,
		},
		{
			tilecontent.Six, tilecontent.Seven,
		},
		{
			tilecontent.Eight, tilecontent.Flag,
		},
		{
			tilecontent.WrongFlag, tilecontent.Mine,
		},
	}

	var result strings.Builder
	for x := range allTiles {
		for y := range allTiles[x] {
			tile := allTiles[x][y]
			result.WriteString(tilerenderer.RenderTileByContent(tile, false) + " ")
			result.WriteString(tilerenderer.RenderTileByContent(tile, true) + " ")
		}
		result.WriteString("\n\r")
	}

	return result.String()
}
