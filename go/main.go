package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/jaevor/go-nanoid"
	"github.com/labstack/echo/v4"
	"github.com/urfave/cli/v2"
)

const (
	DEFAULT_VOLUME_GROUP = "vg0"
	DEFAULT_THIN_POOL_LV = "lv0"
	TARBALL_FILE_PATH    = "/tmp/tarballs"
	MOUNT_PATH           = "/mnt"
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
							volumeName := cCtx.Args().First()
							virtualSizeMB := cCtx.Int("size")
							return volumeCreate(volumeName, virtualSizeMB)
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

func volumeAPIServer(port int) error {
	if err := os.MkdirAll(TARBALL_FILE_PATH, 0777); err != nil {
		return err
	}

	if err := os.MkdirAll(MOUNT_PATH, 0777); err != nil {
		return err
	}

	e := echo.New()

	e.GET("/volumes/:volume/snapshots/:snapshot", func(c echo.Context) error {
	})

	if err := e.Start(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
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

func volumeCreate(volumeName string, virtualSizeMB int) error {
	if volumeName == "" {
		volumeID, err := nanoid.CustomASCII("abcdefghijklmnopqrstuvwxyz0123456789", 19)
		if err != nil {
			return err
		}

		volumeName = "vol_" + volumeID()
	}

	out, err := exec.Command(
		"lvcreate",
		"--thinpool", fmt.Sprintf("%s/%s", DEFAULT_VOLUME_GROUP, DEFAULT_THIN_POOL_LV),
		"--name", volumeName,
		"--virtualsize", fmt.Sprintf("%dM", virtualSizeMB),
	).Output()
	if err != nil {
		return err
	}
	log.Println(strings.TrimSpace(string(out)))

	out, err = exec.Command(
		"mkfs.ext4", fmt.Sprintf("/dev/mapper/%s-%s", DEFAULT_VOLUME_GROUP, volumeName),
	).Output()
	if err != nil {
		return err
	}
	log.Println(strings.TrimSpace(string(out)))

	return nil
}

func volumeDelete() error {
	return nil
}
