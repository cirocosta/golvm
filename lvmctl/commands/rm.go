package commands

import (
	"github.com/cirocosta/golvm/lib"
	"github.com/cirocosta/golvm/lvmctl/utils"
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

var Rm = cli.Command{
	Name:  "rm",
	Usage: "removes existing LVM volumes",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "name",
			Usage: "Name of the logical volume to remove",
		},
		&cli.StringFlag{
			Name:  "volumegroup",
			Usage: "Name of the volume group that the lv is from",
		},
	},
	Action: func(c *cli.Context) (err error) {
		var (
			name        = c.String("name")
			volumegroup = c.String("volumegroup")
			args        []string
		)

		lvm, err := lib.NewLvm(lib.LvmConfig{})
		utils.Abort(err)

		if name == "" || volumegroup == "" {
			cli.ShowCommandHelp(c, "rm")
			utils.Abort(errors.Errorf("All parameters must be set."))
		}

		args, err = lib.BuildLogicalVolumeRemovalArgs(lib.LvRemovalConfig{
			LvName: name,
			VgName: volumegroup,
		})
		utils.Abort(err)

		err = lvm.RemoveLv(args...)
		utils.Abort(err)

		return
	},
}
