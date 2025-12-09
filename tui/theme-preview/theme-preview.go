package themepreview

import (
	glyphs "sweep/shared/vars/glyphs"

	tilerenderer "sweep/tui/tile-renderer"
)

func RenderThemePreview() string {
	allTiles := [][]string{
		{
			"0", "1",
		},
		{
			"2", "3",
		},
		{
			"4", "5",
		},
		{
			"6", "7",
		},
		{
			glyphs.MINE, glyphs.FLAG,
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
