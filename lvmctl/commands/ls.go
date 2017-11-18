package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/cirocosta/golvm/lib"
	"github.com/cirocosta/golvm/lvmctl/utils"
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

var Ls = cli.Command{
	Name:  "ls",
	Usage: "lists existing LVM volumes",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "volumegroup",
			Usage: "the vg to look for volumes",
		},
	},
	Action: func(c *cli.Context) (err error) {
		var (
			volumegroup = c.String("volumegroup")
		)

		if volumegroup == "" {
			cli.ShowCommandHelp(c, "ls")
			utils.Abort(errors.New(
				"a volumegroup must be specified"))
		}

		lvm, err := lib.NewLvm(lib.LvmConfig{})
		utils.Abort(err)

		w := new(tabwriter.Writer)
		w.Init(os.Stdout, 0, 8, 0, '\t', 0)

		lvs, err := lvm.ListLogicalVolumes()
		utils.Abort(err)

		fmt.Println("")
		fmt.Println("LOGICAL VOLUMES")
		fmt.Fprintln(w, "NAME\tVG\tSIZE\tPOOL\t")
		for _, lv := range lvs {
			if lv.VgName != volumegroup {
				continue
			}

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
