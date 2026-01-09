package cursor

import (
	"fmt"
	colors "sweep/config/colors"
	glyphs "sweep/shared/vars/glyphs"
	styles "sweep/tui/styles"
)

type Cursor struct {
	Color     colors.Color `json:"color"`
	LeftHalf  string       `json:"left half"`
	RightHalf string       `json:"right half"`
}

type CursorHalfTooLongError struct {
	whichHalf string
}

func (e *CursorHalfTooLongError) Error() string {
	return fmt.Sprintf("(cursor.%v half) cursor %v half must be one character", e.whichHalf, e.whichHalf)
}

func (c *Cursor) Validate() (bool, []error) {
	errors := []error{}
	if !c.Color.IsValid() {
		errors = append(errors, &colors.InvalidColorError{ConfigModule: "config", Option: "color", Value: c.Color})
	}
	if len(c.LeftHalf) > 1 {
		errors = append(errors, &CursorHalfTooLongError{"left"})
	}
	if len(c.RightHalf) > 1 {
		errors = append(errors, &CursorHalfTooLongError{"right"})
	}

	return len(errors) == 0, errors
}

func (c *Cursor) Apply() {
	if c.Color.IsSet() {
		styles.SetCursorColor(string(c.Color))
	}
	if c.LeftHalf != "" {
		glyphs.CursorLeftHalf = c.LeftHalf
	}
	if c.RightHalf != "" {
		glyphs.CursorRightHalf = c.RightHalf
	}
}
