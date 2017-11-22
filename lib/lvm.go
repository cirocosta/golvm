package lib

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// Lvm encapsulates a series of methods for
// dealing with LVM management.
// It's mostly stateless except for a logger.
type Lvm struct {
	logger zerolog.Logger
}

// LvmConfig provides the configuration details for
// the Lvm helper.
type LvmConfig struct{}

func NewLvm(cfg LvmConfig) (l Lvm, err error) {
	l.logger = zerolog.New(os.Stdout).With().
		Str("from", "lvm").
		Logger()

	return
}

// ParseLvAttr takes an 'attr' string from 'lvs' command and
// parses it so that it can be consumed via the LvAttr struct.
// In case of any unexpected tokens or malformed attr, fails
// with an error.
func ParseLvAttr(attr string) (parsedAttr *LvAttr, err error) {
	if attr == "" {
		err = errors.Errorf("attr must not be empty")
		return
	}

	if len(attr) != 10 {
		err = errors.Errorf(
			"malformed attr '%s' - must be 10chars",
			attr)
		return
	}

	var (
		mapping map[string]string
		present bool
		val     string
		chars   = strings.Split(attr, "")
	)

	parsedAttr = new(LvAttr)
	for ndx, character := range chars {
		mapping = lvAttrMapper[ndx]
		val, present = mapping[character]
		if !present {
			err = errors.Errorf(
				"unexpected character '%' for lv attr '%d'",
				character, ndx)
			return
		}

		switch ndx {
		case 0:
			parsedAttr.VolumeType = val
		case 1:
			parsedAttr.Permissions = val
		case 2:
			parsedAttr.AllocationPolicy = val
		case 3:
			parsedAttr.FixedMinor = val
		case 4:
			parsedAttr.State = val
		case 5:
			parsedAttr.DeviceState = val
		case 6:
			parsedAttr.TargetType = val
		case 7:
			parsedAttr.OverrideNewBlocksWithZero = val
		case 8:
			parsedAttr.VolumeHealth = val
		case 9:
			parsedAttr.SkipActivation = val
		}
	}

	return
}

// PickBestVolumeGroup picks the volume group that best
// accomodates a given size.
// 'size' specifies the size to be accomodated - if 0, any
// volume with free space fits it.
func PickBestVolumeGroup(size float64, vols []*VolumeGroup) (bestVol *VolumeGroup, err error) {
	if vols == nil {
		err = errors.Errorf("can't pick best volume from nil list of vols")
		return
	}

	for _, vol := range vols {
		if vol.Free > size {
			if bestVol == nil {
				bestVol = vol
				continue
			}

			if bestVol.Free > vol.Free {
				continue
			}

			bestVol = vol
		}
	}

	return
}

// ListPhysicalVolumes gathers a list of all the physical
// volumes that can be found by the LVM controller.
// It parses the output from the `pvs` command and returns
// a list of PhysicalVolume structs.
func (l Lvm) ListPhysicalVolumes() (vols []*PhysicalVolume, err error) {
	var output []byte

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

	vols, err = DecodePhysicalVolumesResponse(output)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to decode physical volumes response")
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

// ListVolumeGroups lists all groups that can be reached
// by the LVM controller. As a result it parses the response
// of the 'vgs' command and returns a list of VolumeGroup structs.
func (l Lvm) ListVolumeGroups() (vols []*VolumeGroup, err error) {
	var output []byte

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

	vols, err = DecodeVolumeGroupsResponse(output)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to devoce volume groups response")
		return
	}

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

// ListLogicalVolumes retrieves a list of LogicalVolume structs
// from the result of parsing the response of the 'lvs' command.
func (l Lvm) ListLogicalVolumes() (vols []*LogicalVolume, err error) {
	var output []byte

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

	vols, err = DecodeLogicalVolumesResponse(output)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to devoce logical volumes response")
		return
	}

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

// LvCreationConfig is a simplified configuration
// struct to be passed to logical volume creation
// methods.
type LvCreationConfig struct {
	Name        string
	Size        string
	Snapshot    string
	KeyFile     string
	ThinPool    string
	VolumeGroup string
	FsType      string
}

// BuildVolumeMountArgs builds a list of arguments to
// be used on the 'mount' command to properly mount
// a given device to a location in the filesystem hierarchy.
func (l Lvm) BuildVolumeMountArgs(location string) (err error) {
	if location == "" {
		err = errors.Errorf("a location must be specified")
		return
	}

	return
}

// BuildMakeFsArgs builds a list of arguments to be used
// on the 'mkfs' command to properly create a filesystem on
// a device. It supports two types of FS:
//	-	xfs
//	-	ext4
func (l Lvm) BuildMakeFsArgs(fsType, device string) (args []string, err error) {
	if fsType == "" || device == "" {
		err = errors.Errorf("both fstype and device must be specified")
		return
	}

	args = []string{"-t"}

	switch fsType {
	case "ext4":
	case "xfs":
	default:
		err = errors.Errorf("unsupported fs type %s", fsType)
		return
	}

	args = append(args, fsType, device)

	return
}

// MakeFs runs the 'mkfs' command with the arguments provided.
func (l Lvm) MakeFs(args ...string) (err error) {
	_, err = l.Run("mkfs", args...)
	return
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

type LvRemovalConfig struct {
	LvName string
	VgName string
}

// DeleteLogicalVolume deletes a logical volume if it exists
func (l Lvm) BuildLogicalVolumeRemovalArgs(cfg LvRemovalConfig) (args []string, err error) {
	if cfg.LvName == "" {
		err = errors.Errorf(
			"the logical volume name must be specified")
		return
	}

	if cfg.VgName == "" {
		err = errors.Errorf(
			"the volume group name must be specified")
		return
	}

	args = []string{"--force", cfg.VgName + "/" + cfg.LvName}
	return
}

// CreateLv runs the 'lvremove' command with
// the arguments provided.
func (l Lvm) RemoveLv(args ...string) (err error) {
	_, err = l.Run("lvremove", args...)
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

	out, err = cmd.CombinedOutput()
	if err != nil {
		err = errors.Wrapf(err,
			"failed to execute command '%s' with args '%+v'. Output:\n%s\n",
			name, args, string(out))
		return
	}

	return
}
