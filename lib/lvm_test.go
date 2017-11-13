package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParsePhysycalVolumesOutput(t *testing.T) {
	var testCases = []struct {
		desc  string
		input string
		shouldError bool
	}{
		{
			desc: "empty",
		},
	}

	var err error

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			err = ParsePhysicalVolumesResponse(tc.input)
			if tc.shouldError {
				require.Error(t, err)
				return
			}
		})
	}
}
