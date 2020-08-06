package service_test

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/eujoy/erbuilder/internal/app/service"
	"github.com/eujoy/erbuilder/internal/config"
	"github.com/eujoy/erbuilder/internal/domain"
	"github.com/eujoy/erbuilder/internal/pkg/util"
	"github.com/eujoy/erbuilder/internal/pkg/writer"
	"github.com/eujoy/erbuilder/test/mock"
	"github.com/udhos/equalfile"
	"github.com/urfave/cli/v2"
)

func TestNew(t *testing.T) {
	options := domain.Options{}
	actualService := service.New(options, &mock.Survey{}, util.New(), writer.New(util.New(), "some/path", "some_filename"))

	if reflect.TypeOf(&service.Service{}) != reflect.TypeOf(actualService) {
		t.Errorf("Expected to get type '%v' but got '%v'.", reflect.TypeOf(&service.Service{}), reflect.TypeOf(actualService))
	}
}

func TestGenerate(t *testing.T) {
	options := domain.Options{
		IDField:        "id",
		OutputPath:     "./../../../test",
		Tag:            "db",
		Title:          "example_db",
		ColumnNameCase: "snake_case",
		TableNameCase:  "snake_case",
		Config:         config.New(),
	}

	type setupFunc func(options domain.Options) *service.Service

	testCases := map[string]struct {
		setupFn            setupFunc
		providedOptions    domain.Options
		filenameSuffix     string
		expectedOutputFile string
	}{
		"Generate .er file based on a list of provided files": {
			setupFn: defaultGenerateTestSetupFunc,
			providedOptions: func(filenameSuffix string) domain.Options {
				fileListStringSlice := cli.StringSlice{}
				err := fileListStringSlice.Set("./../../../test/example.go")
				if err != nil {
					t.Errorf("Expected to get nil as error but got '%v'.", err)
				}

				testOptions := options
				testOptions.OutputFilename = fmt.Sprintf("test-er-diagram-%v", filenameSuffix)
				testOptions.FileList = fileListStringSlice
				return testOptions
			}("normal-execution-with-provided-list-of-files"),
			filenameSuffix:     "normal-execution-with-provided-list-of-files",
			expectedOutputFile: "./../../../test/example-er-diagram.er",
		},
		"Generate .er file based on a directory provided": {
			setupFn: defaultGenerateTestSetupFunc,
			providedOptions: func(filenameSuffix string) domain.Options {
				testOptions := options
				testOptions.Directory = "./../../../test"
				return testOptions
			}("normal-execution-with-directory-provided"),
			filenameSuffix:     "normal-execution-with-directory-provided",
			expectedOutputFile: "./../../../test/example-er-diagram.er",
		},
		"Generate .er file based on a directory provided including some common fields": {
			setupFn: defaultGenerateTestSetupFunc,
			providedOptions: func(filenameSuffix string) domain.Options {
				commonFieldsStringSlice := cli.StringSlice{}
				for _, cf := range []string{"created_at", "updated_at", "deleted_at"} {
					err := commonFieldsStringSlice.Set(cf)
					if err != nil {
						t.Errorf("Expected to get nil as error but got '%v'.", err)
					}
				}

				testOptions := options
				testOptions.OutputFilename = fmt.Sprintf("test-er-diagram-%v", filenameSuffix)
				testOptions.Directory = "./../../../test"
				testOptions.CommonFields = commonFieldsStringSlice
				return testOptions
			}("normal-execution-with-directory-provided-and-common-fields"),
			filenameSuffix:     "normal-execution-with-directory-provided-and-common-fields",
			expectedOutputFile: "./../../../test/example-er-diagram-with-common-fields.er",
		},
		"Generate .er file from a directory including some extra tables definition": {
			setupFn: defaultGenerateTestSetupFunc,
			providedOptions: func(filenameSuffix string) domain.Options {
				testOptions := options
				testOptions.OutputFilename = fmt.Sprintf("test-er-diagram-%v", filenameSuffix)
				testOptions.Directory = "./../../../test"
				testOptions.ExtraTablesDefinition = `[{"name":"schema_migrations","columns":[{"name":"id","type":"integer","is_primary_key":true,"is_foreign_key":false,"is_extra_field":false},{"name":"version","type":"varchar","is_primary_key":false,"is_foreign_key":false,"is_extra_field":false}],"color":"#ebe486"}]`
				return testOptions
			}("include-extra-tables-definition"),
			filenameSuffix:     "include-extra-tables-definition",
			expectedOutputFile: "./../../../test/example-er-diagram-with-extra-tables.er",
		},
		//
		// @todo Add test scenario for generating file and also include survey for extra tables.
		//
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actualService := tc.setupFn(tc.providedOptions)
			validateGenerateExecution(t, actualService, tc.providedOptions.OutputFilename, tc.expectedOutputFile)
		})
	}
}

func defaultGenerateTestSetupFunc(options domain.Options) *service.Service {
	return service.New(options, nil, util.New(), writer.New(util.New(), options.OutputPath, options.OutputFilename))
}

func validateGenerateExecution(t *testing.T, actualService *service.Service, actualOutputFileName, expectedOutputFile string) {
	err := actualService.Generate()
	if err != nil {
		t.Errorf("Expected to get nil as error but got '%v'.", err)
	}

	cmp := equalfile.New(nil, equalfile.Options{})
	equal, err := cmp.CompareFile(expectedOutputFile, fmt.Sprintf("./../../../test/%v.er", actualOutputFileName))
	if err != nil {
		t.Errorf("Expected to get nil as error but got '%v'.", err)
	}

	if !equal {
		t.Errorf("Expected that the generated file and the example test file are the same but they were not.")
	}

	err = os.Remove(fmt.Sprintf("./../../../test/%v.er", actualOutputFileName))
	if err != nil {
		t.Errorf("Expected to get nil as error when deleting the test file but got '%v'.", err)
	}
}

func TestBuild(t *testing.T) {
	defaultTableAnswer := domain.TableAnswer{
		Name:  "my_table",
		Color: "#aabbcc",
	}

	type setupFunc func(expectedError error) *service.Service
	options := domain.Options{ExtraTablesSurvey: true}

	testCases := map[string]struct {
		expectedDefinition []domain.Table
		expectedError      error
		setupFn            setupFunc
	}{
		"Normal execution providing one table with one column": {
			expectedDefinition: []domain.Table{
				{
					Name: "my_table",
					ColumnList: []domain.Column{
						{
							Name:         "my_column",
							Type:         "integer",
							IsPrimaryKey: true,
							IsForeignKey: false,
							IsExtraField: false,
						},
					},
					Color: "#aabbcc",
				},
			},
			expectedError: nil,
			setupFn: func(expectedError error) *service.Service {
				mockSurvey := mock.Survey{}
				mockSurvey.On("AskTableDetails").Return(defaultTableAnswer, expectedError)
				mockSurvey.On("AskColumnDetails").Return(domain.ColumnAnswer{
					Name:         "my_column",
					Type:         "integer",
					IsPrimaryKey: true,
					IsForeignKey: false,
					AddMore:      "Nothing",
				}, expectedError)

				return service.New(options, &mockSurvey, util.New(), nil)
			},
		},
		"Failed to get table definition": {
			expectedDefinition: []domain.Table{},
			expectedError:      errors.New("failed to get table details"),
			setupFn: func(expectedError error) *service.Service {
				mockSurvey := mock.Survey{}
				mockSurvey.On("AskTableDetails").Return(domain.TableAnswer{}, expectedError)
				mockSurvey.On("AskColumnDetails").Times(0)

				return service.New(options, &mockSurvey, util.New(), nil)
			},
		},
		"Failed to get column definition": {
			expectedDefinition: []domain.Table{},
			expectedError:      errors.New("failed to get column details"),
			setupFn: func(expectedError error) *service.Service {
				mockSurvey := mock.Survey{}
				mockSurvey.On("AskTableDetails").Return(defaultTableAnswer, nil)
				mockSurvey.On("AskColumnDetails").Return(domain.ColumnAnswer{}, expectedError)

				return service.New(options, &mockSurvey, util.New(), nil)
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			srv := tc.setupFn(tc.expectedError)
			actualDefinition, actualError := srv.Build()

			if !reflect.DeepEqual(tc.expectedError, actualError) {
				t.Errorf("Expected to get '%v' as error but got '%v'.", tc.expectedError, actualError)
			}

			if !reflect.DeepEqual(tc.expectedDefinition, actualDefinition) {
				t.Errorf("Expected to get '%v' as response but got '%v'.", tc.expectedDefinition, actualDefinition)
			}
		})
	}
}
