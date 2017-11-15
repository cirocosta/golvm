package commands

import (
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

var Rm = cli.Command{
	Name: "rm",
	Usage: "removes existing LVM volumes",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "name",
			Usage: "Name of the volume to create",
		},
	},
	Action: func(c *cli.Context) (err error) {
		var (
			name = c.String("name")
		)

		if name == "" {
			cli.ShowCommandHelp(c, "create")
			err = errors.Errorf("All parameters must be set.")
			return
		}

		return
	},
}
