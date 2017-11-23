package lib

import (
	"bufio"
	"os"
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

func ParseMountsFile(filename string) (infos []*MountInfo, err error) {
	var (
		line string
		file *os.File
		info *MountInfo
	)

	file, err = os.Open(filename)
	if err != nil {
		err = errors.Wrapf(err,
			"failed opening procs mount file %s",
			filename)
		return
	}
	defer file.Close()

	infos = make([]*MountInfo, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line = scanner.Text()
		if line == "" {
			continue
		}

		info, err = ParseMountLine(line)
		if err != nil {
			err = errors.Wrapf(err,
				"failed to parse line '%s' from mounts file %s",
				line, filename)
			return
		}

		infos = append(infos, info)
	}

	err = scanner.Err()
	if err != nil {
		err = errors.Wrapf(err,
			"failed to read lines from proc mounts file %s",
			filename)
		return
	}

	return
}
