package driver

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
	"golang.org/x/sys/unix"
)

type DirManager struct {
	root string
}

type DirManagerConfig struct {
	Root string
}

var (
	NameRegex = regexp.MustCompile(`^[a-zA-Z0-9][\w\-]{1,250}$`)

	ErrInvalidName = errors.Errorf("Invalid name")
	ErrNotFound    = errors.Errorf("Volume not found")
)

func NewDirManager(cfg DirManagerConfig) (manager DirManager, err error) {
	if cfg.Root == "" {
		err = errors.Errorf("A root must be specified.")
		return
	}

	if !filepath.IsAbs(cfg.Root) {
		err = errors.Errorf(
			"Root (%s) must be an absolute path",
			cfg.Root)
		return
	}

	err = unix.Access(cfg.Root, unix.W_OK)
	if err != nil {
		err = errors.Wrapf(err,
			"Root (%s) must be writable.",
			cfg.Root)
		return
	}

	manager.root = cfg.Root
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
	if !isValidName(path) {
		err = ErrInvalidName
		return
	}

	absPath = filepath.Join(m.root, path)
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