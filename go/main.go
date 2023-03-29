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
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name: "port",
							},
						},
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
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name: "remote-snapshot",
							},
							&cli.IntFlag{
								Name: "size",
							},
						},
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
