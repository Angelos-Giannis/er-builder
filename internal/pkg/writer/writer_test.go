package writer_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/udhos/equalfile"

	"github.com/eujoy/erbuilder/internal/app/service"
	"github.com/eujoy/erbuilder/internal/pkg/util"
	"github.com/eujoy/erbuilder/internal/pkg/writer"
	"github.com/eujoy/erbuilder/test"
)

func TestNewWriter(t *testing.T) {
	actualWriter := writer.New(util.New(), "some/path", "some_filename")

	if reflect.TypeOf(&writer.Writer{}) != reflect.TypeOf(actualWriter) {
		t.Errorf("Expected to get type '%v' but got '%v'.", reflect.TypeOf(&service.Service{}), reflect.TypeOf(actualWriter))
	}
}

func TestWriteFile(t *testing.T) {
	dataBuilder := test.NewDataBuilder()
	diagram := dataBuilder.GetWriterTestDiagram()

	writer := writer.New(util.New(), "./../../../test", "test-writer-example-er-diagram")
	writer.WriteFile(diagram)

	cmp := equalfile.New(nil, equalfile.Options{})
	equal, err := cmp.CompareFile("./../../../test/example-er-diagram.er", "./../../../test/test-writer-example-er-diagram.er")
	if err != nil {
		t.Errorf("Expected to get nil as error but got '%v'.", err)
	}

	if !equal {
		t.Errorf("Expected that the generated file and the example test file are the same but they were not.")
	}

	err = os.Remove("./../../../test/test-writer-example-er-diagram.er")
	if err != nil {
		t.Errorf("Expected to get nil as error when deleting the test file but got '%v'.", err)
	}
}
