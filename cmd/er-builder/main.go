package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Angelos-Giannis/erd-builder/internal/app/service"
	"github.com/Angelos-Giannis/erd-builder/internal/config"
	"github.com/Angelos-Giannis/erd-builder/internal/domain"
	"github.com/urfave/cli/v2"
)

const (
	// define the configuration file for the system.
	configurationFile = "configuration.yaml"
)

func main() {
	cfg, err := config.New(configurationFile)
	if err != nil {
		fmt.Printf("Error parsing configuration : %v\n", err)
		os.Exit(1)
	}

	var app = cli.NewApp()
	info(app, cfg)

	var options domain.Options

	app.Commands = []*cli.Command{
		&cli.Command{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "Generate the .er file based on the provided structures.",
			Flags: []cli.Flag{
				options.GetCommonField(),
				options.GetDirectoryFlag(),
				options.GetFileList(),
				options.GetIDField(),
				options.GetOutputFilename(),
				options.GetOutputPath(),
				options.GetTag(),
				options.GetTitle(),
			},
			Action: func(c *cli.Context) error {
				srv := service.New(options)
				return srv.Generate()
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// info sets up the information of the tool.
func info(app *cli.App, cfg config.Config) {
	var appAuthors []*cli.Author
	for _, author := range cfg.Application.Authors {
		newAuthor := cli.Author{
			Name:  author.Name,
			Email: author.Email,
		}
		appAuthors = append(appAuthors, &newAuthor)
	}

	app.Authors = appAuthors
	app.Name = cfg.Application.Name
	app.Usage = cfg.Application.Usage
	app.Version = cfg.Application.Version
}
