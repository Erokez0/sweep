package actions_test

import (
	actions "sweep/shared/consts/action"
	"testing"
)

func Test(t *testing.T) {
	type TestCase struct {
		keyPress string
		expected error
	}
	testCases := []TestCase{
		{
			keyPress: "gg",
			expected: nil,
		},
	}
	for _, testCase := range testCases {
		_, err := actions.GetAction(testCase.keyPress)
		if err != testCase.expected {
			t.Error(err)
		}
	}
}
