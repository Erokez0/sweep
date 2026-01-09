package glyphs

import (
	"fmt"
	tilecontent "sweep/shared/consts/tile-content"
)

type Glyphs map[string]string

type InvalidGlyphOption struct {
	tileName string
}

func (e *InvalidGlyphOption) Error() string {
	return fmt.Sprintf("(glyphs) %v is not a valid option", e.tileName)
}

func (e *InvalidGlyphOption) Is(target error) bool {
	return e.Error() == target.Error()
}

func (g Glyphs) Validate() (bool, []error) {
	errors := make([]error, 0)
	for tileName := range g {
		if _, err := tilecontent.FromString(tileName); err != nil {
			errors = append(errors, &InvalidGlyphOption{tileName})
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
