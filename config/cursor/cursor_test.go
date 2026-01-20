package cursor

import (
	"errors"
	"sweep/config/colors"
	"testing"
)

func Test_Validate(t *testing.T) {
	type TestCase struct {
		cursor  Cursor
		isValid bool
		errs    []error
	}

	testCases := []TestCase{

		{
			cursor: Cursor{
				Color:     "foo",
				LeftHalf:  "",
				RightHalf: "",
			},
			isValid: false,
			errs: []error{
				&colors.InvalidColorError{
					ConfigModule: configModule,
					Option:       option,
					Value:        "foo",
				},
			},
		},
		{
			cursor: Cursor{
				Color:     "#FF",
				LeftHalf:  "bar",
				RightHalf: "",
			},
			isValid: false,
			errs: []error{
				&CursorHalfTooLongError{
					leftHalf,
				},
			},
		},
		{
			cursor: Cursor{
				Color:     "#FF",
				LeftHalf:  "",
				RightHalf: "baz",
			},
			isValid: false,
			errs: []error{
				&CursorHalfTooLongError{
					rightHalf,
				},
			},
		},
		{
			cursor: Cursor{
				Color:     "foo",
				LeftHalf:  "bar",
				RightHalf: "baz",
			},
			isValid: false,
			errs: []error{
				&colors.InvalidColorError{
					ConfigModule: configModule,
					Option:       option,
					Value:        "foo",
				},
				&CursorHalfTooLongError{
					leftHalf,
				},
				&CursorHalfTooLongError{
					rightHalf,
				},
			},
		},
	}

	for _, testCase := range testCases {
		isValid, errs := testCase.cursor.Validate()
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
			t.Errorf("[Assertion failed] errors should be equal\nExpected: %v\nActual: %v\n", testCase.errs, errs)
		}

	}
}
