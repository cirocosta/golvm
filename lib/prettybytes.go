package lib

import (
	"github.com/pkg/errors"
	"strings"

	units "github.com/docker/go-units"
)

func MustFromHumanSize(size string) uint64 {
	bytes, err := units.FromHumanSize(size)
	if err != nil {
		err = errors.Wrapf(err,
			"Couldn't convert string in human size (size=%s) to bytes",
			size)
		panic(err)
	}

	return uint64(bytes)
}

func FromHumanSize(size string) (bytes uint64, err error) {
	size = strings.TrimSpace(size)

	bytesInt, err := units.RAMInBytes(size)
	if err != nil {
		err = errors.Wrapf(err,
			"Couldn't convert string in human size (size=%s) to bytes",
			size)
		return
	}

	bytes = uint64(bytesInt)
	return
}

func HumanSize(size uint64) string {
	return units.HumanSize(float64(size))
}
