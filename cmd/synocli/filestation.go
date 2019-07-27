package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

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

	_, err := syno.FileStation.Stat(path)
	if err != nil {
		log.Panic(err)
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
