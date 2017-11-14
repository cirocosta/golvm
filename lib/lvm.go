package lib

import (
	"encoding/json"
	"os"
	"os/exec"

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
		"--nosuffix",
		"--noheadings",
		"--report-format=json")
	if err != nil {
		err = errors.Wrapf(err,
			"failed to retrieve physical volumes")
		return
	}

	return
}

func DecodePhysicalVolumesResponse(response []byte) (infos []*PhysicalInfo, err error) {
	if response == nil {
		err = errors.Errorf("response can't be nil")
		return
	}

	if len(response) == 0 {
		err = errors.Errorf("can't decode empty response")
		return
	}

	var report = new(PhysicalVolumesReport)
	err = json.Unmarshal(response, report)
	if err != nil {
		err = errors.Wrapf(err, "errored decoding vgs response")
		return
	}

	if len(report.Report) != 1 {
		err = errors.Errorf(
			"unexpected number of responses decoded - %s",
			response)
		return
	}

	infos = report.Report[0].Pv
	return
}

func (l Lvm) VolumeGroups() (output string, err error) {
	l.logger.Debug().Msg("retrieving volume groups")

	output, err = l.run("vgs",
		"--units=m",
		"--nosuffix",
		"--noheadings",
		"--report-format=json")
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
		"--nosuffix",
		"--noheadings",
		"--report-format=json")
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
