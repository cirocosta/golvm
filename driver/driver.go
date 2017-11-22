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
	vgWhiteList []string
	sync.Mutex
}

type DriverConfig struct {
	Lvm             *lib.Lvm
	DirManager      *DirManager
	VgWhitelistFile string
}

func NewDriver(cfg DriverConfig) (d Driver, err error) {
	var whitelist []string

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

	for _, vg := range whitelist {
		d.logger.Info().
			Str("vg", vg).
			Msg("vg whitelisted")
	}

	d.lvm = cfg.Lvm
	d.dirManager = cfg.DirManager
	d.vgWhiteList = whitelist
	d.logger.Info().Msg("driver initialized")

	return
}

// Create creates a new logical volume if it doesn't
// exist yet.
func (d Driver) Create(req *v.CreateRequest) (err error) {
	var (
		size        string
		thinpool    string
		snapshot    string
		keyfile     string
		volumegroup string
		fstype      string
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
	fstype, _ = req.Options["fstype"]

	config := &lib.LvCreationConfig{
		Name:        req.Name,
		Size:        size,
		ThinPool:    thinpool,
		Snapshot:    snapshot,
		KeyFile:     keyfile,
		VolumeGroup: volumegroup,
		FsType:      fstype,
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

	resp = &v.GetResponse{
		Volume: &v.Volume{
			Name: req.Name,
		},
	}

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
