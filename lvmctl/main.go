package main

import (
	"os"

	"github.com/cirocosta/golvm/lvmctl/commands"
	"gopkg.in/urfave/cli.v2"
)

var (
	version string = "master-dev"
)

func main() {
	var app = &cli.App{
		Name:    "lvmctl",
		Version: version,
		Usage:   "Controls the 'golvm' volume plugin",
		Commands: []*cli.Command{
			&commands.Check,
			&commands.Create,
			&commands.Get,
			&commands.Ls,
			&commands.Rm,
		},
	}

	app.Run(os.Args)
}
