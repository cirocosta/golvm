package driver

import (
	"os"
	"sync"

	"github.com/cirocosta/golvm/lib"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	v "github.com/docker/go-plugins-helpers/volume"
)

type Driver struct {
	lvm         *lib.Lvm
	dirManager  *DirManager
	logger      zerolog.Logger
	vgWhiteList map[string]bool
	mountsFile  string

	sync.Mutex
}

type DriverConfig struct {
	Lvm             *lib.Lvm
	DirManager      *DirManager
	VgWhitelistFile string
	MountsFile      string
}

func NewDriver(cfg DriverConfig) (d Driver, err error) {
	var whitelist map[string]bool

	if cfg.MountsFile == "" {
		err = errors.Errorf("MountsFile must be specified")
		return
	}

	if cfg.Lvm == nil {
		err = errors.Errorf("Lvm must be specified")
		return
	}

	if cfg.DirManager == nil {
		err = errors.Errorf("DirManager must be specified")
		return
	}

	if cfg.VgWhitelistFile == "" {
		err = errors.Errorf("a VgWhitelistFile must be specified")
		return
	}

	d.logger = zerolog.New(os.Stdout).
		With().
		Str("from", "driver").
		Logger()

	whitelist, err = ReadVgWhitelist(cfg.VgWhitelistFile)
	if err != nil {
		d.logger.Error().
			Err(err).
			Str("file", cfg.VgWhitelistFile).
			Msg("couldn't read vgs from whitelist file")
	}

	for vg, _ := range whitelist {
		d.logger.Info().
			Str("vg", vg).
			Msg("vg whitelisted")
	}

	d.lvm = cfg.Lvm
	d.dirManager = cfg.DirManager
	d.vgWhiteList = whitelist
	d.mountsFile = cfg.MountsFile
	d.logger.Info().Msg("driver initialized")

	return
}

// Create creates a new logical volume if it doesn't
// exist yet.
// It takes few possible options:
//	-	size:		size (or virtualsize) to allocate for
//				the volume
//	-	thinpool:	creates a thinly-provisioned volume
//				from the thinpool specified
//	-	snapshot:	creates a snapshot volume out of the
//				specified volume
//	-	keyfile:	file to use on luks encryption
//	-	volumegroup:	volume group to use or pick the
//				best one from the pool of whitelisted
//				volumegroups.
//	-	fstype:		type of filesystem to use in
//				the volume (ext4 or xfs)
func (d Driver) Create(req *v.CreateRequest) (err error) {
	var (
		size        string
		thinpool    string
		snapshot    string
		keyfile     string
		volumegroup string
		fstype      string
	)

	d.logger.Debug().
		Str("name", req.Name).
		Interface("opts", req.Options).
		Msg("starting creation")

	d.Lock()
	defer d.Unlock()

	size, _ = req.Options["size"]
	thinpool, _ = req.Options["thinpool"]
	snapshot, _ = req.Options["snapshot"]
	keyfile, _ = req.Options["keyfile"]
	volumegroup, _ = req.Options["volumegroup"]
	fstype, _ = req.Options["fstype"]

	if volumegroup == "" {
		vgs, err := d.lvm.ListVolumeGroups()
		err = errors.Wrapf(err,
			"failed to list volume groups")

		var validVgs = make([]*lib.VolumeGroup, 0)
		for _, potentialVg := range vgs {
			_, present := d.vgWhiteList[potentialVg.Name]
			if present {
				validVgs = append(validVgs, potentialVg)
			}
		}

		vg, err := lib.PickBestVolumeGroup(0, validVgs)
		err = errors.Wrapf(err,
			"failed to pick the best volume group")

		if vg == nil {
			err = errors.Errorf(
				"didn't find suitable vg for specified size")
		}

		volumegroup = vg.Name
	}

	err = d.lvm.CreateLv(lib.LvCreationConfig{
		Name:        req.Name,
		Size:        size,
		ThinPool:    thinpool,
		Snapshot:    snapshot,
		KeyFile:     keyfile,
		VolumeGroup: volumegroup,
		FsType:      fstype,
	})
	if err != nil {
		err = errors.Wrapf(err, "failed to create logical volume")
		return
	}

	d.logger.Debug().
		Str("name", req.Name).
		Msg("finished creation")

	return
}

func (d Driver) List() (resp *v.ListResponse, err error) {
	var vols []*lib.LogicalVolume

	d.logger.Debug().
		Msg("starting list")

	d.Lock()
	defer d.Unlock()

	d.logger.Debug().
		Msg("listing volumes")

	vols, err = d.lvm.ListLogicalVolumes()
	if err != nil {
		err = errors.Wrapf(err, "couldn't list volumes")
		return
	}

	var volumesList = make([]*v.Volume, 0)
	for _, vol := range vols {
		volumesList = append(volumesList, &v.Volume{
			Name: vol.LvName,
		})
	}

	resp = &v.ListResponse{
		Volumes: volumesList,
	}
	return
}

func (d Driver) Get(req *v.GetRequest) (resp *v.GetResponse, err error) {
	var (
		mountpoint string
		vol        *lib.LogicalVolume
	)

	d.logger.Debug().
		Str("name", req.Name).
		Msg("starting get")

	d.Lock()
	defer d.Unlock()

	vol, err = d.lvm.GetLogicalVolume(req.Name)
	if err != nil {
		err = errors.Wrapf(err,
			"errored searching for volume named %s",
			req.Name)
		return
	}

	if vol == nil {
		err = errors.Errorf(
			"Couldn't find path for volume %s",
			req.Name)
		return
	}

	mountpoint, _, err = d.dirManager.Get(req.Name)
	if err != nil {
		err = errors.Errorf("couldn't look for mountpoint of volume %s", req.Name)
		return
	}

	resp = &v.GetResponse{
		Volume: &v.Volume{
			Name:       req.Name,
			Mountpoint: mountpoint,
		},
	}

	return
}

func (d Driver) Remove(req *v.RemoveRequest) (err error) {
	var (
		vol             *lib.LogicalVolume
		mountpointFound bool
	)

	d.logger.Debug().
		Str("name", req.Name).
		Msg("starting removal")

	d.Lock()
	defer d.Unlock()

	vol, err = d.lvm.GetLogicalVolume(req.Name)
	if err != nil {
		err = errors.Wrapf(err,
			"errored retrieving logical volume")
		return
	}

	if vol == nil {
		err = errors.Errorf(
			"logical volume %s not found in LVM",
			req.Name)
		return
	}

	_, mountpointFound, err = d.dirManager.Get(req.Name)
	if err != nil {
		err = errors.Wrapf(err,
			"errored searching for mountpoint of volume %s",
			req.Name)
		return
	}

	if mountpointFound {
		// unmount
		// remove
	}

	err = d.lvm.RemoveLv(lib.LvRemovalConfig{
		LvName: vol.LvName,
		VgName: vol.VgName,
	})

	// check if there's still a directory
	//	if true: umount and then delete the directory
	// remove the volume

	return
}

func (d Driver) Path(req *v.PathRequest) (resp *v.PathResponse, err error) {
	var (
		mountpoint string
		vol        *lib.LogicalVolume
	)

	d.logger.Debug().
		Str("name", req.Name).
		Msg("starting path")

	d.Lock()
	defer d.Unlock()

	vol, err = d.lvm.GetLogicalVolume(req.Name)
	if err != nil {
		err = errors.Wrapf(err,
			"errored searching for volume named %s",
			req.Name)
		return
	}

	if vol == nil {
		err = errors.Errorf(
			"Couldn't find path for volume %s",
			req.Name)
		return
	}

	mountpoint, _, err = d.dirManager.Get(req.Name)
	if err != nil {
		err = errors.Errorf(
			"couldn't look for mountpoint of volume %s",
			req.Name)
		return
	}

	resp = &v.PathResponse{
		Mountpoint: mountpoint,
	}
	return
}

func (d Driver) IsLocationMounted(location string) (isMounted bool, err error) {
	var infos []*lib.MountInfo

	infos, err = lib.ParseMountsFile(d.mountsFile)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to parse mounts from file %s",
			d.mountsFile)
		return
	}

	for _, info := range infos {
		if info.Location == location {
			isMounted = true
			return
		}
	}

	return
}

func (d Driver) Mount(req *v.MountRequest) (resp *v.MountResponse, err error) {
	var (
		vol         *lib.LogicalVolume
		mountpoint  string
		found       bool
		isFormatted bool
		isMounted   bool
	)

	d.logger.Debug().
		Str("name", req.Name).
		Str("ID", req.ID).
		Msg("starting mount")

	d.Lock()
	defer d.Unlock()

	vol, err = d.lvm.GetLogicalVolume(req.Name)
	if err != nil {
		err = errors.Wrapf(err,
			"couldn't found volume %s to mount",
			req.Name)
		return
	}

	mountpoint, found, err = d.dirManager.Get(req.Name)
	if err != nil {
		err = errors.Errorf(
			"couldn't look for mountpoint of volume %s",
			req.Name)
		return
	}

	if !found {
		mountpoint, err = d.dirManager.Create(req.Name)
		if err != nil {
			err = errors.Wrapf(err,
				"couldn't create mountpoint for volume %s",
				req.Name)
			return
		}
	} else {
		isMounted, err = d.IsLocationMounted(mountpoint)
		if err != nil {
			err = errors.Wrapf(err, "failed retrieving mount info list")
			return
		}

		if isMounted {
			resp = &v.MountResponse{
				Mountpoint: mountpoint,
			}
			return
		}
	}

	if vol.LvDmPath == "" {
		err = errors.Errorf(
			"can't find the device for volume %s",
			req.Name)
		return
	}

	isFormatted, err = d.lvm.IsDeviceFormatted(vol.LvDmPath)
	if err != nil {
		err = errors.Errorf(
			"couldn't check if device %s is formated",
			vol.LvDmPath)
		return
	}

	if !isFormatted {
		err = d.lvm.FormatDevice(vol.LvDmPath, "ext4")
		if err != nil {
			err = errors.Errorf(
				"couldn't format device %s as %s",
				vol.LvDmPath, "ext4")
		}
	}

	err = d.lvm.Mount(vol.LvDmPath, mountpoint)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to mount device %s to location %s",
			vol.LvDmPath,
			mountpoint)
		return
	}

	resp = &v.MountResponse{
		Mountpoint: mountpoint,
	}

	return
}

func (d Driver) Unmount(req *v.UnmountRequest) (err error) {
	var (
		mountpoint string
		found      bool
	)

	d.logger.Debug().
		Str("name", req.Name).
		Str("ID", req.ID).
		Msg("starting unmount")

	d.Lock()
	defer d.Unlock()

	mountpoint, found, err = d.dirManager.Get(req.Name)
	if err != nil {
		err = errors.Errorf(
			"errored looking for mountpoint of volume %s",
			req.Name)
		return
	}

	if !found {
		return
	}

	err = d.lvm.Unmount(mountpoint)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to unmount volume %s from %s",
			req.Name, mountpoint)
		return
	}

	return
}

func (d Driver) Capabilities() (resp *v.CapabilitiesResponse) {
	d.logger.Debug().
		Msg("starting capabilities")

	resp = &v.CapabilitiesResponse{
		Capabilities: v.Capability{
			Scope: "global",
		},
	}

	return
}
