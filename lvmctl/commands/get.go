package commands

import (
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

var Get = cli.Command{
	Name: "get",
	Usage: "inspects existing LVM volumes",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "name",
			Usage: "Name of the volume to create",
		},
		&cli.StringFlag{
			Name:  "root, r",
			Usage: "Root of the volume creation",
		},
	},
	Action: func(c *cli.Context) (err error) {
		var (
			name = c.String("name")
			root = c.String("root")
		)

		if name == "" || root == "" {
			cli.ShowCommandHelp(c, "create")
			err = errors.Errorf("All parameters must be set.")
			return
		}

		return
	},
}
