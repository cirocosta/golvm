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
	HostMountPoint = "/mnt/lvmvol"
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

	d.lvm = cfg.Lvm
	d.logger = zerolog.New(os.Stdout).
		With().
		Str("from", "driver").
		Logger()

	return
}

func (d Driver) Create(req *v.CreateRequest) (err error) {
	d.Lock()
	defer d.Unlock()

	d.logger.Debug().
		Str("name", req.Name).
		Interface("opts", req.Options).
		Msg("starting creation")

	return
}

func (d Driver) List() (resp *v.ListResponse, err error) {
	d.Lock()
	defer d.Unlock()

	//        resp.Volumes = make([]*v.Volume, len(dirs))
	//        for idx, dir := range dirs {
	//        	resp.Volumes[idx] = &v.Volume{
	//        		Name: dir,
	//        	}
	//        }

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
	resp.Capabilities = v.Capability{
		Scope: "global",
	}
	return
}
