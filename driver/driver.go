package driver

import (
	"os"
	"sync"

	"github.com/cirocosta/golvm/lib"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	v "github.com/docker/go-plugins-helpers/volume"
)

const (
	HostMountPoint            = "/mnt/lvmvol"
	VolumeGroupsWhiteListFile = "/mnt/lvmvol/whitelist"
)

type Driver struct {
	lvm    *lib.Lvm
	logger zerolog.Logger
	sync.Mutex
}

type DriverConfig struct {
	Lvm *lib.Lvm
}

func NewDriver(cfg DriverConfig) (d Driver, err error) {
	if cfg.Lvm == nil {
		err = errors.Errorf("Lvm must be specified")
		return
	}

	d.logger = zerolog.New(os.Stdout).
		With().
		Str("from", "driver").
		Logger()

	d.lvm = cfg.Lvm
	d.logger.Info().Msg("driver initialized")

	return
}

func (d Driver) Create(req *v.CreateRequest) (err error) {
	var (
		size        string
		thinpool    string
		snapshot    string
		keyfile     string
		volumegroup string
		args        []string
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

	config := &lib.LvCreationConfig{
		Name:        req.Name,
		Size:        size,
		ThinPool:    thinpool,
		Snapshot:    snapshot,
		KeyFile:     keyfile,
		VolumeGroup: volumegroup,
	}

	args, err = d.lvm.BuildLogicalVolumeCretionArgs(config)
	if err != nil {
		err = errors.Wrapf(err, "couldn't build volume creation args")
		return
	}

	err = d.lvm.CreateLv(args...)
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
	d.Lock()
	defer d.Unlock()

	// resp.Volume = &v.Volume{
	// 	Name:       req.Name,
	// 	Mountpoint: mp,
	// }

	return
}

func (d Driver) Remove(req *v.RemoveRequest) (err error) {
	d.Lock()
	defer d.Unlock()

	return
}

func (d Driver) Path(req *v.PathRequest) (resp *v.PathResponse, err error) {
	d.Lock()
	defer d.Unlock()

	// 	resp.Mountpoint = mp

	return
}

func (d Driver) Mount(req *v.MountRequest) (resp *v.MountResponse, err error) {
	d.Lock()
	defer d.Unlock()

	return
}

func (d Driver) Unmount(req *v.UnmountRequest) (err error) {
	d.Lock()
	defer d.Unlock()

	return
}

func (d Driver) Capabilities() (resp *v.CapabilitiesResponse) {
	resp = &v.CapabilitiesResponse{
		Capabilities: v.Capability{
			Scope: "global",
		},
	}

	return
}
