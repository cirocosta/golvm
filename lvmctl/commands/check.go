package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/cirocosta/golvm/lib"
	"gopkg.in/urfave/cli.v2"
)

func abort (err error) {
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

		pvs, err := lvm.ListPhysicalVolumes()
		abort(err)

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)
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

		// 1. list pv
		// 2. list vg
		// 3. list lvs

		return
	},
}
