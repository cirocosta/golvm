package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	respPvs1 = `
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

	respVgs1 = `
{
    "report": [
        {
            "vg": [
                {
                    "lv_count": "0",
                    "pv_count": "1",
                    "snap_count": "0",
                    "vg_attr": "wz--n-",
                    "vg_free": "48.00",
                    "vg_name": "myvg",
                    "vg_size": "48.00"
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
		expected    []*PhysicalVolume
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
			expected:    []*PhysicalVolume{},
			shouldError: true,
		},
		{
			desc:  "valid response should return valid",
			input: []byte(respPvs1),
			expected: []*PhysicalVolume{
				&PhysicalVolume{
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
		infos  []*PhysicalVolume
		actual *PhysicalVolume
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

func TestParseVolumeGroupsOutput(t *testing.T) {
	var testCases = []struct {
		desc        string
		input       []byte
		expected    []*VolumeGroup
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
			expected:    []*VolumeGroup{},
			shouldError: true,
		},
		{
			desc:  "valid response should return valid",
			input: []byte(respVgs1),
			expected: []*VolumeGroup{
				&VolumeGroup{
					Attr:      "wz--n-",
					Name:      "myvg",
					Free:      48,
					Size:      48,
					LvCount:   0,
					PvCount:   1,
					SnapCount: 0,
				},
			},
			shouldError: false,
		},
	}

	var (
		err    error
		infos  []*VolumeGroup
		actual *VolumeGroup
	)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			infos, err = DecodeVolumeGroupsResponse(tc.input)
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
					expectedInfo.Name,
					actual.Name)
				assert.Equal(t,
					expectedInfo.Free,
					actual.Free)
				assert.Equal(t,
					expectedInfo.Size,
					actual.Size)
				assert.Equal(t,
					expectedInfo.LvCount,
					actual.LvCount)
				assert.Equal(t,
					expectedInfo.PvCount,
					actual.PvCount)
				assert.Equal(t,
					expectedInfo.SnapCount,
					actual.SnapCount)
			}

		})
	}
}
