package lib

import (
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

// NewLvm instantiates a new LVm controller instance.
func NewLvm(cfg LvmConfig) (l Lvm, err error) {
	l.logger = zerolog.New(os.Stdout).With().
		Str("from", "lvm").
		Logger()

	return
}

// GetLogicalVolume retrieves a single logical volume
// by its `lv_name`.
// Note.:	if the same `lv_name` exists in two volume groups,
//		the first found is returned.
func (l Lvm) GetLogicalVolume(name string) (vol *LogicalVolume, err error) {
	vols, err := l.ListLogicalVolumes()
	if err != nil {
		err = errors.Wrapf(err,
			"couldn't list logical volumes")
		return
	}

	for _, vol = range vols {
		if vol.LvName != name {
			continue
		}

		return
	}

	vol = nil
	return
}

// ListPhysicalVolumes gathers a list of all the physical
// volumes that can be found by the LVM controller.
// It parses the output from the `pvs` command and returns
// a list of PhysicalVolume structs.
func (l Lvm) ListPhysicalVolumes() (vols []*PhysicalVolume, err error) {
	var output []byte

	l.logger.Debug().
		Msg("listing physical volumes")

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

// ListVolumeGroups lists all groups that can be reached
// by the LVM controller. As a result it parses the response
// of the 'vgs' command and returns a list of VolumeGroup structs.
func (l Lvm) ListVolumeGroups() (vols []*VolumeGroup, err error) {
	var output []byte

	l.logger.Debug().
		Msg("listing volume groups")

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

// FormatDevice format a `device` with a filesystem of
// a particular `fstype`.
// Allowed `fsType`s are:
//	-	ext4
//	-	xfs
func (l Lvm) FormatDevice(device, fsType string) (err error) {
	var args []string

	args, err = BuildMakeFsArgs(fsType, device)
	if err != nil {
		err = errors.Wrapf(err,
			"couldn't build args for formatting device")
		return
	}

	_, err = l.Run("mkfs", args...)
	return
}

// IsDeviceFormatted checks whether a given `device`
// is already formatted with a filesystem.
func (l Lvm) IsDeviceFormatted(device string) (isFormatted bool, err error) {
	var response []byte

	args, err := BuildGetDeviceFormatArgs(device)
	if err != nil {
		return
	}

	response, err = l.Run("lsblk", args...)
	if err != nil {
		return
	}

	isFormatted = string(response) == ""
	return
}

// ListLogicalVolumes retrieves a list of LogicalVolume structs
// from the result of parsing the response of the 'lvs' command.
func (l Lvm) ListLogicalVolumes() (vols []*LogicalVolume, err error) {
	var output []byte

	l.logger.Debug().
		Msg("retrieving logical volumes")

	output, err = l.Run("lvs",
		"--units=m",
		"--nosuffix",
		"--noheadings",
		"--options=lv_all",
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

	for _, vol := range vols {
		if vol.VgName == "" && vol.LvFullName != "" {
			parts := strings.SplitN(vol.LvFullName, "/", 2)
			if len(parts) != 2 {
				continue
			}
			vol.VgName = parts[0]
		}
	}

	return
}

// Mount runs the 'mount' command with the arguments provided.
func (l Lvm) Mount(device, location string) (err error) {
	_, err = l.Run("mount", device, location)
	return
}

// CreateLv runs the 'lvremove' command with
// the arguments provided.
func (l Lvm) RemoveLv(cfg LvRemovalConfig) (err error) {
	var args []string

	args, err = BuildLogicalVolumeRemovalArgs(cfg)
	if err != nil {
		err = errors.Wrapf(err, "failed to create lv removal args")
		return
	}

	_, err = l.Run("lvremove", args...)
	return
}

// CreateLv runs the 'lvcreate' command with
// the arguments provided.
func (l Lvm) CreateLv(cfg LvCreationConfig) (err error) {
	var args []string

	args, err = BuildLogicalVolumeCretionArgs(cfg)
	if err != nil {
		err = errors.Wrapf(err, "failed to create lv cretion args")
		return
	}

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
