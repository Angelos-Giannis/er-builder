package domain

import "github.com/urfave/cli/v2"

// Options describe the allowed options of the cli tool.
type Options struct {
	CommonFields   cli.StringSlice
	Directory      string
	FileList       cli.StringSlice
	IDField        string
	OutputFilename string
	OutputPath     string
	Tag            string
	Title          string
}

// GetCommonField returns the definition for common_field flag
func (o *Options) GetCommonField() *cli.StringSliceFlag {
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
