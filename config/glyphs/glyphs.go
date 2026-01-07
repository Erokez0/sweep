package glyphs

import (
	"fmt"
	tilecontent "sweep/shared/consts/tile-content"
)

type Glyphs map[string]string

func (g Glyphs) Validate() (bool, []string) {
	errors := make([]string, 0)
	for tileName := range g {
		if _, err := tilecontent.FromString(tileName); err != nil {
			errors = append(errors, fmt.Sprintf("(glyphs) %v is not a valid option of tile content", tileName))
		}
	}

	return len(errors) == 1, errors
}

func (g Glyphs) Apply() {
	for tileName, glyph := range g {
		tileContent, _ := tilecontent.FromString(tileName)
		tilecontent.SetGlyph(tileContent, glyph)
	}
}
