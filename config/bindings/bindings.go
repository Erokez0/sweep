package bindings

import (
	"fmt"

	actions "sweep/shared/consts/action"
	keyactionmap "sweep/shared/vars/key-action-map"
	regexes "sweep/shared/vars/regexes"
)

type Bindings map[actions.Action][]string

func (b Bindings) Apply() {
	for action, bindings := range b {
		for _, key  := range bindings {
			keyactionmap.KeyActionMap[key] = action
		}
	}
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
