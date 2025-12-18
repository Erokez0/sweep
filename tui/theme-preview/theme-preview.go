package themepreview

import (
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

	result := ""
	for x := range allTiles {
		for y := range allTiles[x] {
			tile := allTiles[x][y]
			result += tilerenderer.RenderTileByContent(tile, false) + " "
			result += tilerenderer.RenderTileByContent(tile, true) + " "
		}
		result += "\n\r"
	}

	return result
}
