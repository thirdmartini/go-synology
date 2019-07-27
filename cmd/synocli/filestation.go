package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/urfave/cli"

	"github.com/thirdmartini/go-synology"
)

var fileStationCommands = []cli.Command{
	{
		Name:        "list",
		ShortName:   "ls",
		Usage:       "List devices",
		Description: "Lists all devices",
		Action:      fileStationList,
		Flags: []cli.Flag{
			pathFlag,
		},
	},
	{
		Name:        "stat",
		ShortName:   "stat",
		Usage:       "Stat a specific path",
		Description: "Stat a specific path",
		Action:      fileStationStat,
		Flags: []cli.Flag{
			pathFlag,
		},
	},
	{
		Name:        "download",
		ShortName:   "dl",
		Usage:       "Download a file",
		Description: "Download a file",
		Action:      fileStationDownload,
		Flags: []cli.Flag{
			pathFlag,
			destFlag,
		},
	},
	{
		Name:        "md5",
		ShortName:   "md5",
		Usage:       "Calculate the MD5 of a file on the NAS",
		Description: "Calculate the MD5 of a file on the NAS",
		Action:      fileStationMD5,
		Flags: []cli.Flag{
			pathFlag,
		},
	},
}

func fileStationList(ctx *cli.Context) {
	syno := mustGetSyno(ctx)

	path := ctx.String(pathFlag.Name)

	var files []synology.FileInfo
	var err error
	if path == "" || path == "/" {
		files, err = syno.FileStation.ListShares()
	} else {
		files, err = syno.FileStation.List(path)
	}
	if err != nil {
		log.Panic(err)
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, ' ', 0)

	fmt.Fprintf(w, "  \tOwner\tGroup\tSize\tName\t\n")
	for _, fi := range files {
		flags := " "
		if fi.Isdir {
			flags = "DIR"
		}
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t\n",
			flags,
			fi.Stat.Owner.User,
			fi.Stat.Owner.Group,
			humanize.Bytes(fi.Stat.Size),
			fi.Name)
	}
	w.Flush()
}

func fileStationStat(ctx *cli.Context) {
	syno := mustGetSyno(ctx)

	path := ctx.String(pathFlag.Name)

	files, err := syno.FileStation.Stat(path)
	if err != nil {
		log.Panic(err)
	}

	for _, fi := range files {
		fmt.Printf("File: %s\n", fi.Path)

		kind := "Regular File"
		if fi.Isdir {
			kind = "Directory"
		}
		fmt.Printf("Size: %10s  File Type: %s\n", humanize.Bytes(fi.Stat.Size), kind)
		fmt.Printf("Owner: %10s  Group: %s\n", fi.Stat.Owner.User, fi.Stat.Owner.Group)
		fmt.Printf("Access: %s\n", time.Unix(int64(fi.Stat.Time.Atime), 0).String())
		fmt.Printf("Modify: %s\n", time.Unix(int64(fi.Stat.Time.Mtime), 0).String())
		fmt.Printf("Change: %s\n", time.Unix(int64(fi.Stat.Time.Ctime), 0).String())
	}

}

func fileStationDownload(ctx *cli.Context) {
	syno := mustGetSyno(ctx)

	path := ctx.String(pathFlag.Name)

	err := syno.FileStation.Download(path, os.Stdout)
	if err != nil {
		log.Panic(err)
	}
}

func fileStationMD5(ctx *cli.Context) {
	syno := mustGetSyno(ctx)

	path := ctx.String(pathFlag.Name)

	hash, err := syno.FileStation.MD5(path)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("%s %s\n", hash, path)
}
