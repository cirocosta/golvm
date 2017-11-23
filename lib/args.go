package lib

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

func BuildGetDeviceFormatArgs(device string) (args []string, err error) {
	if device == "" {
		err = errors.Errorf("a device must be specified")
		return
	}

	args = []string{
		"--noheadings",
		"--discard",
		"--output=FSTYPE",
		device,
	}

	return
}

// BuildMakeFsArgs builds a list of arguments to be used
// on the 'mkfs' command to properly create a filesystem on
// a device. It supports two types of FS:
//	-	xfs
//	-	ext4
func BuildMakeFsArgs(fsType, device string) (args []string, err error) {
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

// CreateLogicalVolume creates a logical volume using the definition passed.
// A VolumeGroup my be specified or not.
// notes.:
//	-	a snapshot is a volume that is created from
//		another volume that must already exist.
func BuildLogicalVolumeCretionArgs(cfg *LvCreationConfig) (args []string, err error) {
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
func BuildLogicalVolumeRemovalArgs(cfg LvRemovalConfig) (args []string, err error) {
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
