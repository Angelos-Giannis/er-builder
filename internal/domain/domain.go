package domain

import (
	"errors"
	"fmt"

	"github.com/Angelos-Giannis/erbuilder/internal/config"

	"github.com/urfave/cli/v2"
)

// Options describe the allowed options of the cli tool.
type Options struct {
	CommonFields    cli.StringSlice
	Directory       string
	FileList        cli.StringSlice
	IDField         string
	OutputFilename  string
	OutputPath      string
	Tag             string
	Title           string
	ColumnNameCase  string
	TableNameCase   string
	TableNamePlural bool

	Config config.Config
}

// NewOptions creates and returns a new options structure.
func NewOptions(cfg config.Config) Options {
	return Options{
		Config: cfg,
	}
}

// Validate the provided values to confirm that they are all correct.
func (o *Options) Validate() error {
	if o.Directory == "" && len(o.FileList.Value()) == 0 {
		return errors.New("Need to provide at least one of 'directory' or 'file_list'")
	}

	if !o.validateWithAllowedValues(o.ColumnNameCase, o.Config.Settings.AllowedColumnNameCaseValues) {
		return fmt.Errorf(
			"The provided value for column name case is not valid. Allowed values : %v",
			o.Config.Settings.AllowedColumnNameCaseValues,
		)
	}

	if !o.validateWithAllowedValues(o.TableNameCase, o.Config.Settings.AllowedTableNameCaseValues) {
		return fmt.Errorf(
			"The provided value for table name case is not valid. Allowed values : %v",
			o.Config.Settings.AllowedTableNameCaseValues,
		)
	}

	return nil
}

// validateWithAllowedValues checks if the provided string value of a field is in the list of allowed ones.
func (o *Options) validateWithAllowedValues(providedValue string, allowedValues []string) bool {
	for _, allowed := range allowedValues {
		if providedValue == allowed {
			return true
		}
	}
	return false
}

// GetCommonFields returns the definition for common_field flag.
func (o *Options) GetCommonFields() *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Name:        "common_field",
		Aliases:     []string{"c"},
		Usage:       "Common field for all the tables which do not have the provided tag in place.",
		Value:       nil,
		Destination: &o.CommonFields,
		Required:    false,
	}
}

// GetDirectoryFlag returns the definition for directory flag.
func (o *Options) GetDirectoryFlag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "directory",
		Aliases:     []string{"d"},
		Usage:       "Directory to retrieve the files from.",
		Value:       "",
		Destination: &o.Directory,
		Required:    false,
	}
}

// GetFileList returns the definition for file_list flag.
func (o *Options) GetFileList() *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Name:        "file_list",
		Aliases:     []string{"l"},
		Usage:       "List of files to parse.",
		Value:       nil,
		Destination: &o.FileList,
		Required:    false,
	}
}

// GetIDField returns the definition for id_field flag.
func (o *Options) GetIDField() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "id_field",
		Usage:       "Id field to be used for all the tables.",
		Value:       "",
		Destination: &o.IDField,
		Required:    false,
	}
}

// GetOutputFilename returns the definition for output_filename flag.
func (o *Options) GetOutputFilename() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "output_filename",
		Aliases:     []string{"of"},
		Usage:       "Define the generated output filename (will be used for both the .er and the image file).",
		Value:       "er-diagram",
		Destination: &o.OutputFilename,
		Required:    false,
	}
}

// GetOutputPath returns the definition for output_path flag.
func (o *Options) GetOutputPath() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "output_path",
		Aliases:     []string{"o"},
		Usage:       "The path were to store the .er file.",
		Value:       ".",
		Destination: &o.OutputPath,
		Required:    false,
	}
}

// GetTag returns the definition for tag flag.
func (o *Options) GetTag() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "tag",
		Aliases:     []string{"t"},
		Usage:       "Tag value to consume from the structs.",
		Value:       "db",
		Destination: &o.Tag,
		Required:    false,
	}
}

// GetTitle returns the definition for title flag.
func (o *Options) GetTitle() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "title",
		Usage:       "Title to be included in the exported image.",
		Value:       "Database Schema",
		Destination: &o.Title,
		Required:    false,
	}
}

// GetColumnNameCase returns the definition for column_name_case flag.
func (o *Options) GetColumnNameCase() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "column_name_case",
		Aliases:     []string{"cnc"},
		Usage:       fmt.Sprintf("Define the case definition for the column names. (Allowed values : %v)", o.Config.Settings.AllowedTableNameCaseValues),
		Value:       "snake_case",
		Destination: &o.ColumnNameCase,
		Required:    false,
	}
}

// GetTableNameCase returns the definition for table_name_case flag.
func (o *Options) GetTableNameCase() *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "table_name_case",
		Aliases:     []string{"tnc"},
		Usage:       fmt.Sprintf("Define the case definition for the table names. (Allowed values : %v)", o.Config.Settings.AllowedTableNameCaseValues),
		Value:       "snake_case",
		Destination: &o.TableNameCase,
		Required:    false,
	}
}

// GetTableNamePlural returns the definition for title flag.
func (o *Options) GetTableNamePlural() *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:        "table_in_plural",
		Aliases:     []string{"tp"},
		Usage:       "Define whether the table name should be in plural.",
		Value:       false,
		Destination: &o.TableNamePlural,
		Required:    false,
	}
}
