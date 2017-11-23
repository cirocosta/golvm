package commands

import (
	"fmt"

	"github.com/cirocosta/golvm/lib"
	"github.com/cirocosta/golvm/lvmctl/utils"
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

var Get = cli.Command{
	Name:  "get",
	Usage: "inspects existing LVM volumes",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "name",
			Usage: "Name of the volume to inspect",
		},
	},
	Action: func(c *cli.Context) (err error) {
		var (
			name          = c.String("name")
			desiredVolume *lib.LogicalVolume
			attr          *lib.LvAttr
		)

		if name == "" {
			cli.ShowCommandHelp(c, "get")
			err = errors.Errorf("All parameters must be set.")
			return
		}

		lvm, err := lib.NewLvm(lib.LvmConfig{})
		utils.Abort(err)

		lvs, err := lvm.ListLogicalVolumes()
		utils.Abort(err)

		for _, lv := range lvs {
			if lv.LvName != name {
				continue
			}

			desiredVolume = lv
			break
		}

		if desiredVolume == nil {
			utils.Abort(errors.Errorf(
				"volume named %s not found", name))
		}

		fmt.Printf("NAME\t\t%s\n", desiredVolume.LvFullName)
		fmt.Printf("POOL\t\t%s\n", desiredVolume.PoolLv)
		fmt.Printf("SIZE\t\t%f\n", desiredVolume.LvSize)
		fmt.Printf("ATTR\t\t%s\n", desiredVolume.LvAttr)

		attr, err = lib.ParseLvAttr(desiredVolume.LvAttr)
		utils.Abort(err)

		fmt.Printf("ALLOC\t\t%s\n", attr.AllocationPolicy)
		fmt.Printf("STATE\t\t%s\n", attr.State)
		fmt.Printf("DEV_STATE\t\t%s\n", attr.DeviceState)
		fmt.Printf("VOL_TYPE\t\t%s\n", attr.VolumeType)
		fmt.Printf("VOL_HEALTH\t\t%s\n", attr.VolumeHealth)
		fmt.Printf("TARGET_TYPE\t\t%s\n", attr.TargetType)

		return
	},
}
