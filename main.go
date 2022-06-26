package main

import (
	"github.com/mg98/scriptup/pkg/scriptup"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := NewApp()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// NewApp creates a new command line app
func NewApp() *cli.App {
	app := cli.NewApp()
	app.Name = "scriptup"
	app.Usage = "A lightweight and agnostic script migration tool."
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "env",
			Aliases: []string{"e"},
			Value:   "dev",
			Usage:   "specify which configuration to use",
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:    "new",
			Aliases: []string{"n"},
			Usage:   "Generate a new migration file",
			Action: func(c *cli.Context) error {
				cfg := scriptup.GetConfig(c.String("env"))
				return scriptup.NewMigrationFile(cfg, c.Args().Get(0))
			},
		},
		{
			Name:    "up",
			Aliases: []string{"u"},
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:  "steps",
					Value: -1,
					Usage: "Limit the number of migrations to run",
				},
			},
			Usage: "Execute recent scripts that have not been migrated yet",
			Action: func(c *cli.Context) error {
				cfg := scriptup.GetConfig(c.String("env"))
				return scriptup.MigrateUp(cfg, c.Int("steps"))
			},
		},
		{
			Name:    "down",
			Aliases: []string{"d"},
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:  "steps",
					Value: -1,
					Usage: "Limit the number of migrations to run",
				},
			},
			Usage: "Undo recently performed migrations",
			Action: func(c *cli.Context) error {
				cfg := scriptup.GetConfig(c.String("env"))
				return scriptup.MigrateDown(cfg, c.Int("steps"))
			},
		},
		{
			Name:    "status",
			Aliases: []string{"s"},
			Usage:   "Get status about open migrations and which was run last",
			Action: func(c *cli.Context) error {
				cfg := scriptup.GetConfig(c.String("env"))
				return scriptup.Status(cfg)
			},
		},
	}

	return app
}
