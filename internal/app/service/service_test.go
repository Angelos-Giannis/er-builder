package service_test

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/AlecAivazis/survey/v2/terminal"
	expect "github.com/Netflix/go-expect"
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
		OutputPath:     "./../../../test",
		Tag:            "db",
		Title:          "example_db",
		ColumnNameCase: "snake_case",
		TableNameCase:  "snake_case",
		Config:         config.New(),
	}

	testCases := map[string]struct {
		providedOptions    domain.Options
		filenameSuffix     string
		expectedOutputFile string
	}{
		"Generate .er file based on a list of provided files": {
			providedOptions: func() domain.Options {
				fileListStringSlice := cli.StringSlice{}
				err := fileListStringSlice.Set("./../../../test/example.go")
				if err != nil {
					t.Errorf("Expected to get nil as error but got '%v'.", err)
				}

				testOptions := options

				testOptions.FileList = fileListStringSlice
				return testOptions
			}(),
			filenameSuffix:     "normal-execution-with-provided-list-of-files",
			expectedOutputFile: "./../../../test/example-er-diagram.er",
		},
		"Generate .er file based on a directory provided": {
			providedOptions: func() domain.Options {
				testOptions := options
				testOptions.Directory = "./../../../test"
				return testOptions
			}(),
			filenameSuffix:     "normal-execution-with-directory-provided",
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

				testOptions := options
				testOptions.Directory = "./../../../test"
				testOptions.CommonFields = commonFieldsStringSlice
				return testOptions
			}(),
			filenameSuffix:     "normal-execution-with-directory-provided-and-common-fields",
			expectedOutputFile: "./../../../test/example-er-diagram-with-common-fields.er",
		},
		"Generate .er file from a directory including some extra tables definition": {
			providedOptions: func() domain.Options {
				testOptions := options
				testOptions.Directory = "./../../../test"
				testOptions.ExtraTablesDefinition = `[{"name":"schema_migrations","columns":[{"name":"id","type":"integer","is_primary_key":true,"is_foreign_key":false,"is_extra_field":false},{"name":"version","type":"varchar","is_primary_key":false,"is_foreign_key":false,"is_extra_field":false}],"color":"#ebe486"}]`
				return testOptions
			}(),
			filenameSuffix:     "include-extra-tables-definition",
			expectedOutputFile: "./../../../test/example-er-diagram-with-extra-tables.er",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.providedOptions.OutputFilename = fmt.Sprintf("test-er-diagram-%v", tc.filenameSuffix)
			actualService := service.New(tc.providedOptions, util.New(), writer.New(util.New(), tc.providedOptions.OutputPath, tc.providedOptions.OutputFilename))
			validateGenerateExecution(t, actualService, tc.providedOptions.OutputFilename, tc.expectedOutputFile)
		})
	}
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
	t.Skip()
	c, err := expect.NewTestConsole(t)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	cmd := exec.Command("go", "run", "../../../main.go", "build", "-e")
	cmd.Stdin = c.Tty()
	cmd.Stdout = c.Tty()
	// cmd.Stderr = c.Tty()

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		c.ExpectString("What is the name table name?")
		c.Log("aaa")
		c.Send("custom_table")
		c.SendLine("")
		c.ExpectString("What should be the background color for the table?")
		c.SendLine("#aabbcc")
		fmt.Println("2")
		c.ExpectString("What is the name of the column?")
		c.SendLine("id")
		fmt.Println("3")
		c.ExpectString("Choose the type of the column:")
		c.Send(string(terminal.KeyArrowUp))
		c.Send(" ")
		fmt.Println("4")
		c.ExpectString("Is this field a primary key? (y/N)")
		c.SendLine("y")
		fmt.Println("5")
		c.ExpectString("Is this field a foreign key? (y/N)")
		c.SendLine("N")
		fmt.Println("6")
		c.ExpectString("Do you want to add some more:")
		c.SendLine("Column")
		fmt.Println("7")
		c.ExpectString("What is the name of the column?")
		c.SendLine("version")
		fmt.Println("8")
		c.ExpectString("Choose the type of the column:")
		c.Send(" ")
		fmt.Println("9")
		c.ExpectString("Is this field a primary key? (y/N)")
		c.SendLine("N")
		fmt.Println("10")
		c.ExpectString("Is this field a foreign key? (y/N)")
		c.SendLine("N")
		fmt.Println("11")
		c.ExpectString("Do you want to add some more:")
		c.Send(string(terminal.KeyArrowDown))
		c.Send(" ")
		fmt.Println("12")
		c.ExpectEOF()
	}()

	err = cmd.Wait()
	fmt.Println("aaa")
	if err != nil {
		log.Fatal(err)
	}

	// options := domain.Options{ExtraTablesSurvey: true}
	// srv := service.New(options, util.New(), nil)
	// actualDefinition, actualError := srv.Build()

	// fmt.Println(actualDefinition)
	// fmt.Println(actualError)
}
