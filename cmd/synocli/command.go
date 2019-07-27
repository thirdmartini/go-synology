package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/urfave/cli"

	"github.com/thirdmartini/go-synology"
)

const (
	FormatHuman = "human"
	FormatJSON  = "json"
)

var (
	addressFlag = cli.StringFlag{
		Name:  "address",
		Value: "",
		Usage: "Address of the Synology NAS",
	}

	userFlag = cli.StringFlag{
		Name:  "auth.user",
		Value: "",
		Usage: "Username to login",
	}

	passwordFlag = cli.StringFlag{
		Name:  "auth.password",
		Value: "",
		Usage: "Password",
	}

	debugEnabledFlag = cli.BoolFlag{
		Name:  "debug.enabled",
		Usage: "Enable debug logging",
	}

	pathFlag = cli.StringFlag{
		Name:  "path",
		Value: "",
		Usage: "Path to list",
	}

	formatFlag = cli.StringFlag{
		Name:  "format",
		Value: "human",
		Usage: "Output display format (json|human)",
	}

	destFlag = cli.StringFlag{
		Name:  "dest",
		Value: "",
		Usage: "Destination path",
	}
)

var globalFlags = []cli.Flag{
	addressFlag,
	userFlag,
	passwordFlag,
	debugEnabledFlag,
}

// Commands are the top-level commands for rachio-cli
var commands = []cli.Command{
	{
		Name:        "api",
		ShortName:   "api",
		Usage:       "API info Commands",
		Description: "API info Commands",
		Subcommands: apiCommands,
	},

	{
		Name:        "filestation",
		ShortName:   "fs",
		Usage:       "FileStation commands",
		Description: "FileStation commands",
		Subcommands: fileStationCommands,
	},
}

func mustGetSyno(ctx *cli.Context) *synology.Synology {
	addr := ctx.GlobalString(addressFlag.Name)
	user := ctx.GlobalString(userFlag.Name)
	pass := ctx.GlobalString(passwordFlag.Name)

	if addr == "" {
		log.Panic("--address required. (EI: https://10.0.0.10:5001)")
	}

	if user == "" {
		log.Panic("--auth.user required")
	}

	if pass == "" {
		log.Panic("--auth.password required")
	}

	syno, err := synology.Login(addr, user, pass)
	if err != nil {
		log.Panic(err)
	}

	if ctx.GlobalBool(debugEnabledFlag.Name) {
		logger := &log.Logger{}
		logger.SetOutput(os.Stdout)

		syno.WithLogger(logger)
	}
	return syno
}

func mustGetDisplayFormat(ctx *cli.Context) string {
	switch ctx.GlobalString(formatFlag.Name) {
	case "json":
		return "json"
	case "human":
		return "human"
	}

	return "human"
}

func prettyPrintJSON(obj interface{}) {
	marsh := json.NewEncoder(os.Stdout)
	marsh.SetIndent("", "\t")
	marsh.Encode(obj)
}
