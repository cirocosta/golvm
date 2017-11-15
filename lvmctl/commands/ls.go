package commands

import (
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

var Ls = cli.Command{
	Name: "ls",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "root, r",
			Usage: "Root of the volume creation",
		},
	},
	Action: func(c *cli.Context) (err error) {
		var (
			root = c.String("root")
		)

		if root == "" {
			cli.ShowCommandHelp(c, "create")
			err = errors.Errorf("All parameters must be set.")
			return
		}

		return
	},
}
