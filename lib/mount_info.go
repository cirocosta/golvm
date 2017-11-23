package lib

import (
	"strings"

	"github.com/pkg/errors"
)

// ParseMountLine parses a line from /proc/mounts
// and returns a *MountInfo struct.
func ParseMountLine(line string) (info *MountInfo, err error) {
	var parts []string

	if line == "" {
		err = errors.Errorf("can't parse empty line")
		return
	}

	parts = strings.Fields(line)
	if len(parts) < 4 {
		err = errors.Errorf("not enough info in mounts line")
		return
	}

	info = &MountInfo{
		Device:   parts[0],
		Location: parts[1],
		Format:   parts[2],
		Options:  parts[3],
	}
	return
}
