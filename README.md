# er-builder

![version](https://img.shields.io/badge/version-v0.1.0-brightgreen)
![golang-version](https://img.shields.io/badge/Go-1.14-blue)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![master-actions Actions Status](https://github.com/Angelos-Giannis/er-builder/workflows/master-actions/badge.svg)](https://github.com/Angelos-Giannis/er-builder/actions)

The purpose of this tool is to parse a file or a list of files containing the mappping against the database and generate an `.er` file describing the database modeling.

## Installation

In order to be able to use this you simply have to :

```shell
go get github.com/Angelos-Giannis/er-builder
```

## Definition of commands

There is a command called `generate` which actually retrieves the details from the provided files or files under directory and generates the respective .er file based on the definition of the structs in those files.

```shell
NAME:
   main generate - Generate the .er file based on the provided structures.

USAGE:
   main generate [command options] [arguments...]

OPTIONS:
   --common_field value, -c value       Common field for all the tables which do not have the provided tag in place.
   --directory value, -d value          Directory to retrieve the files from.
   --file_list value, -l value          List of files to parse.
   --id_field value                     Id field to be used for all the tables.
   --output_filename value, --of value  Define the generated output filename (will be used for both the .er and the image file). (default: "er-diagram")
   --output_path value, -o value        The path were to store the .er file. (default: ".")
   --tag value, -t value                Tag value to consume from the structs. (default: "db")
   --title value                        Title to be included in the exported image. (default: "Database Schema")
   --help, -h                           show help (default: false)
```

## How to use

To generate a new `.er` file with the database models, a command like the following needs to be executed :

```shell
go run cmd/er-builder/main.go generate --directory "./test/" --output_path "./test/" --output_filename "example-er-diagram" --id_field id --tag db --title example_db
```
