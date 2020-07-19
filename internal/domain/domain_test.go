package domain_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/eujoy/erbuilder/internal/config"
	"github.com/eujoy/erbuilder/internal/domain"
	"github.com/urfave/cli/v2"
)

func TestValidate(t *testing.T) {
	cfg := config.New()

	testCases := map[string]struct {
		options       domain.Options
		expectedError error
	}{
		"Normal setup with valid values": {
			options: domain.Options{
				Directory:      "./../../../test/",
				IDField:        "id",
				FileList:       cli.StringSlice{},
				OutputFilename: "test-example-er-diagram",
				OutputPath:     "./../../../test",
				Tag:            "db",
				Title:          "example_db",
				ColumnNameCase: "snake_case",
				TableNameCase:  "snake_case",
				Config:         cfg,
			},
			expectedError: nil,
		},
		"Attempt execution without providing any of directory and list of files": {
			options: domain.Options{
				Directory:      "",
				IDField:        "id",
				FileList:       cli.StringSlice{},
				OutputFilename: "test-example-er-diagram",
				OutputPath:     "./../../../test",
				Tag:            "db",
				Title:          "example_db",
				ColumnNameCase: "snake_case",
				TableNameCase:  "snake_case",
				Config:         cfg,
			},
			expectedError: errors.New("Need to provide at least one of 'directory' or 'file_list'"),
		},
		"Attempt execution by providing invalid value for column name case": {
			options: domain.Options{
				Directory:      "./../../../test/",
				IDField:        "id",
				FileList:       cli.StringSlice{},
				OutputFilename: "test-example-er-diagram",
				OutputPath:     "./../../../test",
				Tag:            "db",
				Title:          "example_db",
				ColumnNameCase: "invalid_case",
				TableNameCase:  "snake_case",
				Config:         cfg,
			},
			expectedError: fmt.Errorf(
				"The provided value for column name case is not valid. Allowed values : %v",
				cfg.Settings.AllowedColumnNameCaseValues,
			),
		},
		"Attempt execution by providing invalid value for table name case": {
			options: domain.Options{
				Directory:      "./../../../test/",
				IDField:        "id",
				FileList:       cli.StringSlice{},
				OutputFilename: "test-example-er-diagram",
				OutputPath:     "./../../../test",
				Tag:            "db",
				Title:          "example_db",
				ColumnNameCase: "snake_case",
				TableNameCase:  "invalid_case",
				Config:         config.New(),
			},
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

func TestGetCommonFields(t *testing.T) {
	expectedOptionName := "common_field"

	options := &domain.Options{}
	actualFlag := options.GetCommonFields()

	if expectedOptionName != fmt.Sprintf("%v", actualFlag.Name) {
		t.Errorf("Expected to get a flag with name '%v' but got one with flag name '%v'.", expectedOptionName, actualFlag.Name)
	}
	if reflect.TypeOf(&cli.StringSliceFlag{}) != reflect.TypeOf(actualFlag) {
		t.Errorf("Expected to get a flag of type '%v' but got one of type '%v'.", reflect.TypeOf(&cli.StringSliceFlag{}), reflect.TypeOf(actualFlag))
	}
}

func TestGetDirectoryFlag(t *testing.T) {
	expectedOptionName := "directory"

	options := &domain.Options{}
	actualFlag := options.GetDirectoryFlag()

	if expectedOptionName != fmt.Sprintf("%v", actualFlag.Name) {
		t.Errorf("Expected to get a flag with name '%v' but got one with flag name '%v'.", expectedOptionName, actualFlag.Name)
	}
	if reflect.TypeOf(&cli.StringFlag{}) != reflect.TypeOf(actualFlag) {
		t.Errorf("Expected to get a flag of type '%v' but got one of type '%v'.", reflect.TypeOf(&cli.StringFlag{}), reflect.TypeOf(actualFlag))
	}
}

func TestGetFileList(t *testing.T) {
	expectedOptionName := "file_list"

	options := &domain.Options{}
	actualFlag := options.GetFileList()

	if expectedOptionName != fmt.Sprintf("%v", actualFlag.Name) {
		t.Errorf("Expected to get a flag with name '%v' but got one with flag name '%v'.", expectedOptionName, actualFlag.Name)
	}
	if reflect.TypeOf(&cli.StringSliceFlag{}) != reflect.TypeOf(actualFlag) {
		t.Errorf("Expected to get a flag of type '%v' but got one of type '%v'.", reflect.TypeOf(&cli.StringSliceFlag{}), reflect.TypeOf(actualFlag))
	}
}

func TestGetIDField(t *testing.T) {
	expectedOptionName := "id_field"

	options := &domain.Options{}
	actualFlag := options.GetIDField()

	if expectedOptionName != fmt.Sprintf("%v", actualFlag.Name) {
		t.Errorf("Expected to get a flag with name '%v' but got one with flag name '%v'.", expectedOptionName, actualFlag.Name)
	}
	if reflect.TypeOf(&cli.StringFlag{}) != reflect.TypeOf(actualFlag) {
		t.Errorf("Expected to get a flag of type '%v' but got one of type '%v'.", reflect.TypeOf(&cli.StringFlag{}), reflect.TypeOf(actualFlag))
	}
}

func TestGetOutputFilename(t *testing.T) {
	expectedOptionName := "output_filename"

	options := &domain.Options{}
	actualFlag := options.GetOutputFilename()

	if expectedOptionName != fmt.Sprintf("%v", actualFlag.Name) {
		t.Errorf("Expected to get a flag with name '%v' but got one with flag name '%v'.", expectedOptionName, actualFlag.Name)
	}
	if reflect.TypeOf(&cli.StringFlag{}) != reflect.TypeOf(actualFlag) {
		t.Errorf("Expected to get a flag of type '%v' but got one of type '%v'.", reflect.TypeOf(&cli.StringFlag{}), reflect.TypeOf(actualFlag))
	}
}

func TestGetOutputPath(t *testing.T) {
	expectedOptionName := "output_path"

	options := &domain.Options{}
	actualFlag := options.GetOutputPath()

	if expectedOptionName != fmt.Sprintf("%v", actualFlag.Name) {
		t.Errorf("Expected to get a flag with name '%v' but got one with flag name '%v'.", expectedOptionName, actualFlag.Name)
	}
	if reflect.TypeOf(&cli.StringFlag{}) != reflect.TypeOf(actualFlag) {
		t.Errorf("Expected to get a flag of type '%v' but got one of type '%v'.", reflect.TypeOf(&cli.StringFlag{}), reflect.TypeOf(actualFlag))
	}
}

func TestGetTag(t *testing.T) {
	expectedOptionName := "tag"

	options := &domain.Options{}
	actualFlag := options.GetTag()

	if expectedOptionName != fmt.Sprintf("%v", actualFlag.Name) {
		t.Errorf("Expected to get a flag with name '%v' but got one with flag name '%v'.", expectedOptionName, actualFlag.Name)
	}
	if reflect.TypeOf(&cli.StringFlag{}) != reflect.TypeOf(actualFlag) {
		t.Errorf("Expected to get a flag of type '%v' but got one of type '%v'.", reflect.TypeOf(&cli.StringFlag{}), reflect.TypeOf(actualFlag))
	}
}

func TestGetTitle(t *testing.T) {
	expectedOptionName := "title"

	options := &domain.Options{}
	actualFlag := options.GetTitle()

	if expectedOptionName != fmt.Sprintf("%v", actualFlag.Name) {
		t.Errorf("Expected to get a flag with name '%v' but got one with flag name '%v'.", expectedOptionName, actualFlag.Name)
	}
	if reflect.TypeOf(&cli.StringFlag{}) != reflect.TypeOf(actualFlag) {
		t.Errorf("Expected to get a flag of type '%v' but got one of type '%v'.", reflect.TypeOf(&cli.StringFlag{}), reflect.TypeOf(actualFlag))
	}
}

func TestGetColumnNameCase(t *testing.T) {
	expectedOptionName := "column_name_case"

	options := &domain.Options{}
	actualFlag := options.GetColumnNameCase()

	if expectedOptionName != fmt.Sprintf("%v", actualFlag.Name) {
		t.Errorf("Expected to get a flag with name '%v' but got one with flag name '%v'.", expectedOptionName, actualFlag.Name)
	}
	if reflect.TypeOf(&cli.StringFlag{}) != reflect.TypeOf(actualFlag) {
		t.Errorf("Expected to get a flag of type '%v' but got one of type '%v'.", reflect.TypeOf(&cli.StringFlag{}), reflect.TypeOf(actualFlag))
	}
}

func TestGetTableNameCase(t *testing.T) {
	expectedOptionName := "table_name_case"

	options := &domain.Options{}
	actualFlag := options.GetTableNameCase()

	if expectedOptionName != fmt.Sprintf("%v", actualFlag.Name) {
		t.Errorf("Expected to get a flag with name '%v' but got one with flag name '%v'.", expectedOptionName, actualFlag.Name)
	}
	if reflect.TypeOf(&cli.StringFlag{}) != reflect.TypeOf(actualFlag) {
		t.Errorf("Expected to get a flag of type '%v' but got one of type '%v'.", reflect.TypeOf(&cli.StringFlag{}), reflect.TypeOf(actualFlag))
	}
}

func TestGetTableNamePlural(t *testing.T) {
	expectedOptionName := "table_in_plural"

	options := &domain.Options{}
	actualFlag := options.GetTableNamePlural()

	if expectedOptionName != fmt.Sprintf("%v", actualFlag.Name) {
		t.Errorf("Expected to get a flag with name '%v' but got one with flag name '%v'.", expectedOptionName, actualFlag.Name)
	}
	if reflect.TypeOf(&cli.BoolFlag{}) != reflect.TypeOf(actualFlag) {
		t.Errorf("Expected to get a flag of type '%v' but got one of type '%v'.", reflect.TypeOf(&cli.BoolFlag{}), reflect.TypeOf(actualFlag))
	}
}
