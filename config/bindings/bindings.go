package bindings

import (
	"fmt"

	actions "sweep/shared/consts/actions"
	regexes "sweep/shared/vars/regexes"
)

type Bindings map[actions.ActionType][]string

func (b Bindings) Apply() {
	for action, bindings := range b {
		for _, key := range bindings {
			action.SetBinding(key)
		}
	}
}

type InvalidActionError struct {
	action actions.ActionType
}

func (e *InvalidActionError) Error() string {
	return fmt.Sprintf("(bindings) %v is not a valid action", e.action)
}

type InvalidKeyPressPatternError struct {
	action  actions.ActionType
	index   int
	binding string
}

func (e *InvalidKeyPressPatternError) Error() string {
	return fmt.Sprintf("(bindings.%v.%v) %v does bot match key press pattern", e.action, e.index, e.binding)
}

func (b Bindings) Validate() (bool, []error) {
	var errors []error
	for action, bindings := range b {
		if !actions.IsAction(string(action)) {
			errors = append(errors, &InvalidActionError{action})
		}
		for index, binding := range bindings {
			if !regexes.KeyRegex.MatchString(binding) {
				errors = append(errors, &InvalidKeyPressPatternError{action, index, binding})
			}
		}
	}
	return len(errors) == 0, errors
}
