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
	"go.uber.org/zap"
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
				Name: "debug",
				Action: func(cCtx *cli.Context) error {
					hostname, err := os.Hostname()
					if err != nil {
						return err
					}
					fmt.Println(hostname)
					return nil
				},
			},
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
							port := cCtx.Int("port")
							return volumeAPIServer(port)
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
								Name: "snapshot-host",
							},
							&cli.StringFlag{
								Name: "snapshot-location",
							},
							&cli.IntFlag{
								Name: "size",
							},
						},
						Action: func(cCtx *cli.Context) error {
							volumeName := cCtx.Args().First()
							virtualSizeMB := cCtx.Int("size")
							snapshotHost := cCtx.String("snapshot-host")
							snapshotLocation := cCtx.String("snapshot-location")
							remoteSnapshotURL := ""
							if snapshotHost != "" && snapshotLocation != "" {
								remoteSnapshotURL = fmt.Sprintf("http://%s.%s.c.railway-infra-dev.internal:1323/volumes/%s/snapshots/latest", snapshotHost, snapshotLocation, volumeName)
							}
							return volumeCreate(cCtx.Context, volumeName, virtualSizeMB, remoteSnapshotURL)
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

	os.Exit(0)
}

func volumeAPIServer(port int) error {
	zapConfig := zap.NewProductionConfig()
	zapConfig.OutputPaths = []string{"stdout"}
	logger, err := zapConfig.Build()
	if err != nil {
		return err
	}
	defer logger.Sync()

	if err := os.MkdirAll(TARBALL_FILE_PATH, 0777); err != nil {
		return err
	}

	if err := os.MkdirAll(MOUNT_PATH, 0777); err != nil {
		return err
	}

	e := echo.New()

	e.GET("/volumes/:volume/stream", func(c echo.Context) error {
		volumeName := strings.ToLower(c.Param("volume"))

		r, w := io.Pipe()
		defer r.Close()

		r2, w2 := io.Pipe()
		defer r2.Close()

		cmd := exec.CommandContext(
			c.Request().Context(),
			"dd",
			fmt.Sprintf("if=/dev/%s/%s", DEFAULT_VOLUME_GROUP, volumeName),
			"bs=8M",
		)
		cmd.Stdout = w

		zstdCmd := exec.CommandContext(
			c.Request().Context(),
			"zstd",
			"-5",
			"-",
		)
		zstdCmd.Stdin = r
		zstdCmd.Stdout = w2

		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEOctetStream)
		c.Response().WriteHeader(http.StatusOK)

		go func() {
			if err := zstdCmd.Start(); err != nil {
				logger.Sugar().Error(err)
				return
			}

			if err := cmd.Start(); err != nil {
				logger.Sugar().Error(err)
				return
			}

			if err := cmd.Wait(); err != nil {
				logger.Sugar().Error(err)
				return
			}

			w.Close()

			if err := zstdCmd.Wait(); err != nil {
				logger.Sugar().Error(err)
				return
			}

			w2.Close()
		}()

		_, err = io.Copy(c.Response(), r2)
		if err != nil {
			logger.Sugar().Error(err)
			return err
		}

		return nil
	})

	e.GET("/volumes/:volume/snapshots/:snapshot", func(c echo.Context) error {
		volumeName := strings.ToLower(c.Param("volume"))
		snapshotName := strings.ToLower(c.Param("snapshot"))

		var tarballFilePath string

		if snapshotName == "latest" {
			snapshotName = volumeName + "_" + fmt.Sprintf("%d", time.Now().Unix())
			cmd := exec.CommandContext(
				c.Request().Context(),
				"lvcreate",
				"-s",
				"-n", snapshotName,
				"-l", "100%ORIGIN",
				fmt.Sprintf("%s/%s", DEFAULT_VOLUME_GROUP, volumeName),
			)
			out, err := cmd.CombinedOutput()
			if err != nil {
				logger.Sugar().Errorln(strings.Join(cmd.Args, " "))
				logger.Sugar().Error(string(out))
				logger.Sugar().Error(err)
				return err
			}
			log.Println(strings.TrimSpace(string(out)))

			mountPath := filepath.Join(MOUNT_PATH, snapshotName)

			if err := os.MkdirAll(mountPath, 0777); err != nil {
				logger.Sugar().Error(err)
				return err
			}
			defer os.RemoveAll(mountPath)

			out, err = exec.CommandContext(
				c.Request().Context(),
				"mount",
				"/dev/mapper/"+fmt.Sprintf("%s-%s", DEFAULT_VOLUME_GROUP, snapshotName),
				mountPath,
			).Output()
			if err != nil {
				logger.Sugar().Error(err)
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
				"-I", "zstd --fast=5",
				"-cvf", tarballFilePath,
				"-C", mountPath,
				".",
			).Output()
			if err != nil {
				logger.Sugar().Error(string(out))
				logger.Sugar().Error(err)
				return err
			}
			log.Println(strings.TrimSpace(string(out)))

			out, err = exec.CommandContext(
				c.Request().Context(),
				"umount",
				"-f",
				"/dev/mapper/"+fmt.Sprintf("%s-%s", DEFAULT_VOLUME_GROUP, snapshotName),
			).Output()
			if err != nil {
				logger.Sugar().Error(err)
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
				logger.Sugar().Error(err)
				return err
			}
			log.Println(strings.TrimSpace(string(out)))
		} else {
			tarballFilePath = filepath.Join(TARBALL_FILE_PATH, snapshotName+".tar.zst")
		}

		if _, err := os.Stat(tarballFilePath); err != nil {
			logger.Sugar().Error(err)
			return err
		}

		return c.File(tarballFilePath)
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
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
	zapConfig := zap.NewProductionConfig()
	zapConfig.OutputPaths = []string{"stdout"}
	logger, err := zapConfig.Build()
	if err != nil {
		return err
	}
	defer logger.Sync()

	hostname, _ := os.Hostname()
	logger = logger.With(zap.String("hostname", hostname))

	if volumeName == "" {
		volumeID, err := nanoid.CustomASCII("abcdefghijklmnopqrstuvwxyz0123456789", 19)
		if err != nil {
			return err
		}

		volumeName = "vol_" + volumeID()
	}

	cmd := exec.CommandContext(
		ctx,
		"lvcreate",
		"--thinpool", fmt.Sprintf("%s/%s", DEFAULT_VOLUME_GROUP, DEFAULT_THIN_POOL_LV),
		"--name", volumeName,
		"--virtualsize", fmt.Sprintf("%dM", virtualSizeMB),
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		logger.Sugar().Errorln(strings.Join(cmd.Args, " "))
		logger.Sugar().Error(string(out))
		logger.Sugar().Error(err)
		return err
	}
	log.Println(strings.TrimSpace(string(out)))

	out, err = exec.CommandContext(
		ctx,
		"mkfs.ext4", fmt.Sprintf("/dev/mapper/%s-%s", DEFAULT_VOLUME_GROUP, volumeName),
	).Output()
	if err != nil {
		logger.Sugar().Error(err)
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

		mountPath := filepath.Join(MOUNT_PATH, volumeName)

		if err := os.MkdirAll(mountPath, 0777); err != nil {
			return err
		}
		defer os.RemoveAll(mountPath)

		cmd := exec.CommandContext(
			ctx,
			"mount",
			"/dev/mapper/"+fmt.Sprintf("%s-%s", DEFAULT_VOLUME_GROUP, volumeName),
			mountPath,
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			logger.Sugar().Error(strings.Join(cmd.Args, " "))
			logger.Sugar().Error(string(out))
			logger.Sugar().Error(err)
			return err
		}
		log.Println(strings.TrimSpace(string(out)))

		out, err = exec.CommandContext(
			ctx,
			"tar",
			"-I", "zstd",
			"-xvf",
			tarballFilePath,
			"-C",
			mountPath,
		).CombinedOutput()
		if err != nil {
			logger.Sugar().Error(string(out))
			logger.Sugar().Error(err)
			return err
		}
		log.Println(strings.TrimSpace(string(out)))

		out, err = exec.CommandContext(
			ctx,
			"umount",
			"-f",
			"/dev/mapper/"+fmt.Sprintf("%s-%s", DEFAULT_VOLUME_GROUP, volumeName),
		).CombinedOutput()
		if err != nil {
			logger.Sugar().Error(string(out))
			logger.Sugar().Error(err)
			return err
		}
		log.Println(strings.TrimSpace(string(out)))

		if err := os.RemoveAll(tarballFilePath); err != nil {
			return err
		}
	}

	return nil
}

func volumeDelete() error {
	return nil
}
