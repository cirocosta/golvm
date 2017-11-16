package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/cirocosta/golvm/lib"
	"gopkg.in/urfave/cli.v2"
)

func abort(err error) {
	if err == nil {
		return
	}

	fmt.Printf("ERRORED: %+v\nAborting.\n", err)
	os.Exit(1)
}

var Check = cli.Command{
	Name:  "check",
	Usage: "checks verifies the environment",
	Action: func(c *cli.Context) (err error) {
		lvm, err := lib.NewLvm(lib.LvmConfig{})
		abort(err)

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)

		pvs, err := lvm.ListPhysicalVolumes()
		abort(err)

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
		abort(err)

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
		abort(err)

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
