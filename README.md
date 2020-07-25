# erbuilder

[![GoDoc](https://godoc.org/github.com/eujoy/erbuilder?status.png)](https://pkg.go.dev/github.com/eujoy/erbuilder)
![version](https://img.shields.io/badge/version-v0.5.0-brightgreen)
![golang-version](https://img.shields.io/badge/Go-1.14-blue)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

----

[![Go Report Card](https://goreportcard.com/badge/github.com/eujoy/erbuilder)](https://goreportcard.com/report/github.com/eujoy/erbuilder)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/0c302dec02df4805ba27c9eda331ef98)](https://www.codacy.com/manual/eujoy/erbuilder?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=eujoy/erbuilder&amp;utm_campaign=Badge_Grade)
[![master-actions Actions Status](https://github.com/eujoy/erbuilder/workflows/master-actions/badge.svg)](https://github.com/eujoy/erbuilder/actions)
![.github/workflows/release-actions.yaml](https://github.com/eujoy/erbuilder/workflows/.github/workflows/release-actions.yaml/badge.svg)

The purpose of this tool is to parse a file or a list of files containing the mappping against the database and generate an `.er` file describing the database modeling.

## Installation

### Installing via go get

In order to be able to use this you simply have to :

```shell
go get -u github.com/eujoy/erbuilder
```

### Installing via brew

In order to install this tool via brew you simply have to :

```shell
brew tap eujoy/erbuilder
brew install erbuilder
```

## Definition of commands

There is a command called `generate` which actually retrieves the details from the provided files or files under directory and generates the respective .er file based on the definition of the structs in those files.

```shell
NAME:
   main generate - Generate the .er file based on the provided structures.

USAGE:
   main generate [command options] [arguments...]

OPTIONS:
   --common_field value, -c value         Common field for all the tables which do not have the provided tag in place.
   --directory value, -d value            Directory to retrieve the files from.
   --file_list value, -l value            List of files to parse.
   --id_field value                       Id field to be used for all the tables.
   --output_filename value, --of value    Define the generated output filename (will be used for both the .er and the image file). (default: "er-diagram")
   --output_path value, -o value          The path were to store the .er file. (default: ".")
   --tag value, -t value                  Tag value to consume from the structs. (default: "db")
   --title value                          Title to be included in the exported image. (default: "Database Schema")
   --column_name_case value, --cnc value  Define the case definition for the column names. (Allowed values : [snake_case camelCase screaming_snake_case kebab_case]) (default: "snake_case")
   --table_name_case value, --tnc value   Define the case definition for the table names. (Allowed values : [snake_case camelCase screaming_snake_case kebab_case]) (default: "snake_case")
   --table_in_plural, --tp                Define whether the table name should be in plural. (default: false)
   --help, -h                             show help (default: false)
```

## How to use

To generate a new `.er` file with the database models, a command like the following needs to be executed :

```shell
erbuilder generate --directory "./test/" --output_path "./test/" --output_filename "example-er-diagram" --id_field "id" --tag "db" --title "example_db"
```

or by using `go run` in the root folder of the project :

```shell
go run main.go generate --directory "./test/" --output_path "./test/" --output_filename "example-er-diagram" --id_field "id" --tag "db" --title "example_db"
```
