package colors

import (
	"errors"
	"testing"

	tilecontent "sweep/shared/consts/tile-content"
	styles "sweep/tui/styles"
)

func Test_ColorIsSet(t *testing.T) {
	type TestCase struct {
		color    Color
		expected bool
	}

	testCases := []TestCase{
		{
			color:    "",
			expected: false,
		},
		{
			color:    "12",
			expected: true,
		},
	}

	for _, testCase := range testCases {
		actual := testCase.color.IsSet()

		if actual != testCase.expected && testCase.expected {
			t.Errorf("[Assertion failed] %v was not expected to be set", testCase.color)
		}
		if actual != testCase.expected && !testCase.expected {
			t.Errorf("[Assertion failed] %v was expected to be set", testCase.color)
		}
	}
}

func Test_ColorIsValid(t *testing.T) {
	type TestCase struct {
		color    Color
		expected bool
	}

	testCases := []TestCase{
		{
			color:    "-1",
			expected: false,
		},
		{
			color:    "1024",
			expected: false,
		},
		{
			color:    "FF",
			expected: false,
		},
		{
			color:    "FFFFFF",
			expected: false,
		},
		{
			color:    "1",
			expected: true,
		},
		{
			color:    "255",
			expected: true,
		},
		{
			color:    "#FF",
			expected: true,
		},
		{
			color:    "#00",
			expected: true,
		},
		{
			color:    "#FF00",
			expected: true,
		},
		{
			color:    "#00FF00",
			expected: true,
		},
		{
			color:    "",
			expected: true,
		},
	}

	for _, testCase := range testCases {
		actual := testCase.color.IsValid()

		if actual != testCase.expected && testCase.expected {
			t.Errorf("[Assertion failed] %v was expected to be valid", testCase.color)
		}
		if actual != testCase.expected && !testCase.expected {
			t.Errorf("[Assertion failed] %v was not expected to be valid", testCase.color)
		}
	}
}

func Test_Validate(t *testing.T) {
	type Result struct {
		isValid bool
		errors  []error
	}

	type TestCase struct {
		colors   Colors
		expected Result
	}

	testCases := []TestCase{
		{
			colors: Colors{"foo": "#000000"},
			expected: Result{
				isValid: false,
				errors: []error{
					&InvalidColorOptionError{"foo"},
				},
			},
		},

		{
			colors: Colors{"0": "bar"},
			expected: Result{
				isValid: false,
				errors: []error{
					&InvalidColorError{"colors", "0", "bar"},
				},
			},
		},

		{
			colors: Colors{"foo": "bar"},
			expected: Result{
				isValid: false,
				errors: []error{
					&InvalidColorOptionError{"foo"},
					&InvalidColorError{"colors", "foo", "bar"},
				},
			},
		},

		{
			colors: Colors{
				"0":          "8",
				"1":          "12",
				"2":          "10",
				"3":          "3",
				"4":          "9",
				"5":          "13",
				"6":          "5",
				"7":          "1",
				"8":          "14",
				"mine":       "9",
				"wrong flag": "9",
				"flag":       "15",
				"empty":      "",
			},
			expected: Result{
				isValid: true,
				errors:  []error{},
			},
		},
	}

	for _, testCase := range testCases {

		isValid, errs := testCase.colors.Validate()

		actual := Result{
			isValid: isValid,
			errors:  errs,
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

		if actual.isValid != expected.isValid {
			t.Errorf("[Assertion failed] isValid\nExpected: %v\nActual: %v", expected.isValid, actual.isValid)
		}

		if !areErrorsEqual {
			t.Errorf("[Error equality check failed]\nExpected: %v\nActual: %v", expected.errors, actual.errors)
		}

	}
}

func Test_Apply(t *testing.T) {
	testCases := []Colors{
		{
			"0": "#000000",
		},
		{
			"1": "1",
		},
		{
			"2": "2",
		},
		{
			"3": "3",
		},
		{
			"4": "4",
		},
		{
			"5": "5",
		},
		{
			"6": "6",
		},
		{
			"7": "7",
		},
		{
			"8": "8",
		},
		{
			"mine": "9",
		},
		{
			"flag": "#FF",
		},
		{
			"wrong flag": "#FF0000",
		},
		{
			"empty": "",
		},
	}

	for _, testCase := range testCases {

		testCase.Apply()
		for key, color := range testCase {
			expected := styles.CreateTileStyle(string(color)).Value()
			tileContent, _ := tilecontent.FromString(key)
			actual := styles.GetTileStyle(tileContent).Value()
			if actual != expected {
				t.Errorf("[Assertion failed]\n%v should have been set to %v\n%v != %v", key, color, actual, expected)
			}
		}
	}
}
