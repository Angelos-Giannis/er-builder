package domain_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Angelos-Giannis/er-builder/internal/domain"
	"github.com/urfave/cli/v2"
)

func TestGetCommonField(t *testing.T) {
	expectedOptionName := "common_field"

	options := &domain.Options{}
	actualFlag := options.GetCommonField()

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
