package commands

import (
	"gopkg.in/urfave/cli.v2"
)

var Check = cli.Command{
	Name: "check",
	Usage: "checks verifies the environment",
	Action: func(c *cli.Context) (err error) {
		return
	},
}
