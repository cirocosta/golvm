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
			Name:  "volumegroup",
			Usage: "Volume group to use",
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
			thinpool    = c.String("thinpool")
			snapshot    = c.String("snapshot")
			keyfile     = c.String("keyfile")
		)

		lvm, err := lib.NewLvm(lib.LvmConfig{})
		utils.Abort(err)

		if name == "" {
			cli.ShowCommandHelp(c, "create")
			utils.Abort(errors.Errorf("Name parameter not set."))
		}

		if volumegroup == "" {
			vgs, err := lvm.ListVolumeGroups()
			utils.Abort(err)

			vg, err := lib.PickBestVolumeGroup(0, vgs)
			utils.Abort(err)

			if vg == nil {
				utils.Abort(errors.Errorf(
					"didn't find suitable volume group for specified size"))
			}

			volumegroup = vg.Name
		}

		err = lvm.CreateLv(lib.LvCreationConfig{
			Name:        name,
			Size:        size,
			VolumeGroup: volumegroup,
			Snapshot:    snapshot,
			ThinPool:    thinpool,
			KeyFile:     keyfile,
		})
		utils.Abort(err)

		return
	},
}
