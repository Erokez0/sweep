package bindings

import (
	"fmt"
	"slices"
	actions "sweep/shared/consts/action"
	"sweep/shared/vars/regexes"
)

type ActionHandler = func()

type Bindings map[actions.Action][]string

func (b Bindings) IsMoveCursorUp(input string) bool {
	return slices.Contains(b[actions.MoveCursorUp], input)
}
func (b Bindings) IsMoveCursorDown(input string) bool {
	return slices.Contains(b[actions.MoveCursorDown], input)
}
func (b Bindings) IsMoveCursorLeft(input string) bool {
	return slices.Contains(b[actions.MoveCursorLeft], input)
}
func (b Bindings) IsMoveCursorRight(input string) bool {
	return slices.Contains(b[actions.MoveCursorRight], input)
}
func (b Bindings) IsOpenTile(input string) bool {
	return slices.Contains(b[actions.OpenTile], input)
}
func (b Bindings) IsFlagTile(input string) bool {
	return slices.Contains(b[actions.FlagTile], input)
}

func (b Bindings) Validate() (bool, []string) {
	var errors []string
	for action, bindings := range b {
		if !actions.IsAction(string(action)) {
			errors = append(errors, fmt.Sprintf("(bindings) \"%v\" is not a valid action", string(action)))
		}
		for ix, binding := range bindings {
			if !regexes.KeyRegex.MatchString(binding) { 
				errors = append(errors, fmt.Sprintf("(bindings.%v.%v) \"%v\" does not match key press pattern", action, ix, binding))
			}
		}
	}
	return len(errors) == 0, errors
}
