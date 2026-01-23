package actions

import (
	"errors"
	"testing"
)

func Test_GetAction(t *testing.T) {
	type TestCase struct {
		keyPress   string
		expected   error
		quantifier uint16
		prepare    func()
	}

	testCases := []TestCase{
		{
			keyPress:   "gg",
			expected:   nil,
			quantifier: 1,
			prepare: func() {
				MoveCursorToTopRow.SetBinding("gg")
			},
		},
		{
			keyPress:   "G",
			expected:   nil,
			quantifier: 1,
			prepare: func() {
				MoveCursorToTopRow.SetBinding("G")
			},
		},
		{
			keyPress:   "0",
			expected:   nil,
			quantifier: 1,
			prepare: func() {
				MoveCursorToFirstColumn.SetBinding("0")
			},
		},
		{
			keyPress:   "$",
			expected:   nil,
			quantifier: 1,
			prepare: func() {
				MoveCursorToFirstColumn.SetBinding("$")
			},
		},
		{
			keyPress:   "a",
			quantifier: 1,
			expected:   &InvalidBindError{"a"},
			prepare:    func() {},
		},
		{
			keyPress:   "12j",
			quantifier: 12,
			expected:   nil,
			prepare: func() {
				MoveCursorDown.SetBinding("j")
			},
		},
	}

	for n, testCase := range testCases {
		bindingsMap = map[string]ActionType{}
		testCase.prepare()
		action, err := GetAction(testCase.keyPress)
		if !errors.Is(err, testCase.expected) {
			t.Errorf("[Assertion failed] #%v error\nexpected: %v, actual: %v", n+1, testCase.expected, err)
		}
		if action != nil && action.Quantifier != testCase.quantifier {
			t.Errorf("[Assertion failed] #%v quantifier\nexpected: %v, actual: %v", n+1, testCase.quantifier, action.Quantifier)
		}
	}
}

func Test_AnyBindingStartWith(t *testing.T) {
	type TestCase struct {
		bindingsMap map[string]ActionType
		keyStrokes  string
		expected    bool
	}

	testCases := []TestCase{
		{
			bindingsMap: map[string]ActionType{
				"gg": MoveCursorToTopRow,
			},
			keyStrokes: "jj",
			expected:   false,
		},
		{
			bindingsMap: map[string]ActionType{
				"G": MoveCursorToTopRow,
			},
			keyStrokes: "65538g",
			expected:   false,
		},
		{
			bindingsMap: map[string]ActionType{
				"gg": MoveCursorToTopRow,
			},
			keyStrokes: "g",
			expected:   true,
		},
		{
			bindingsMap: map[string]ActionType{
				"0": MoveCursorToFirstColumn,
			},
			keyStrokes: "0",
			expected:   true,
		},
		{
			bindingsMap: map[string]ActionType{
				"G": MoveCursorToBottomRow,
			},
			keyStrokes: "2G",
			expected:   true,
		},
		{
			bindingsMap: map[string]ActionType{
				"j": MoveCursorDown,
			},
			keyStrokes: "22",
			expected:   true,
		},
		{
			bindingsMap: map[string]ActionType{
				"h": MoveCursorLeft,
			},
			keyStrokes: "22",
			expected:   true,
		},
	}

	for n, testCase := range testCases {
		bindingsMap = testCase.bindingsMap
		actual := AnyBindingStartWith(testCase.keyStrokes)
		if actual != testCase.expected {
			t.Errorf("[Assertion failed] #%v\nExpected: %v\nActual: %v\n", n+1, testCase.expected, actual)
		}
	}
}
