package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/urfave/cli"
)

var apiCommands = []cli.Command{
	{
		Name:        "list",
		ShortName:   "ls",
		Usage:       "List APIs reported",
		Description: "List APIs reported",
		Action:      apiList,
		Flags:       []cli.Flag{},
	},
}

func apiList(ctx *cli.Context) {
	syno := mustGetSyno(ctx)

	api, err := syno.API()
	if err != nil {
		log.Panic(err)
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, ' ', 0)

	fmt.Fprintf(w, "API\tMin\tMax\tPath\t\n")

	for k, v := range api {
		fmt.Fprintf(w, "%s\t%d\t%d\t%s\t\n", k, v.MinVersion, v.MaxVersion, v.Path)

	}
	w.Flush()
}
