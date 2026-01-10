package flags

import (
	"errors"
	"fmt"
	"os"
	envkeys "sweep/shared/consts/env-keys"
	"sweep/shared/types"
	"sweep/shared/vars/glyphs"
	"sweep/shared/vars/paths"
	"sweep/tui/styles"
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
		{
			args: []string{},
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

type FillFlagError struct {
	flag types.Flag
}

func (e *FillFlagError) Error() string {
	return fmt.Sprintf("The %v flag was expected to set fill for the styles", e.flag)
}

func (e *FillFlagError) Is(target error) bool {
	return e.Error() == target.Error()
}

type AsciiFlagError struct {
	flag       types.Flag
	glyphKey   string
	glyphValue string
}

func (e *AsciiFlagError) Error() string {
	return fmt.Sprintf("The %v flag was expected to set %v to %v", e.flag, e.glyphKey, e.glyphValue)
}

func Test_Apply(t *testing.T) {
	type TestCase struct {
		flags      Flags
		asExpected func() error
	}

	testCases := []TestCase{
		{
			flags: Flags{FILL},
			asExpected: func() error {
				if !styles.IsFillSet {
					return &FillFlagError{FILL}
				}
				return nil
			},
		},

		{
			flags: Flags{ASCII},
			asExpected: func() error {
				expectations := map[*string]string{
					&glyphs.Mine:      "M",
					&glyphs.Flag:      "F",
					&glyphs.WrongFlag: "W",
					&glyphs.Empty:     " ",
					&glyphs.Zero:      "x",
					&glyphs.One:       "1",
					&glyphs.Two:       "2",
					&glyphs.Three:     "3",
					&glyphs.Four:      "4",
					&glyphs.Five:      "5",
					&glyphs.Six:       "6",
					&glyphs.Seven:     "7",
					&glyphs.Eight:     "8",
				}

				for key, val := range expectations {
					if *key != val {
						return &AsciiFlagError{ASCII, *key, val}
					}
				}
				return nil
			},
		},
	}

	for _, testCase := range testCases {
		testCase.flags.Apply()
		if err := testCase.asExpected(); err != nil {
			t.Error(err)
		}
	}
}

func Test_GetFlagArgument(t *testing.T) {
	type TestCase struct {
		args     []string
		index    int
		expected string
	}

	testCases := []TestCase{
		{
			args:     []string{MINES, "12"},
			index:    0,
			expected: "12",
		},
		{
			args:     []string{"foo", "bar", "baz"},
			index:    1,
			expected: "baz",
		},
		{
			args:     []string{"foo", "bar", "baz", "fizzbuzz"},
			index:    2,
			expected: "fizzbuzz",
		},
	}

	for _, testCase := range testCases {
		actual := getFlagArgument(testCase.args, testCase.index)
		expected := testCase.expected

		if expected != actual {
			t.Errorf("[Assertion failed] was expecting %v to be %v", actual, expected)
		}
	}
}

type UnsetEnvVarError struct {
	envVar string
}

func (e *UnsetEnvVarError) Error() string {
	return fmt.Sprintf("%v was expected to be set", e.envVar)
}
func (e *UnsetEnvVarError) Is(target error) bool {
	return e.Error() == target.Error()
}

type IncorrectEnvVarError struct {
	envVar   string
	expected string
	actual   string
}

func (e *IncorrectEnvVarError) Error() string {
	return fmt.Sprintf("%v value was expected to be %v, but it actually is %v", e.envVar, e.expected, e.actual)
}

func (e *IncorrectEnvVarError) Is(target error) bool {
	return e.Error() == target.Error()
}

func Test_ApplyFromAgs(t *testing.T) {
	type TestCase struct {
		osArgs []string
		test   func() error
	}

	testCases := []TestCase{
		{
			osArgs: []string{"foo", THEME_PREVIEW},
			test: func() error {
				preview, ok := os.LookupEnv(envkeys.Preview)
				if !ok {
					return &UnsetEnvVarError{envkeys.Preview}
				}
				if preview != "true" {
					return &IncorrectEnvVarError{envkeys.Preview, "true", preview}
				}
				return nil
			},
		},
	}

	for _, testCase := range testCases {
		os.Args = testCase.osArgs

		ApplyFromArgs()
		if err := testCase.test(); err != nil {
			t.Error(err)
		}
	}
}
