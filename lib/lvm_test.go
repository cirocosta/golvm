package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePhysycalVolumesOutput(t *testing.T) {
	var testCases = []struct {
		desc        string
		input       string
		expected    *PhysicalInfo
		shouldError bool
	}{
		{
			desc:        "empty",
			input:       "",
			shouldError: true,
		},
		{
			desc:        "malformed",
			input:       "uhsduah",
			shouldError: true,
		},
		{
			desc:  "valid with leading space",
			input: "  /dev/loop0::lvm2:---:50.00:50.00",
			expected: &PhysicalInfo{
				PhysicalVolume:   "/dev/loop0",
				VolumeGroup:      "",
				Fmt:              "lvm2",
				Attr:             "---",
				PhysicalSize:     50,
				PhysicalFreeSize: 50,
			},
			shouldError: false,
		},
	}

	var (
		err    error
		actual *PhysicalInfo
	)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual, err = ParsePhysicalVolumesResponseLine(tc.input)
			if tc.shouldError {
				require.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t,
				tc.expected.Attr,
				actual.Attr)
			assert.Equal(t,
				tc.expected.VolumeGroup,
				actual.VolumeGroup)
			assert.Equal(t,
				tc.expected.Fmt,
				actual.Fmt)
			assert.Equal(t,
				tc.expected.PhysicalSize,
				actual.PhysicalSize)
			assert.Equal(t,
				tc.expected.PhysicalFreeSize,
				actual.PhysicalFreeSize)
		})
	}
}
