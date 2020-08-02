package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/eujoy/erbuilder/internal/app/service"
	"github.com/eujoy/erbuilder/internal/config"
	"github.com/eujoy/erbuilder/internal/domain"
	"github.com/eujoy/erbuilder/internal/pkg/util"
	"github.com/eujoy/erbuilder/internal/pkg/writer"
	"github.com/urfave/cli/v2"
)

// Main actual main function.
func Main() {
	cfg := config.New()

	var app = cli.NewApp()
	info(app, cfg)

	options := domain.NewOptions(cfg)

	app.Commands = []*cli.Command{
		{
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "Generate the .er file based on the provided structures.",
			Flags: []cli.Flag{
				options.GetCommonFields(),
				options.GetDirectoryFlag(),
				options.GetExtraTablesDefinition(),
				options.GetExtraTablesSurvey(),
				options.GetFileList(),
				options.GetIDField(),
				options.GetOutputFilename(),
				options.GetOutputPath(),
				options.GetTag(),
				options.GetTitle(),
				options.GetColumnNameCase(),
				options.GetTableNameCase(),
				options.GetTableNamePlural(),
			},
			Action: func(c *cli.Context) error {
				err := options.Validate()
				if err != nil {
					panic(err)
				}

				util := util.New()
				writer := writer.New(util, options.OutputPath, options.OutputFilename)

				srv := service.New(options, util, writer)
				return srv.Generate()
			},
		},
		{
			Name:    "build",
			Aliases: []string{"b"},
			Usage:   "Use this command to create a json object with extra table definition.",
			Flags: []cli.Flag{
				options.GetExtraTablesSurvey(),
			},
			Action: func(c *cli.Context) error {
				if !options.ExtraTablesSurvey {
					err := errors.New("undefined required flag to run the builder")
					return err
				}

				srv := service.New(options, nil, nil)
				extraDefinition, err := srv.Build()
				if err != nil {
					return err
				}

				byteDefinition, err := json.Marshal(extraDefinition)
				if err != nil {
					return err
				}
				fmt.Println(string(byteDefinition))

				return nil
			},
		},
	}

	err := app.Run(os.Args)
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
