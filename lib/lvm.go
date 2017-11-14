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

func (l Lvm) PhysicalVolumes() (output []byte, err error) {
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

// DecodePhysicalVolumesReponse takes a JSON response from
// the execition of the 'pvs' command and returns a slice of
// PhysicalVolume structs.
func DecodePhysicalVolumesResponse(response []byte) (infos []*PhysicalVolume, err error) {
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
		err = errors.Wrapf(err, "errored decoding pvs response")
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

// DecodeVolumeGroupsRepsponse takes a JSON response from
// the execition of the 'vgs' command and returns a slice of
// VolumeGroup structs.
func DecodeVolumeGroupsResponse(response []byte) (infos []*VolumeGroup, err error) {
	if response == nil {
		err = errors.Errorf("response can't be nil")
		return
	}

	if len(response) == 0 {
		err = errors.Errorf("can't decode empty response")
		return
	}

	var report = new(VolumeGroupsReport)
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

	infos = report.Report[0].Vg
	return
}

// DecodeLogicalVolumesResponse takes a JSON response from
// the execition of the 'lvs' command and returns a slice of
// VolumeGroup structs.
func DecodeLogicalVolumesResponse(response []byte) (infos []*LogicalVolume, err error) {
	if response == nil {
		err = errors.Errorf("response can't be nil")
		return
	}

	if len(response) == 0 {
		err = errors.Errorf("can't decode empty response")
		return
	}

	var report = new(LogicalVolumesReport)
	err = json.Unmarshal(response, report)
	if err != nil {
		err = errors.Wrapf(err, "errored decoding lvs response")
		return
	}

	if len(report.Report) != 1 {
		err = errors.Errorf(
			"unexpected number of responses decoded - %s",
			response)
		return
	}

	infos = report.Report[0].Lv
	return
}

func (l Lvm) VolumeGroups() (output []byte, err error) {
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

// LogicalVolumes retrieves a list of LogicalVolume structs.
func (l Lvm) LogicalVolumes() (output []byte, err error) {
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

// run executes a given command whose executable
// is 'name' and whose arguments are 'args'.
// The executed command inherits the parent environment
// with the addition of LC_NUMERIC set to en_US.UTF-8 in
// order to prevent the use of commas as the floating point
// separator.
func (l Lvm) run(name string, args ...string) (out []byte, err error) {
	l.logger.Debug().
		Str("cmd", name).
		Strs("args", args).
		Msg("executing command")

	cmd := exec.Command(name, args...)
	cmd.Env = append(os.Environ(), "LC_NUMERIC=en_US.UTF-8")

	out, err = cmd.Output()
	if err != nil {
		err = errors.Wrapf(err,
			"failed to execute command %s with args %+v",
			name, args)
		return
	}

	return
}
