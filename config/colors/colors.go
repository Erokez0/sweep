package colors

import (
	"fmt"

	tilecontent "sweep/shared/consts/tile-content"
	regexes "sweep/shared/vars/regexes"
	styles "sweep/tui/styles"
)

type Color string

func (c Color) IsSet() bool {
	return c != ""
}

func (c Color) IsValid() bool {
	return !c.IsSet() || regexes.ColorRegex.MatchString(string(c))
}

type Colors map[string]Color

type InvalidColorOptionError struct {
	option string
}

func (e *InvalidColorOptionError) Error() string {
	return fmt.Sprintf("(colors) %v is not a valid option", e.option)
}

func (e *InvalidColorOptionError) Is(target error) bool {
	return e.Error() == target.Error()
}

type InvalidColorError struct {
	ConfigModule string
	Option       string
	Value        Color
}

func (e *InvalidColorError) Error() string {
	return fmt.Sprintf("(%v.%v) %v does not match ANSI nor HEX RGB", e.ConfigModule, e.Option, e.Value)
}

func (e *InvalidColorError) Is(target error) bool {
	return e.Error() == target.Error()
}

func (c *Colors) Validate() (bool, []error) {
	errors := []error{}
	for key, val := range *c {
		if _, err := tilecontent.FromString(key); err != nil {
			errors = append(errors, &InvalidColorOptionError{key})
		}
		if !val.IsValid() {
			errors = append(errors, &InvalidColorError{"colors", key, val})
		}
	}
	return len(errors) == 0, errors
}

func (c *Colors) Apply() {
	for key, color := range *c {
		tileContent, _ := tilecontent.FromString(key)
		styles.SetTileColor(tileContent, string(color))
	}
}
