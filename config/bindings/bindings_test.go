package bindings

import (
	"errors"
	"sweep/shared/consts/actions"
	"testing"
)

func Test_Validate(t *testing.T) {
	type TestCase struct {
		errs     []error
		isValid  bool
		bindings Bindings
	}

	testCases := []TestCase{

		{
			isValid: false,
			bindings: Bindings{
				"foo": []string{"zz"},
			},
			errs: []error{
				&InvalidActionError{"foo"},
			},
		},
		{
			isValid: false,
			bindings: Bindings{
				actions.MoveCursorLeft: []string{"ctrl+ctrl+foo+bar"},
			},
			errs: []error{
				&InvalidKeyPressPatternError{
					action:  actions.MoveCursorLeft,
					index:   0,
					binding: "ctrl+ctrl+foo+bar",
				},
			},
		},
		{
			isValid: false,
			bindings: Bindings{
				"foo":                   []string{"zz"},
				actions.MoveCursorRight: []string{"ctrl+ctrl+foo+bar"},
			},
			errs: []error{
				&InvalidActionError{"foo"},
				&InvalidKeyPressPatternError{
					action:  actions.MoveCursorRight,
					index:   0,
					binding: "ctrl+ctrl+foo+bar",
				},
			},
		},
		{
			isValid: true,
			bindings: Bindings{
				actions.MoveCursorDown: []string{"j"},
			},
			errs: []error{},
		},
		{
			isValid: true,
			bindings: Bindings{
				actions.MoveCursorToBottomRow: []string{"ctrl+down"},
			},
			errs: []error{},
		},
		{
			isValid: true,
			bindings: Bindings{
				actions.MoveCursorToTopRow: []string{"alt+up"},
			},
			errs: []error{},
		},
		{
			isValid: true,
			bindings: Bindings{
				actions.MoveCursorToFirstColumn: []string{"alt+ctrl+left"},
			},
			errs: []error{},
		},
		{
			isValid: true,
			bindings: Bindings{
				actions.MoveCursorToLastColumn: []string{"ctrl+shift+right"},
			},
			errs: []error{},
		},
	}

	for _, testCase := range testCases {
		isValid, errs := testCase.bindings.Validate()

		if isValid != testCase.isValid {
			t.Errorf("[Assertion failed] isValid\nExpected: %v\nActual: %v\n", testCase.isValid, isValid)
		}

		areErrorsEqual := func() bool {
			if len(errs) != len(testCase.errs) {
				return false
			}
			for ix := range errs {
				if !errors.Is(errs[ix], testCase.errs[ix]) {
					return false
				}
			}
			return true
		}()

		if !areErrorsEqual {
			t.Errorf("[Assertion failed] Errors should be equal\nExpected: %v\nActual: %v\n", testCase.errs, errs)
		}
	}
}

func Test_Apply(t *testing.T) {
	type TestCase struct {
		bindings Bindings
	}

	testCases := []TestCase{
		{
			bindings: Bindings{
				actions.MoveCursorLeft: []string{"l"},
			},
		},
	}

	for _, testCase := range testCases {
		testCase.bindings.Apply()

		for expected, keyPresses := range testCase.bindings {
			for _, keyPress := range keyPresses {

				actual, err := actions.GetAction(keyPress)
				if err != nil {
					t.Errorf("Unexpected error: %v\nExpected action: %v\nKey press: %v", err, expected, keyPress)
				}
				if actual.Kind != expected {
					t.Errorf("[Assertion failed] Kind\nExpected kind: %v\n Actual: %v\n", expected, actual.Kind)
				}
			}

		}
	}
}
