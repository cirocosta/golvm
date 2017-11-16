package utils

import (
	"os"
	"fmt"
)

func Abort(err error) {
	if err == nil {
		return
	}

	fmt.Printf("ERRORED: %+v\nAborting.\n", err)
	os.Exit(1)
}
