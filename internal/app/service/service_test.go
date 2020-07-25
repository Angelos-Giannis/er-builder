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
	fileListStringSlice := cli.StringSlice{}
	fileListStringSlice.Set("./../../../test/example.go")

	options := domain.Options{
		Directory:      "",
		FileList:       fileListStringSlice,
		IDField:        "id",
		OutputFilename: "test-example-er-diagram",
		OutputPath:     "./../../../test",
		Tag:            "db",
		Title:          "example_db",
		ColumnNameCase: "snake_case",
		TableNameCase:  "snake_case",
		Config:         config.New(),
	}
	actualService := service.New(options, util.New(), writer.New(util.New(), options.OutputPath, options.OutputFilename))

	err := actualService.Generate()
	if err != nil {
		t.Errorf("Expected to get nil as error but got '%v'.", err)
	}

	cmp := equalfile.New(nil, equalfile.Options{})
	equal, err := cmp.CompareFile("./../../../test/example-er-diagram.er", "./../../../test/test-example-er-diagram.er")
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
