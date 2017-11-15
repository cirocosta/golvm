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
	l.logger.Debug().
		Msg("retrieving physical volumes")

	output, err = l.Run("pvs",
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

	output, err = l.Run("vgs",
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

	output, err = l.Run("lvs",
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

type LvCreationConfig struct {
	Name            string
	Size            string
	Snapshot        string
	KeyFile         string
	ThinPool        string
	VolumeGroup     string
	ThinVolumeGroup string
}

// CreateLogicalVolume creates a logical volume using the definition passed.
// A VolumeGroup my be specified or not.
// notes.:
//	-	a snapshot is a volume that is created from
//		another volume that must already exist.
func (l Lvm) BuildLogicalVolumeCretionArgs(cfg *LvCreationConfig) (args []string, err error) {
	l.logger.Debug().
		Str("name", cfg.Name).
		Str("size", cfg.Size).
		Str("snapshot", cfg.Snapshot).
		Str("thinpool", cfg.ThinPool).
		Str("keyfile", cfg.KeyFile).
		Str("volume-group", cfg.VolumeGroup).
		Msg("starting logical volume creation args building")

	var (
		isThinSnapshot = false // TODO detect this
		isSnapshot     = cfg.Snapshot != ""
		hasSize        = cfg.Size != ""
		hasKeyFile     = cfg.KeyFile != ""
		hasThinPool    = cfg.ThinPool != ""
		finfo          os.FileInfo
	)

	if cfg.Name == "" {
		err = errors.Errorf("Name must be set")
		return
	}

	if cfg.VolumeGroup == "" {
		err = errors.Errorf("VolumeGroup must be specified")
	}

	if hasKeyFile {
		finfo, err = os.Stat(cfg.KeyFile)
		if err != nil {
			err = errors.Wrapf(err,
				"failed to inspect keyfile %s",
				cfg.KeyFile)
			return
		}

		if finfo.IsDir() {
			err = errors.Errorf(
				"keyfile %s must be a file, not a dir",
				cfg.KeyFile)
			return
		}

		_, err = exec.LookPath("cryptsetup")
		if err != nil {
			err = errors.Wrapf(err,
				"cryptsetup not found in PATH")
			return
		}
	}

	if isSnapshot {
		if hasKeyFile {
			err = errors.Errorf("can't have snapshot with keyfile")
			return
		}

		// TODO:check if the volumegroup is thinly provisioned
		//	if is ==> mark as thinsnap
	}

	if !hasSize && !isThinSnapshot {
		err = errors.Errorf("a size must be provided")
		return
	}

	if hasSize && isThinSnapshot {
		err = errors.Errorf("can't specify size for thin snapshots")
		return
	}

	args = []string{"--setactivationskip", "n"}
	args = append(args, "--name", cfg.Name)

	switch {
	case isSnapshot:
		args = append(args, "--snapshot")

		if hasSize {
			args = append(args, "--size", cfg.Size)
		}

		args = append(args, cfg.VolumeGroup+"/"+cfg.Snapshot)
	case hasThinPool:
		args = append(args, "--virtualsize", cfg.Size)
		args = append(args, "--thin")
		args = append(args, cfg.VolumeGroup+"/"+cfg.ThinPool)
	default:
		args = append(args, "--size", cfg.Size)
		args = append(args, cfg.VolumeGroup)
	}

	return
}

// DeleteLogicalVolume deletes a logical volume if it exists
func (l Lvm) DeleteLogicalVolume(def *LogicalVolume) (err error) {
	return
}

// CreateLv runs the 'lvcreate' command with
// the arguments provided.
func (l Lvm) CreateLv(args ...string) (err error) {
	_, err = l.Run("lvcreate", args...)
	return
}

// run executes a given command whose executable
// is 'name' and whose arguments are 'args'.
// The executed command inherits the parent environment
// with the addition of LC_NUMERIC set to en_US.UTF-8 in
// order to prevent the use of commas as the floating point
// separator.
func (l Lvm) Run(name string, args ...string) (out []byte, err error) {
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
