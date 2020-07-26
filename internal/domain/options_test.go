package domain_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/eujoy/erbuilder/internal/config"
	"github.com/eujoy/erbuilder/internal/domain"
	"github.com/eujoy/erbuilder/test"
)

func TestNewOptions(t *testing.T) {
	expectedOptions := domain.Options{}
	actualOptions := domain.NewOptions(config.New())

	if reflect.TypeOf(expectedOptions) != reflect.TypeOf(actualOptions) {
		t.Errorf("Expected to get type '%v' but got '%v'.", reflect.TypeOf(expectedOptions), reflect.TypeOf(actualOptions))
	}
}

func TestValidate(t *testing.T) {
	cfg := config.New()
	dataBuilder := test.NewDataBuilder()

	testCases := map[string]struct {
		options       domain.Options
		expectedError error
	}{
		"Normal setup with valid values": {
			options:       dataBuilder.GetOptionsForOptionTest(cfg),
			expectedError: nil,
		},
		"Attempt execution without providing any of directory and list of files": {
			options: func() domain.Options {
				options := dataBuilder.GetOptionsForOptionTest(cfg)
				options.Directory = ""
				return options
			}(),
			expectedError: errors.New("Need to provide at least one of 'directory' or 'file_list'"),
		},
		"Attempt execution by providing invalid value for column name case": {
			options: func() domain.Options {
				options := dataBuilder.GetOptionsForOptionTest(cfg)
				options.ColumnNameCase = "invalid_case"
				return options
			}(),
			expectedError: fmt.Errorf(
				"The provided value for column name case is not valid. Allowed values : %v",
				cfg.Settings.AllowedColumnNameCaseValues,
			),
		},
		"Attempt execution by providing invalid value for table name case": {
			options: func() domain.Options {
				options := dataBuilder.GetOptionsForOptionTest(cfg)
				options.TableNameCase = "invalid_case"
				return options
			}(),
			expectedError: fmt.Errorf(
				"The provided value for table name case is not valid. Allowed values : %v",
				cfg.Settings.AllowedTableNameCaseValues,
			),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actualError := tc.options.Validate()
			if tc.expectedError != nil && tc.expectedError.Error() != actualError.Error() {
				t.Errorf("Expected to get '%v' as error but got '%v'.", tc.expectedError, actualError)
			}
		})
	}
}

func TestOptionFlags(t *testing.T) {
	options := &domain.Options{}

	t.Run("Test GetCommonFields", func(t *testing.T) {
		actualFlag := options.GetCommonFields()
		validateFlagIsAsExpected(t, "common_field", actualFlag.Name, "stringSliceFlag", reflect.TypeOf(actualFlag).Name())
	})

	t.Run("Test GetDirectoryFlag", func(t *testing.T) {
		actualFlag := options.GetDirectoryFlag()
		validateFlagIsAsExpected(t, "directory", actualFlag.Name, "stringFlag", reflect.TypeOf(actualFlag).Name())
	})

	t.Run("Test GetFileList", func(t *testing.T) {
		actualFlag := options.GetFileList()
		validateFlagIsAsExpected(t, "file_list", actualFlag.Name, "stringSliceFlag", reflect.TypeOf(actualFlag).Name())
	})

	t.Run("Test GetIDField", func(t *testing.T) {
		actualFlag := options.GetIDField()
		validateFlagIsAsExpected(t, "id_field", actualFlag.Name, "stringFlag", reflect.TypeOf(actualFlag).Name())
	})

	t.Run("Test GetOutputFilename", func(t *testing.T) {
		actualFlag := options.GetOutputFilename()
		validateFlagIsAsExpected(t, "output_filename", actualFlag.Name, "stringFlag", reflect.TypeOf(actualFlag).Name())
	})

	t.Run("Test GetOutputPath", func(t *testing.T) {
		actualFlag := options.GetOutputPath()
		validateFlagIsAsExpected(t, "output_path", actualFlag.Name, "stringFlag", reflect.TypeOf(actualFlag).Name())
	})

	t.Run("Test GetTag", func(t *testing.T) {
		actualFlag := options.GetTag()
		validateFlagIsAsExpected(t, "tag", actualFlag.Name, "stringFlag", reflect.TypeOf(actualFlag).Name())
	})

	t.Run("Test GetTitle", func(t *testing.T) {
		actualFlag := options.GetTitle()
		validateFlagIsAsExpected(t, "title", actualFlag.Name, "stringFlag", reflect.TypeOf(actualFlag).Name())
	})

	t.Run("Test GetColumnNameCase", func(t *testing.T) {
		actualFlag := options.GetColumnNameCase()
		validateFlagIsAsExpected(t, "column_name_case", actualFlag.Name, "stringFlag", reflect.TypeOf(actualFlag).Name())
	})

	t.Run("Test GetTableNameCase", func(t *testing.T) {
		actualFlag := options.GetTableNameCase()
		validateFlagIsAsExpected(t, "table_name_case", actualFlag.Name, "stringFlag", reflect.TypeOf(actualFlag).Name())
	})

	t.Run("Test GetTableNamePlural", func(t *testing.T) {
		actualFlag := options.GetTableNamePlural()
		validateFlagIsAsExpected(t, "table_in_plural", actualFlag.Name, "boolFlag", reflect.TypeOf(actualFlag).Name())
	})
}

func validateFlagIsAsExpected(t *testing.T, expectedFlagName, actualFlagName, expectedFlagType, actualFlagType string) {
	checkExpectedFlagName(t, expectedFlagName, actualFlagName)
	checkExpectedFlagType(t, expectedFlagType, actualFlagType)
}

func checkExpectedFlagName(t *testing.T, expectedFlagName, actualFlagName string) {
	if expectedFlagName != fmt.Sprintf("%v", actualFlagName) {
		t.Errorf("Expected to get a flag with name '%v' but got one with flag name '%v'.", expectedFlagName, actualFlagName)
	}
}

func checkExpectedFlagType(t *testing.T, expectedFlagType, actualFlagType string) {
	if expectedFlagType == actualFlagType {
		t.Errorf("Expected to get a flag of type '%v' but got one of type '%v'.", expectedFlagType, actualFlagType)
	}
}
