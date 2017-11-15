package commands

import (
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

var Create = cli.Command{
	Name: "create",
	Usage: "create an LVM volume",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "name",
			Usage: "Name of the volume to create",
		},
		&cli.StringFlag{
			Name:  "size",
			Usage: "Maximum size of the volume",
		},
		&cli.StringFlag{
			Name:  "thinpool",
			Usage: "Name of the thinpool to base the volume",
		},
		&cli.StringFlag{
			Name:  "snapshot",
			Usage: "Volume to get an snapshot from",
		},
		&cli.StringFlag{
			Name:  "keyfile",
			Usage: "Keyfile to encrypt the volume",
		},
		&cli.StringFlag{
			Name:  "root, r",
			Usage: "Root of the volume creation",
		},
	},
	Action: func(c *cli.Context) (err error) {
		var (
			name = c.String("name")
			// 		size     = c.String("size")
			// 		snapshot = c.String("snapshot")
			// 		thinpool = c.String("thinpool")
			// 		keyfile  = c.String("keyfile")
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
