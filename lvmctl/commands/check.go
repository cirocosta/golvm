package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/cirocosta/golvm/lib"
	"github.com/cirocosta/golvm/lvmctl/utils"
	"gopkg.in/urfave/cli.v2"
)

var Check = cli.Command{
	Name:  "check",
	Usage: "verifies the environment by dumping all that was found",
	Action: func(c *cli.Context) (err error) {
		lvm, err := lib.NewLvm(lib.LvmConfig{})
		utils.Abort(err)

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)

		pvs, err := lvm.ListPhysicalVolumes()
		utils.Abort(err)

		fmt.Println("")
		fmt.Println("PHYSICAL VOLUMES")
		fmt.Fprintln(w, "NAME\tVG\tSIZE\tFREE\t")
		for _, pv := range pvs {
			fmt.Fprintf(w, "%s\t%s\t%.2f\t%.2f\n",
				pv.PhysicalVolume,
				pv.VolumeGroup,
				pv.PhysicalSize,
				pv.PhysicalSizeFree)
		}
		w.Flush()

		vgs, err := lvm.ListVolumeGroups()
		utils.Abort(err)

		fmt.Println("")
		fmt.Println("VOLUME GROUPS")
		fmt.Fprintln(w, "NAME\tSIZE\tFREE\t")
		for _, vg := range vgs {
			fmt.Fprintf(w, "%s\t%.2f\t%.2f\n",
				vg.Name,
				vg.Size,
				vg.Free)
		}
		w.Flush()

		lvs, err := lvm.ListLogicalVolumes()
		utils.Abort(err)

		fmt.Println("")
		fmt.Println("LOGICAL VOLUMES")
		fmt.Fprintln(w, "NAME\tVG\tSIZE\tPOOL\t")
		for _, lv := range lvs {
			fmt.Fprintf(w, "%s\t%s\t%.2f\t%s\n",
				lv.LvName,
				lv.VgName,
				lv.LvSize,
				lv.PoolLv)
		}
		w.Flush()

		return
	},
}
