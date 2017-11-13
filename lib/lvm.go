package lib

import (
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Lvm struct {
	logger zerolog.Logger
}

type LvmConfig struct{}

func NewLvm(cfg LvmConfig) (l Lvm, err error) {
	l.logger = zerolog.New(os.Stdout).With().
		Str("from", "lvm").
		Logger()

	return
}

func (l Lvm) PhysicalVolumes() (output string, err error) {
	l.logger.Debug().Msg("retrieving physical volumes")

	output, err = l.run("pvs",
		"--units=m",
		"--separator=:",
		"--nosuffix",
		"--noheadings")
	if err != nil {
		err = errors.Wrapf(err,
			"failed to retrieve physical volumes")
		return
	}

	return
}

func ParsePhysicalVolumesResponseLine(line string) (p *PhysicalInfo, err error) {
	var (
		parts            []string
		physicalSize     float64
		physicalSizeFree float64
	)

	if line == "" {
		err = errors.Errorf("can't parse infomation from empty line")
		return
	}

	line = strings.TrimSpace(line)

	parts = strings.SplitN(line, ":", 6)
	if len(parts) != 6 {
		err = errors.Errorf(
			"malformed line %s - expected 6 fields delimited by colon",
			line)
		return
	}

	physicalSize, err = strconv.ParseFloat(parts[4], 64)
	if err != nil {
		err = errors.Wrapf(err,
			"couldn't convert physical size %s to float",
			parts[4])
		return
	}

	physicalSizeFree, err = strconv.ParseFloat(parts[5], 64)
	if err != nil {
		err = errors.Wrapf(err,
			"couldn't convert physical size free %s to float",
			parts[5])
		return
	}

	p = &PhysicalInfo{
		PhysicalVolume:   parts[0],
		VolumeGroup:      parts[1],
		Fmt:              parts[2],
		Attr:             parts[3],
		PhysicalSize:     physicalSize,
		PhysicalFreeSize: physicalSizeFree,
	}

	return
}

func (l Lvm) VolumeGroups() (output string, err error) {
	l.logger.Debug().Msg("retrieving volume groups")

	output, err = l.run("vgs",
		"--units=m",
		"--separator=:",
		"--nosuffix",
		"--noheadings")
	if err != nil {
		err = errors.Wrapf(err,
			"failed to retrieve volume groups")
		return
	}

	return
}

func (l Lvm) LogicalVolumes() (output string, err error) {
	l.logger.Debug().Msg("retrieving logical volumes")

	output, err = l.run("lvs",
		"--units=m",
		"--separator=:",
		"--nosuffix",
		"--noheadings")
	if err != nil {
		err = errors.Wrapf(err,
			"failed to retrieve logical volumes")
		return
	}
	return
}

func (l Lvm) run(name string, args ...string) (output string, err error) {
	l.logger.Debug().
		Str("cmd", name).
		Strs("args", args).
		Msg("executing command")

	cmd := exec.Command(name, args...)
	out, err := cmd.Output()
	if err != nil {
		err = errors.Wrapf(err,
			"failed to execute command %s with args %+v",
			name, args)
		return
	}

	output = string(out)
	return
}
