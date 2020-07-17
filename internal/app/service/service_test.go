package service_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/Angelos-Giannis/er-builder/internal/app/service"
	"github.com/Angelos-Giannis/er-builder/internal/domain"
	"github.com/udhos/equalfile"
	"github.com/urfave/cli/v2"
)

func TestNew(t *testing.T) {
	options := domain.Options{}
	actualService := service.New(options)

	if reflect.TypeOf(&service.Service{}) != reflect.TypeOf(actualService) {
		t.Errorf("Expected to get type '%v' but got '%v'.", reflect.TypeOf(&service.Service{}), reflect.TypeOf(actualService))
	}
}

func TestGenerate(t *testing.T) {
	options := domain.Options{
		Directory:      "./../../../test/",
		IDField:        "id",
		FileList:       cli.StringSlice{},
		OutputFilename: "test-example-er-diagram",
		OutputPath:     "./../../../test",
		Tag:            "db",
		Title:          "example_db",
	}
	actualService := service.New(options)

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
