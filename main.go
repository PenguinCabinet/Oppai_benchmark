package main

import (
	"log"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{},
		Commands: []*cli.Command{
			{
				Name:    "GUI",
				Aliases: []string{"g"},
				Usage:   "run on gui",
				Action: func(c *cli.Context) error {
					GUI_main()
					return nil
				},
			},
			{
				Name:    "CLI",
				Aliases: []string{"c"},
				Usage:   "run on cli",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "gpu",
						Value: false,
						Usage: "GPU",
					},
				},
				Action: func(c *cli.Context) error {
					CLI_main(c.Bool("gpu"))
					return nil
				},
			},
		},
		Action: func(c *cli.Context) error {
			GUI_main()
			return nil
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
