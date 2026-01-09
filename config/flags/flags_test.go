package flags

import (
	"errors"
	"os"
	"sweep/shared/types"
	"testing"
)

func Test_FlagUint16Argument_Error(t *testing.T) {
	type TestCase struct {
		args  []string
		index int
		err   error
	}
	testCases := []TestCase{
		{
			args:  []string{MINES},
			index: 0,
			err:   &NoArgumentProvidedFlagError{MINES},
		},

		{
			args:  []string{MINES_SHORT, WIDTH_SHORT},
			index: 0,
			err:   &MustBeUin16FlagError{MINES_SHORT},
		},

		{
			args:  []string{HEIGHT, "-13"},
			index: 0,
			err:   &MustBeUin16FlagError{HEIGHT},
		},

		{
			args:  []string{HEIGHT_SHORT, "null"},
			index: 0,
			err:   &MustBeUin16FlagError{HEIGHT_SHORT},
		},

		{
			args:  []string{WIDTH, "true"},
			index: 0,
			err:   &MustBeUin16FlagError{WIDTH},
		},
	}

	for _, testCase := range testCases {
		err := validateFlagUint16Argument(testCase.args, testCase.index)
		expected := testCase.err
		if err == nil {
			t.Errorf("Should not be valid\nExpected error: %v", expected)
		}
		if !errors.Is(err, expected) {
			t.Errorf("Should not be valid\nExpected error: %v\nActual error: %v", expected, err)
		}
	}
}

func Test_FlagUint16Argument_Ok(t *testing.T) {
	type TestCase struct {
		args  []string
		index int
	}
	testCases := []TestCase{
		{
			args:  []string{MINES, "1"},
			index: 0,
		},

		{
			args:  []string{MINES_SHORT, "2"},
			index: 0,
		},

		{
			args:  []string{HEIGHT, "3"},
			index: 0,
		},

		{
			args:  []string{HEIGHT_SHORT, "4"},
			index: 0,
		},

		{
			args:  []string{WIDTH, "5"},
			index: 0,
		},

		{
			args:  []string{WIDTH_SHORT, "6"},
			index: 0,
		},
	}

	for _, testCase := range testCases {
		err := validateFlagUint16Argument(testCase.args, testCase.index)
		if err != nil {
			t.Errorf("Should be valid\nargs: %v, index: %v\nError: %v", testCase.args, testCase.index, err)
		}
	}

}

func Test_Validate(t *testing.T) {
	type Result struct {
		isValid bool
		errors  []error
	}
	type TestCase struct {
		args     []types.Flag
		expected Result
	}
	const INVALID string = "--invalid-flag"

	testCases := []TestCase{
		{
			args: []string{MINES},
			expected: Result{
				errors:  []error{&NoArgumentProvidedFlagError{MINES}},
				isValid: false,
			},
		},
		{
			args: []string{WIDTH, "-12"},
			expected: Result{
				errors:  []error{&MustBeUin16FlagError{WIDTH}},
				isValid: false,
			},
		},
		{
			args: []string{INVALID},
			expected: Result{
				errors:  []error{&InvalidFlagError{INVALID}},
				isValid: false,
			},
		},
		{
			args: []string{FILL},
			expected: Result{
				errors:  []error{},
				isValid: true,
			},
		},
	}

	for _, testCase := range testCases {
		var flags Flags = testCase.args

		cwd, _ := os.Getwd()
		os.Args = []string{cwd}
		isValid, errs := flags.Validate()
		actual := Result{
			errors:  errs,
			isValid: isValid,
		}
		expected := &testCase.expected

		areErrorsEqual := func(expected, actual []error) bool {
			if len(actual) != len(expected) {
				return false
			}

			areEqual := true
			for ix := range expected {

				if !errors.Is(actual[ix], expected[ix]) {
					return false
				}
			}
			return areEqual
		}(actual.errors, expected.errors)

		if actual.isValid != testCase.expected.isValid {
			t.Errorf("[Assertion failed] isValid\nExpected: %v\nActual: %v", expected.isValid, actual.isValid)
		}

		if !areErrorsEqual {
			t.Errorf("[Error equality check failed]\nExpected: %v\nActual: %v", expected.errors, actual.errors)
		}
	}
}
