package commands

import (
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

var Resize = cli.Command{
	Name:  "resize",
	Usage: "resizes a volume",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "name",
			Usage: "Name of the volume to resize",
		},
		&cli.StringFlag{
			Name:  "size",
			Usage: "Desired size to get reduced / expanded to",
		},
	},
	Action: func(c *cli.Context) (err error) {
		var (
			name = c.String("name")
			size = c.String("size")
		)

		if name == "" || size == "" {
			cli.ShowCommandHelp(c, "resize")
			err = errors.Errorf("All parameters must be set.")
			return
		}

		//	1.	get the desired size
		//	2.	check if exists;

		return
	},
}
