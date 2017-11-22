package driver

import (
	"bufio"
	"os"

	"github.com/pkg/errors"
)

func ReadVgWhitelist(filename string) (vgs []string, err error) {
	var (
		file *os.File
	)

	file, err = os.Open(filename)
	if err != nil {
		err = errors.Wrapf(err,
			"can't open whitelist file %s", filename)
		return
	}
	defer file.Close()

	vgs = make([]string, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		vgs = append(vgs, scanner.Text())
	}

	err = scanner.Err()
	if err != nil {
		err = errors.Wrapf(err,
			"failed to read whitelist %s lines",
			filename)
		return
	}

	return
}
