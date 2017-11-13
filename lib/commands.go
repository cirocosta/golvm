package lib

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type Lvm struct {
	logger zerolog.Logger
}

func NewLvm(cfg LvmConfig) (l Lvm, err error) {
	l.logger = zerolog.New(os.Stdout).
		With().
		Str("from", "lvm").
		Logger()

	return
}

func (repo Lvm) GetPhysicalVolumes() (output string, err error) {
	output, err = repo.runCommand("pvs", "--units=m", "--separator=:", "--nosuffix", "--noheadings")
	return
}

func (repo Lvm) VolumeGroups() (output string, err error) {
	output, err = repo.runCommand("vgs", "--units=m", "--separator=:", "--nosuffix", "--noheadings")
	return
}

func (repo Lvm) LogicalVolumes() (output string, err error) {
	output, err = repo.runCommand("lvs", "--units=m", "--separator=:", "--nosuffix", "--noheadings")
	return
}

func (repo Lvm) run(name string, args ...string) (output string, err error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.Output()
	output = fmt.Sprintf("%s", out)
	if err != nil {
		return
	}
	return
}
