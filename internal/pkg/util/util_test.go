package util_test

import (
	"reflect"
	"testing"

	"github.com/Angelos-Giannis/erbuilder/internal/pkg/util"
)

func TestGetCaseOfString(t *testing.T) {
	utility := util.New()
	testCases := map[string]struct {
		inputInitialValue  string
		inputConvertToCase string
		expectedOutput     string
	}{
		"Normal execution to convert to snake case": {
			inputInitialValue:  "SomeString",
			inputConvertToCase: "snake_case",
			expectedOutput:     "some_string",
		},
		"Normal execution to convert to camel case": {
			inputInitialValue:  "Some_string",
			inputConvertToCase: "camelCase",
			expectedOutput:     "someString",
		},
		"Normal execution to convert to screaming snake case": {
			inputInitialValue:  "SomeString",
			inputConvertToCase: "screaming_snake_case",
			expectedOutput:     "SOME_STRING",
		},
		"Normal execution to convert to kebab case": {
			inputInitialValue:  "SomeString",
			inputConvertToCase: "kebab_case",
			expectedOutput:     "some-string",
		},
		"Normal execution to convert to non expected case, so that to fall to the default one (snake case)": {
			inputInitialValue:  "SomeString",
			inputConvertToCase: "invalid-case",
			expectedOutput:     "some_string",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actualOutput := utility.GetCaseOfString(tc.inputInitialValue, tc.inputConvertToCase)
			if !reflect.DeepEqual(tc.expectedOutput, actualOutput) {
				t.Errorf("Expected to get '%v' as response but got '%v'.", tc.expectedOutput, actualOutput)
			}
		})
	}
}

func TestGetValueCount(t *testing.T) {
	utility := util.New()
	testCases := map[string]struct {
		inputToPlural     bool
		inputInitialValue string
		expectedOutput    string
	}{
		"Normal execution to get the singular of a singular string": {
			inputToPlural:     false,
			inputInitialValue: "example",
			expectedOutput:    "example",
		},
		"Normal execution to get the plural of a singular string": {
			inputToPlural:     true,
			inputInitialValue: "example",
			expectedOutput:    "examples",
		},
		"Normal execution to get the singular of a plural string": {
			inputToPlural:     false,
			inputInitialValue: "examples",
			expectedOutput:    "example",
		},
		"Normal execution to get the plural of a plural string": {
			inputToPlural:     true,
			inputInitialValue: "examples",
			expectedOutput:    "examples",
		},
		"Normal execution to get the plural of a singular snake case string": {
			inputToPlural:     true,
			inputInitialValue: "some_example",
			expectedOutput:    "some_examples",
		},
		"Normal execution to get the singular of a plural snake case string": {
			inputToPlural:     false,
			inputInitialValue: "some_examples",
			expectedOutput:    "some_example",
		},
		"Normal execution to get the plural of a singular camel case string": {
			inputToPlural:     true,
			inputInitialValue: "someExample",
			expectedOutput:    "someExamples",
		},
		"Normal execution to get the singular of a plural camel case string": {
			inputToPlural:     false,
			inputInitialValue: "someExamples",
			expectedOutput:    "someExample",
		},
		"Normal execution to get the singular of a plural screaming snake case string": {
			inputToPlural:     false,
			inputInitialValue: "SOME_EXAMPLES",
			expectedOutput:    "SOME_EXAMPLE",
		},
		"Normal execution to get the plural of a singular screaming snake case string": {
			inputToPlural:     true,
			inputInitialValue: "SOME_EXAMPLE",
			expectedOutput:    "SOME_EXAMPLES",
		},
		"Normal execution to get the singular of a plural kebab case string": {
			inputToPlural:     false,
			inputInitialValue: "some-examples",
			expectedOutput:    "some-example",
		},
		"Normal execution to get the plural of a singular kebab case string": {
			inputToPlural:     true,
			inputInitialValue: "some-example",
			expectedOutput:    "some-examples",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actualOutput := utility.GetValueCount(tc.inputToPlural, tc.inputInitialValue)
			if !reflect.DeepEqual(tc.expectedOutput, actualOutput) {
				t.Errorf("Expected to get '%v' as response but got '%v'.", tc.expectedOutput, actualOutput)
			}
		})
	}
}
