package cursor

import (
	colors "sweep/config/colors"
	"sweep/shared/vars/glyphs"
	"sweep/tui/styles"
)

type Cursor struct {
	Color     colors.Color `json:"color"`
	LeftHalf  string       `json:"left half"`
	RightHalf string       `json:"right half"`
}

func (c *Cursor) Validate() (bool, []string) {
	errors := []string{}
	if !c.Color.IsValid() {
		errors = append(errors, "(cursor.color) cursor color does not match ANSI nor HEX RGB")
	}
	if len(c.LeftHalf) > 1 {
		errors = append(errors, "(cursor.left half) cursor left half is longer than one character")
	}
	if len(c.RightHalf) > 1 {
		errors = append(errors, "(cursor.right half) cursor right half is longer than one character")
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