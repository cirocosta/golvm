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
		Commands: []*cli.Command{
			&commands.Create,
			&commands.Get,
			&commands.Ls,
			&commands.Rm,
		},
		Name:    "lvmctl",
		Version: version,
		Usage:   "Controls the 'golvm' volume plugin",
	}

	app.Run(os.Args)
}
