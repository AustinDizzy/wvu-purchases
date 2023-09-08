package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "wvupurchases",
		Usage: "A command line tool for aggregating and reporting on purchase data from West Virginia University",
		Commands: []*cli.Command{
			{
				Name:   "ingest",
				Usage:  "Ingest a purchase file (.xlsx) to the database",
				Action: ingest,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "overwrite",
						Usage: "overwrite existing data in the database",
						Value: false,
					},
					&cli.StringFlag{
						Name:     "type",
						Usage:    "the type of purchase record to ingest, must be either pcard or procurement",
						Required: true,
					},
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "db",
				Usage:    "the database file (.db)",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "enable verbose logging",
				Value: false,
			},
		},
	}

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
