package driver

import (
	"bufio"
	"os"

	"github.com/pkg/errors"
)

func ReadVgWhitelist(filename string) (vgs map[string]bool, err error) {
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

	vgs = make(map[string]bool)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		vgs[scanner.Text()] = true
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
