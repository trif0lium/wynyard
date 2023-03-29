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
							return volumeMount()
						},
					},
					{
						Name: "tree",
						Action: func(cCtx *cli.Context) error {
							return volumeTree()
						},
					},
					{
						Name: "list",
						Action: func(cCtx *cli.Context) error {
							return volumeList()
						},
					},
					{
						Name: "describe",
						Action: func(cCtx *cli.Context) error {
							return volumeDescribe()
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
							return volumeCreate()
						},
					},
					{
						Name: "delete",
						Action: func(cCtx *cli.Context) error {
							return volumeDelete()
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

func volumeMount() error {
	return nil
}

func volumeTree() error {
	return nil
}

func volumeList() error {
	return nil
}

func volumeDescribe() error {
	return nil
}

func volumeCreate() error {
	return nil
}

func volumeDelete() error {
	return nil
}
