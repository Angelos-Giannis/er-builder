/*
Package erbuilder reads across a list of files and after retrieving the structs found in there, it generates a '.er' file based on a specific tag.

Installation can be done by running :

	go get -u github.com/eujoy/erbuilder

To use the tool and generate a .er file, the following command needs to be executed (change values accordingly) :

	erbuilder generate --directory "./test/" --output_path "./test/" --output_filename "example-er-diagram" --id_field "id" --tag "db" --title "example_db"

Also, this can be used as an external package in any service. To use this tool as part of your project, you need to do something like the following :

	import (
		"github.com/eujoy/erbuilder/internal/app/service"
		"github.com/eujoy/erbuilder/internal/domain"
	)

	// Config describes the configuration of the service.
	type Config struct {
		Settings settings
	}

	type settings struct {
		AllowedColumnNameCaseValues []string
		AllowedTableNameCaseValues  []string
	}

	func main() {
		cfg := Config{
			Settings: settings{
				AllowedColumnNameCaseValues: []string{"snake_case", "camelCase", "screaming_snake_case", "kebab_case"},
				AllowedTableNameCaseValues:  []string{"snake_case", "camelCase", "screaming_snake_case", "kebab_case"},
			},
		}

		// Specify your values for the options.
		options := domain.Options{
			Directory:      "./../../../test/",
			IDField:        "id",
			FileList:       cli.StringSlice{},
			OutputFilename: "test-example-er-diagram",
			OutputPath:     "./../../../test",
			Tag:            "db",
			Title:          "example_db",
			ColumnNameCase: "snake_case",
			TableNameCase:  "snake_case",
			Config:         cfg,
		}
		actualService := service.New(options)
		_ = actualService.Generate()
	}

For more details, please visit : https://github.com/eujoy/erbuilder.
*/
package main
