package main

import (
	"os"

	"github.com/cirocosta/golvm/driver"
	"github.com/cirocosta/golvm/lib"
	"github.com/rs/zerolog"

	v "github.com/docker/go-plugins-helpers/volume"
)

const (
	socketAddress = "/run/docker/plugins/golvm.sock"
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
	if err != nil {
		logger.Error().
			Err(err).
			Msg("failed to initialize lvm manager")
		os.Exit(1)
	}

	d, err := driver.NewDriver(driver.DriverConfig{
		Lvm: &l,
	})
	if err != nil {
		logger.Error().
			Err(err).
			Msg("failed to initialize lvm volume driver")
		os.Exit(1)
	}

	handler := v.NewHandler(d)

	logger.Info().
		Str("address", socketAddress).
		Str("version", version).
		Msg("listening on unix socket")

	err = handler.ServeUnix(socketAddress, 0)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("errored serving on unix socket")
		os.Exit(1)
	}
}
