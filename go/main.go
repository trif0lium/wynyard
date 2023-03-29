package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "wynyard",
		Commands: []*cli.Command{
			{
				Name: "volume",
				Subcommands: []*cli.Command{
					{
						Name: "api-server",
						Action: func(cCtx *cli.Context) error {
							return nil
						},
					},
					{
						Name: "mount",
						Action: func(cCtx *cli.Context) error {
							return nil
						},
					},
					{
						Name: "tree",
						Action: func(cCtx *cli.Context) error {
							return nil
						},
					},
					{
						Name: "list",
						Action: func(cCtx *cli.Context) error {
							return nil
						},
					},
					{
						Name: "describe",
						Action: func(cCtx *cli.Context) error {
							return nil
						},
					},
					{
						Name: "create",
						Action: func(cCtx *cli.Context) error {
							return nil
						},
					},
					{
						Name: "delete",
						Action: func(cCtx *cli.Context) error {
							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
