package actions

import (
	"errors"
	"testing"
)

func Test(t *testing.T) {
	type TestCase struct {
		keyPress   string
		expected   error
		multiplier uint16
		prepare    func()
	}

	testCases := []TestCase{
		{
			keyPress:   "gg",
			expected:   nil,
			multiplier: 1,
			prepare: func() {
				MoveCursorToTopRow.SetBinding("gg")
			},
		},
		{
			keyPress:   "G",
			expected:   nil,
			multiplier: 1,
			prepare: func() {
				MoveCursorToTopRow.SetBinding("G")
			},
		},
		{
			keyPress:   "0",
			expected:   nil,
			multiplier: 1,
			prepare: func() {
				MoveCursorToFirstColumn.SetBinding("0")
			},
		},
		{
			keyPress:   "$",
			expected:   nil,
			multiplier: 1,
			prepare: func() {
				MoveCursorToFirstColumn.SetBinding("$")
			},
		},
		{
			keyPress:   "a",
			multiplier: 1,
			expected:   &InvalidBindError{"a"},
			prepare:    func() {},
		},
		{
			keyPress:   "12j",
			multiplier: 12,
			expected:   nil,
			prepare: func() {
				MoveCursorDown.SetBinding("j")
			},
		},
	}

	for _, testCase := range testCases {
		bindingsMap = map[string]ActionType{}
		testCase.prepare()
		action, err := GetAction(testCase.keyPress)
		if !errors.Is(err, testCase.expected) {
			t.Errorf("[Assertion failed] error\nexpected: %v, actual: %v", testCase.expected, err)
		}
		if action != nil && action.Multiplier != testCase.multiplier {
			t.Errorf("[Assertion failed] multiplier\nexpected: %v, actual: %v", testCase.multiplier, action.Multiplier)
		}
	}
}
