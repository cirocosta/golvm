package commands

import (
	"github.com/cirocosta/golvm/lib"
	"github.com/cirocosta/golvm/lvmctl/utils"
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

var Create = cli.Command{
	Name:  "create",
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
			name        = c.String("name")
			size        = c.String("size")
			volumegroup = c.String("volumegroup")
			// 		snapshot = c.String("snapshot")
			// 		thinpool = c.String("thinpool")
			// 		keyfile  = c.String("keyfile")
			args []string
		)

		lvm, err := lib.NewLvm(lib.LvmConfig{})
		utils.Abort(err)

		if name == "" {
			cli.ShowCommandHelp(c, "create")
			utils.Abort(errors.Errorf("Name parameter not set."))
		}

		if volumegroup == "" {
			// pick the volume group that best fits
			//	1.	if size is specified - the one that fits
			//		and has the lowest amount
			//		if size NOT specified - pick the one with
			//		the most free.
		}

		args, err = lvm.BuildLogicalVolumeCretionArgs(&lib.LvCreationConfig{
			Name: name,
			Size: size,
		})
		utils.Abort(err)

		err = lvm.CreateLv(args...)
		utils.Abort(err)

		return
	},
}
