package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPickBestVolumeGroup(t *testing.T) {
	var testCases = []struct {
		desc        string
		size        float64
		vols        []*VolumeGroup
		expected    *VolumeGroup
		shouldError bool
	}{
		{
			desc:        "nil should error",
			size:        0,
			vols:        nil,
			expected:    nil,
			shouldError: true,
		},
		{
			desc:        "should return nil if empty list",
			size:        0,
			vols:        []*VolumeGroup{},
			expected:    nil,
			shouldError: false,
		},
		{
			desc: "should return the most free of all if 0",
			size: 0,
			vols: []*VolumeGroup{
				&VolumeGroup{
					Name: "vg1",
					Size: 30,
					Free: 10,
				},
				&VolumeGroup{
					Name: "vg2",
					Size: 30,
					Free: 20,
				},
				&VolumeGroup{
					Name: "vg3",
					Size: 30,
					Free: 5,
				},
			},
			expected: &VolumeGroup{
				Name: "vg2",
				Size: 30,
				Free: 20,
			},
			shouldError: false,
		},
		{
			desc: "should return the most free suiting size > 0",
			size: 15,
			vols: []*VolumeGroup{
				&VolumeGroup{
					Name: "vg1",
					Size: 60,
					Free: 10,
				},
				&VolumeGroup{
					Name: "vg2",
					Size: 60,
					Free: 20,
				},
				&VolumeGroup{
					Name: "vg3",
					Size: 60,
					Free: 50,
				},
			},
			expected: &VolumeGroup{
				Name: "vg3",
				Size: 60,
				Free: 50,
			},
			shouldError: false,
		},
	}

	var (
		bestVol *VolumeGroup
		err     error
	)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			bestVol, err = PickBestVolumeGroup(tc.size, tc.vols)
			if tc.shouldError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			if tc.expected != nil {
				assert.Equal(t, tc.expected.Name, bestVol.Name)
			} else {
				assert.Nil(t, bestVol)
			}
		})
	}
}
