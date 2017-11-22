package main

import (
	"os"

	"github.com/cirocosta/golvm/driver"
	"github.com/cirocosta/golvm/lib"
	"github.com/cirocosta/golvm/lvmctl/utils"
	"github.com/rs/zerolog"

	v "github.com/docker/go-plugins-helpers/volume"
)

const (
	socketAddress   = "/run/docker/plugins/golvm.sock"
	volumeMountRoot = "/mnt"
	hostMountPoint  = "/mnt/lvmvol"
	vgWhitelistFile = "/mnt/lvmvol/whitelist"
)

var (
	err     error
	version string = "master-dev"
	logger         = zerolog.New(os.Stdout).
		With().Str("from", "main").
		Logger()
)

func main() {
	l, err := lib.NewLvm(lib.LvmConfig{})
	utils.Abort(err)

	dm, err := driver.NewDirManager(driver.DirManagerConfig{
		Root: volumeMountRoot,
	})
	utils.Abort(err)

	d, err := driver.NewDriver(driver.DriverConfig{
		Lvm:             &l,
		DirManager:      &dm,
		VgWhitelistFile: vgWhitelistFile,
	})
	utils.Abort(err)

	handler := v.NewHandler(d)

	logger.Info().
		Str("address", socketAddress).
		Str("version", version).
		Msg("listening on unix socket")

	err = handler.ServeUnix(socketAddress, 0)
	utils.Abort(err)
}
