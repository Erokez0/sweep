package colors

import (
	"fmt"
	tilecontent "sweep/shared/consts/tile-content"
	regexes "sweep/shared/vars/regexes"
	"sweep/tui/styles"
)

type Color string

func (c Color) IsSet() bool {
	return c != ""
}

func (c Color) IsValid() bool {
	return !c.IsSet() || regexes.ColorRegex.MatchString(string(c))
}

type Colors map[string]Color

func (c *Colors) Validate() (bool, []string) {
	errors := []string{}
	for key, val := range *c {
		if !val.IsValid() {
			errors = append(errors, fmt.Sprintf("(colors.%v) %v does not match ANSI nor HEX RGB", key, val))
		}
		if _, err := tilecontent.FromString(key); err != nil {
			errors = append(errors, fmt.Sprintf("(colors) %v is not a valid option", key))
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
