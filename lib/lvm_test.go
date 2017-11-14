package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	respVgs1 = `
{
    "report": [
        {
            "pv": [
                {
                    "pv_attr": "a--",
                    "pv_fmt": "lvm2",
                    "pv_free": "48.00",
                    "pv_name": "/dev/loop0",
                    "pv_size": "48.00",
                    "vg_name": "myvg"
                }
            ]
        }
    ]
}`
)

func TestParsePhysycalVolumesOutput(t *testing.T) {
	var testCases = []struct {
		desc        string
		input       []byte
		expected    []*PhysicalInfo
		shouldError bool
	}{
		{
			desc:        "nil input should fail",
			input:       nil,
			shouldError: true,
		},
		{
			desc:        "empty input should fail",
			input:       []byte(""),
			expected:    []*PhysicalInfo{},
			shouldError: true,
		},
		{
			desc:  "valid response should return valid",
			input: []byte(respVgs1),
			expected: []*PhysicalInfo{
				&PhysicalInfo{
					Attr:             "a--",
					Fmt:              "lvm2",
					PhysicalSize:     48,
					PhysicalVolume:   "/dev/loop0",
					PhysicalSizeFree: 48,
					VolumeGroup:      "myvg",
				},
			},
			shouldError: false,
		},
	}

	var (
		err    error
		infos  []*PhysicalInfo
		actual *PhysicalInfo
	)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			infos, err = DecodePhysicalVolumesResponse(tc.input)
			if tc.shouldError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, len(tc.expected), len(infos))

			for ndx, expectedInfo := range tc.expected {
				actual = infos[ndx]

				assert.Equal(t,
					expectedInfo.Attr,
					actual.Attr)
				assert.Equal(t,
					expectedInfo.VolumeGroup,
					actual.VolumeGroup)
				assert.Equal(t,
					expectedInfo.Fmt,
					actual.Fmt)
				assert.Equal(t,
					expectedInfo.PhysicalSize,
					actual.PhysicalSize)
				assert.Equal(t,
					expectedInfo.PhysicalSizeFree,
					actual.PhysicalSizeFree)
			}

		})
	}
}
