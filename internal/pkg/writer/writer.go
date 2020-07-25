package writer

import (
	"fmt"
	"os"
	"sort"

	"github.com/eujoy/erbuilder/internal/domain"
	"github.com/eujoy/erbuilder/internal/pkg/util"
)

// Writer describes the writer package.
type Writer struct {
	outputPath     string
	outputFilename string
	outputFile     *os.File
	util           *util.Util
}

// New creates and returns a new writer instrance.
func New(util *util.Util, outputPath, outputFilename string) *Writer {
	return &Writer{
		outputPath:     outputPath,
		outputFilename: outputFilename,
		outputFile:     nil,
		util:           util,
	}
}

// WriteFile creates and writes the diagram to the desired file.
func (w *Writer) WriteFile(diagram domain.Diagram) error {
	var err error
	w.outputFile, err = os.Create(fmt.Sprintf("%v/%v.er", w.outputPath, w.outputFilename))
	if err != nil {
		return err
	}
	defer w.outputFile.Close()

	if diagram.Title != "" {
		_, err = w.outputFile.WriteString(fmt.Sprintf("title {label: \"%v\"}\n\n", diagram.Title))
		if err != nil {
			return err
		}
	}

	err = w.writeTableDetails(diagram.TableList)
	if err != nil {
		return err
	}

	err = w.writeForeignKeyReferences(diagram.ReferenceList)
	if err != nil {
		return err
	}

	return nil
}

// writeTableDetails writes the table details in the output file.
func (w *Writer) writeTableDetails(tableList []domain.Table) error {
	_, err := w.outputFile.WriteString("# Definition of tables.\n")
	if err != nil {
		return err
	}

	sort.SliceStable(tableList, func(i, j int) bool {
		return tableList[i].Name < tableList[j].Name
	})

	for _, table := range tableList {
		_, err := w.outputFile.WriteString(fmt.Sprintf("[%v]\n", table.Name))
		if err != nil {
			return err
		}

		err = w.writeColumnsOfTable(table.ColumnList)
		if err != nil {
			return err
		}
	}

	return nil
}

// writeColumns writes the column details in the output file.
func (w *Writer) writeColumnsOfTable(columnList []domain.Column) error {
	sort.SliceStable(columnList, func(i, j int) bool {
		if columnList[j].IsPrimaryKey {
			return false
		}

		return w.util.GetCaseOfString(columnList[i].Name, "camelCase") < w.util.GetCaseOfString(columnList[j].Name, "camelCase") ||
			columnList[i].IsPrimaryKey ||
			!columnList[i].IsExtraField
	})

	for _, column := range columnList {
		idPrefix := ""
		fkPrefix := ""
		if column.IsPrimaryKey {
			idPrefix = "*"
		}
		if column.IsForeignKey {
			fkPrefix = "+"
		}

		_, err := w.outputFile.WriteString(
			fmt.Sprintf(
				"\t%v%v%v {label: \"%v\"}\n",
				idPrefix,
				fkPrefix,
				column.Name,
				column.Type,
			),
		)
		if err != nil {
			return err
		}
	}

	_, err := w.outputFile.WriteString("\n")
	if err != nil {
		return err
	}

	return nil
}

// writeForeignKeyReferences writes the foreign key references in the output file.
func (w *Writer) writeForeignKeyReferences(referenceList []domain.Reference) error {
	if len(referenceList) == 0 {
		return nil
	}

	_, err := w.outputFile.WriteString("\n# Definition of foreign keys.\n")
	if err != nil {
		return err
	}

	sort.SliceStable(referenceList, func(i, j int) bool {
		return referenceList[i].FromTableName < referenceList[j].FromTableColumn
	})

	for _, foreignKey := range referenceList {
		_, err := w.outputFile.WriteString(
			fmt.Sprintf(
				"%v %v %v {label: \"%v\"}\n",
				foreignKey.FromTableName,
				foreignKey.TypeOfReference,
				foreignKey.ToTableName,
				foreignKey.FromTableColumn,
			),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
