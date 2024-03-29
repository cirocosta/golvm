package driver

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

// DirManager is responsible for managing directories
// under a specific root.
// It's main functionalities include:
//	-	creating directories
//	-	inspecting directories
//	-	removing directories
type DirManager struct {
	root string
}

// DirManagerConfig provides the minimum configuration for
// instantiating a DirManager.
type DirManagerConfig struct {
	Root string
}

var (
	NameRegex = regexp.MustCompile(`^[a-zA-Z0-9][\w\-]{1,250}$`)

	ErrInvalidName = errors.Errorf("Invalid name")
	ErrNotFound    = errors.Errorf("Volume not found")
)

// NewDirManager instantiates a new DirManager.
func NewDirManager(cfg DirManagerConfig) (manager DirManager, err error) {
	if cfg.Root == "" {
		err = errors.Errorf("A root must be specified.")
		return
	}

	if !filepath.IsAbs(cfg.Root) {
		err = errors.Errorf(
			"Root %s must be an absolute path",
			cfg.Root)
		return
	}

	_, err = os.Stat(cfg.Root)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(cfg.Root, 0755)
			if err != nil {
				err = errors.Wrapf(err,
					"Failed to create directory %s",
					cfg.Root)
				return
			}
		} else {
			err = errors.Wrapf(err,
				"Errored inspecting directory %s",
				cfg.Root)
			return
		}
	}

	err = unix.Access(cfg.Root, unix.W_OK)
	if err != nil {
		err = errors.Wrapf(err,
			"Root %s must be writable.",
			cfg.Root)
		return
	}

	manager.root = cfg.Root
	return
}

// Mountpoint retrieves the full mountpoint that a
// given 'name' can receive.
func (m DirManager) Mountpoint(name string) (mp string, err error) {
	if name == "" {
		err = errors.Errorf("a name must be provided")
		return
	}

	if !isValidName(name) {
		err = ErrInvalidName
		return
	}

	mp = filepath.Join(m.root, name)
	return
}

func (m DirManager) List() (directories []string, err error) {
	files, err := ioutil.ReadDir(m.root)
	if err != nil {
		err = errors.Wrapf(err,
			"Couldn't list files/directories from %s", m.root)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			directories = append(directories, file.Name())
		}
	}

	return
}

func (m DirManager) Get(name string) (absPath string, found bool, err error) {
	if !isValidName(name) {
		err = ErrInvalidName
		return
	}

	files, err := ioutil.ReadDir(m.root)
	if err != nil {
		err = errors.Wrapf(err,
			"Couldn't list files/directories from %s", m.root)
		return
	}

	for _, file := range files {
		if file.IsDir() && file.Name() == name {
			found = true
			absPath = filepath.Join(m.root, name)
			return
		}
	}

	return
}

func (m DirManager) Create(path string) (absPath string, err error) {
	absPath, err = m.Mountpoint(path)
	if err != nil {
		err = errors.Wrapf(err,
			"Couldn't retrieve full name for path %s", path)
		return
	}

	err = os.MkdirAll(absPath, 0755)
	if err != nil {
		err = errors.Wrapf(err,
			"Couldn't create directory %s", absPath)
		return
	}

	return
}

func (m DirManager) Delete(name string) (err error) {
	if !isValidName(name) {
		err = ErrInvalidName
		return
	}

	abs, found, err := m.Get(name)
	if err != nil {
		err = errors.Wrapf(err,
			"Errored retrieving abs path for name %s",
			name)
		return
	}

	if !found {
		err = ErrNotFound
		return
	}

	err = os.RemoveAll(abs)
	if err != nil {
		err = errors.Wrapf(err,
			"Errored removing volume named %s at path %s",
			name, abs)
		return
	}

	return
}

func isValidName(name string) bool {
	if name == "" {
		return false
	}

	return NameRegex.MatchString(name)
}
