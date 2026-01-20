package cursor

import (
	"fmt"
	colors "sweep/config/colors"
	glyphs "sweep/shared/vars/glyphs"
	styles "sweep/tui/styles"
)

const configModule string = "cursor"
const option string = "color"

const leftHalf = "left"
const rightHalf = "right"

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

func (e *CursorHalfTooLongError) Is(target error) bool {
	return e.Error() == target.Error()
}

func (c *Cursor) Validate() (bool, []error) {
	errors := []error{}
	if !c.Color.IsValid() {
		errors = append(errors, &colors.InvalidColorError{ConfigModule: configModule, Option: option, Value: c.Color})
	}
	if len(c.LeftHalf) > 1 {
		errors = append(errors, &CursorHalfTooLongError{leftHalf})
	}
	if len(c.RightHalf) > 1 {
		errors = append(errors, &CursorHalfTooLongError{rightHalf})
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
