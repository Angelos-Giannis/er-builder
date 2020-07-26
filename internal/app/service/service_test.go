package service_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/eujoy/erbuilder/internal/app/service"
	"github.com/eujoy/erbuilder/internal/config"
	"github.com/eujoy/erbuilder/internal/domain"
	"github.com/eujoy/erbuilder/internal/pkg/util"
	"github.com/eujoy/erbuilder/internal/pkg/writer"
	"github.com/udhos/equalfile"
	"github.com/urfave/cli/v2"
)

func TestNew(t *testing.T) {
	options := domain.Options{}
	actualService := service.New(options, util.New(), writer.New(util.New(), "some/path", "some_filename"))

	if reflect.TypeOf(&service.Service{}) != reflect.TypeOf(actualService) {
		t.Errorf("Expected to get type '%v' but got '%v'.", reflect.TypeOf(&service.Service{}), reflect.TypeOf(actualService))
	}
}

func TestGenerate(t *testing.T) {
	options := domain.Options{
		IDField:        "id",
		OutputFilename: "test-example-er-diagram",
		OutputPath:     "./../../../test",
		Tag:            "db",
		Title:          "example_db",
		ColumnNameCase: "snake_case",
		TableNameCase:  "snake_case",
		Config:         config.New(),
	}

	testCases := map[string]struct {
		providedOptions    domain.Options
		expectedOutputFile string
	}{
		"Generate .er file based on a list of provided files": {
			providedOptions: func() domain.Options {
				fileListStringSlice := cli.StringSlice{}
				err := fileListStringSlice.Set("./../../../test/example.go")
				if err != nil {
					t.Errorf("Expected to get nil as error but got '%v'.", err)
				}

				options.FileList = fileListStringSlice
				return options
			}(),
			expectedOutputFile: "./../../../test/example-er-diagram.er",
		},
		"Generate .er file based on a directory provided": {
			providedOptions: func() domain.Options {
				options.Directory = "./../../../test"
				return options
			}(),
			expectedOutputFile: "./../../../test/example-er-diagram.er",
		},
		"Generate .er file based on a directory provided including some common fields": {
			providedOptions: func() domain.Options {
				commonFieldsStringSlice := cli.StringSlice{}
				for _, cf := range []string{"created_at", "updated_at", "deleted_at"} {
					err := commonFieldsStringSlice.Set(cf)
					if err != nil {
						t.Errorf("Expected to get nil as error but got '%v'.", err)
					}
				}

				options.Directory = "./../../../test"
				options.CommonFields = commonFieldsStringSlice
				return options
			}(),
			expectedOutputFile: "./../../../test/example-er-diagram-with-common-fields.er",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actualService := service.New(tc.providedOptions, util.New(), writer.New(util.New(), tc.providedOptions.OutputPath, tc.providedOptions.OutputFilename))
			validateGenerateExecution(t, actualService, tc.expectedOutputFile)
		})
	}
}

func validateGenerateExecution(t *testing.T, actualService *service.Service, expectedOutputFile string) {
	err := actualService.Generate()
	if err != nil {
		t.Errorf("Expected to get nil as error but got '%v'.", err)
	}

	cmp := equalfile.New(nil, equalfile.Options{})
	equal, err := cmp.CompareFile(expectedOutputFile, "./../../../test/test-example-er-diagram.er")
	if err != nil {
		t.Errorf("Expected to get nil as error but got '%v'.", err)
	}

	if !equal {
		t.Errorf("Expected that the generated file and the example test file are the same but they were not.")
	}

	err = os.Remove("./../../../test/test-example-er-diagram.er")
	if err != nil {
		t.Errorf("Expected to get nil as error when deleting the test file but got '%v'.", err)
	}
}
