package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/jaevor/go-nanoid"
	"github.com/labstack/echo/v4"
	"github.com/urfave/cli/v2"
	"golang.org/x/net/context"
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
		volumeName := strings.ToLower(c.Param("volume"))
		snapshotName := strings.ToLower(c.Param("snapshot"))

		var tarballFilePath string

		if snapshotName == "latest" {
			snapshotName = volumeName + "_" + fmt.Sprintf("%d", time.Now().Unix())
			out, err := exec.CommandContext(
				c.Request().Context(),
				"lvcreate",
				"-s",
				"-n", snapshotName,
				"-l", "%ORIGIN",
				fmt.Sprintf("%s/%s", DEFAULT_VOLUME_GROUP, volumeName),
			).Output()
			if err != nil {
				return err
			}
			log.Println(strings.TrimSpace(string(out)))

			mountPath := filepath.Join(MOUNT_PATH, snapshotName)

			if err := os.MkdirAll(mountPath, 0777); err != nil {
				return err
			}
			defer os.RemoveAll(mountPath)

			out, err = exec.CommandContext(
				c.Request().Context(),
				"mount",
				"/dev/mapper/"+snapshotName,
				mountPath,
			).Output()
			if err != nil {
				return err
			}
			log.Println(strings.TrimSpace(string(out)))

			tarballFilePath = filepath.Join(TARBALL_FILE_PATH, snapshotName+".tar.zst")
			if err := os.RemoveAll(tarballFilePath); err != nil {
				return err
			}
			out, err = exec.CommandContext(
				c.Request().Context(),
				"tar",
				"-I", "'zstd --fast=5'",
				"-cvf", tarballFilePath,
				"-C", mountPath,
				".",
			).Output()
			if err != nil {
				return err
			}
			log.Println(strings.TrimSpace(string(out)))

			out, err = exec.CommandContext(
				c.Request().Context(),
				"umount",
				"-f",
				mountPath,
			).Output()
			if err != nil {
				return err
			}
			log.Println(strings.TrimSpace(string(out)))

			out, err = exec.CommandContext(
				c.Request().Context(),
				"lvremove",
				"-f",
				fmt.Sprintf("%s/%s", DEFAULT_VOLUME_GROUP, volumeName),
			).Output()
			if err != nil {
				return err
			}
			log.Println(strings.TrimSpace(string(out)))
		} else {
			tarballFilePath = filepath.Join(TARBALL_FILE_PATH, snapshotName+".tar.zst")
		}

		if _, err := os.Stat(tarballFilePath); err != nil {
			return err
		}

		return c.File(tarballFilePath)
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

func volumeCreate(ctx context.Context, volumeName string, virtualSizeMB int, remoteSnapshotURL string) error {
	if volumeName == "" {
		volumeID, err := nanoid.CustomASCII("abcdefghijklmnopqrstuvwxyz0123456789", 19)
		if err != nil {
			return err
		}

		volumeName = "vol_" + volumeID()
	}

	out, err := exec.CommandContext(
		ctx,
		"lvcreate",
		"--thinpool", fmt.Sprintf("%s/%s", DEFAULT_VOLUME_GROUP, DEFAULT_THIN_POOL_LV),
		"--name", volumeName,
		"--virtualsize", fmt.Sprintf("%dM", virtualSizeMB),
	).Output()
	if err != nil {
		return err
	}
	log.Println(strings.TrimSpace(string(out)))

	out, err = exec.CommandContext(
		ctx,
		"mkfs.ext4", fmt.Sprintf("/dev/mapper/%s-%s", DEFAULT_VOLUME_GROUP, volumeName),
	).Output()
	if err != nil {
		return err
	}
	log.Println(strings.TrimSpace(string(out)))

	if remoteSnapshotURL != "" {
		tarballFilePath := filepath.Join(TARBALL_FILE_PATH, volumeName+".tar.zst")
		if err := os.RemoveAll(tarballFilePath); err != nil {
			return err
		}

		tarballFile, err := os.Create(tarballFilePath)
		if err != nil {
			return err
		}
		defer tarballFile.Close()

		resp, err := http.Get(remoteSnapshotURL)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		_, err = io.Copy(tarballFile, resp.Body)
		if err != nil {
			return err
		}
	}

	return nil
}

func volumeDelete() error {
	return nil
}
