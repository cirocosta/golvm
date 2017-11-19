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
	respLvs1 = `
{
    "report": [
        {
            "lv": [
                {
                    "convert_lv": "",
                    "copy_percent": "",
                    "data_percent": "",
                    "lv_attr": "-wi-a-----",
                    "lv_name": "lv1",
                    "lv_size": "12.00",
                    "metadata_percent": "",
                    "mirror_log": "",
                    "move_pv": "",
                    "origin": "",
                    "pool_lv": "",
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

func TestParseLogicalVolumesOutput(t *testing.T) {
	var testCases = []struct {
		desc        string
		input       []byte
		expected    []*LogicalVolume
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
			expected:    []*LogicalVolume{},
			shouldError: true,
		},
		{
			desc:  "valid response should return valid",
			input: []byte(respLvs1),
			expected: []*LogicalVolume{
				&LogicalVolume{
					LvName: "lv1",
					LvSize: 12.0,
					LvAttr: "-wi-a-----",
				},
			},
			shouldError: false,
		},
	}

	var (
		err    error
		infos  []*LogicalVolume
		actual *LogicalVolume
	)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			infos, err = DecodeLogicalVolumesResponse(tc.input)
			if tc.shouldError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, len(tc.expected), len(infos))

			for ndx, expectedInfo := range tc.expected {
				actual = infos[ndx]

				assert.Equal(t,
					expectedInfo.LvName,
					actual.LvName)
				assert.Equal(t,
					expectedInfo.LvSize,
					actual.LvSize)
				assert.Equal(t,
					expectedInfo.LvAttr,
					actual.LvAttr)
			}

		})
	}
}

func TestBuildLogicalVolumeRemovalArgs(t *testing.T) {
	var testCases = []struct {
		desc        string
		cfg         LvRemovalConfig
		expected    []string
		shouldError bool
	}{
		{
			desc:        "without a lvName should fail",
			cfg:         LvRemovalConfig{},
			expected:    []string{},
			shouldError: true,
		},
		{
			desc: "without a vgName should fail",
			cfg: LvRemovalConfig{
				LvName: "lv",
			},
			expected:    []string{},
			shouldError: true,
		},
		{
			desc: "works with vg and lv names",
			cfg: LvRemovalConfig{
				LvName: "lv",
				VgName: "vg",
			},
			expected: []string{
				"--force", "vg/lv",
			},
			shouldError: false,
		},
	}

	var (
		err  error
		l    Lvm
		args []string
	)

	l, err = NewLvm(LvmConfig{})
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			args, err = l.BuildLogicalVolumeRemovalArgs(tc.cfg)
			if tc.shouldError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, len(tc.expected), len(args))
			for ndx, arg := range args {
				assert.Equal(t, tc.expected[ndx], arg)
			}
		})
	}

}

func TestBuildLogicalVolumeCreationArgs(t *testing.T) {
	var testCases = []struct {
		desc        string
		cfg         *LvCreationConfig
		expected    []string
		shouldError bool
	}{
		{
			desc:        "without a name should fail",
			cfg:         &LvCreationConfig{},
			expected:    []string{},
			shouldError: true,
		},
		{
			desc: "without a vg should fail",
			cfg: &LvCreationConfig{
				Name: "haha",
			},
			expected:    []string{},
			shouldError: true,
		},
		{
			desc: "fails with only name and vg ",
			cfg: &LvCreationConfig{
				Name:        "name",
				VolumeGroup: "volumegroup",
			},
			expected:    []string{},
			shouldError: true,
		},
		{
			desc: "snap fails with only name and vg",
			cfg: &LvCreationConfig{
				Name:        "name",
				VolumeGroup: "volumegroup",
				Snapshot:    "snapshot",
			},
			expected:    []string{},
			shouldError: true,
		},
		{
			desc: "snap fails with only name and vg",
			cfg: &LvCreationConfig{
				Name:        "name",
				VolumeGroup: "volumegroup",
				Snapshot:    "snapshot",
			},
			expected:    []string{},
			shouldError: true,
		},
		{
			desc: "vol works with name, size and vg",
			cfg: &LvCreationConfig{
				Name:        "name",
				VolumeGroup: "volumegroup",
				Size:        "22M",
			},
			expected: []string{
				"--setactivationskip", "n",
				"--name", "name",
				"--size", "22M",
				"volumegroup",
			},
			shouldError: false,
		},
	}

	var (
		err  error
		l    Lvm
		args []string
	)

	l, err = NewLvm(LvmConfig{})
	require.NoError(t, err)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			args, err = l.BuildLogicalVolumeCretionArgs(tc.cfg)
			if tc.shouldError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, len(tc.expected), len(args))
			for ndx, arg := range args {
				assert.Equal(t, tc.expected[ndx], arg)
			}
		})
	}

}

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

func TestParseLvAttr(t *testing.T) {
	var testCases = []struct {
		desc        string
		input       string
		expected    *LvAttr
		shouldError bool
	}{
		{
			desc:        "fails with empty",
			input:       "",
			expected:    nil,
			shouldError: true,
		},
		{
			desc:        "fails with not enough chars",
			input:       "aa",
			expected:    nil,
			shouldError: true,
		},
		{
			desc:        "fails with more than enough chars",
			input:       "aaaaaaaaaaaasssssaaaaaa",
			expected:    nil,
			shouldError: true,
		},
		{
			desc:        "fails with unexpected chars",
			input:       "Ç---------",
			expected:    &LvAttr{},
			shouldError: true,
		},
		{
			desc:        "fails with unexpected chars",
			input:       "a         ",
			expected:    &LvAttr{},
			shouldError: true,
		},
		{
			desc:  "proper volume type parsing",
			input: "t---------",
			expected: &LvAttr{
				VolumeType:                "thin pool",
				Permissions:               "-",
				AllocationPolicy:          "-",
				FixedMinor:                "-",
				State:                     "-",
				DeviceState:               "-",
				TargetType:                "-",
				OverrideNewBlocksWithZero: "-",
				VolumeHealth:              "-",
				SkipActivation:            "-",
			},
			shouldError: false,
		},
		{
			desc:  "proper overall parsing",
			input: "twi-aotz--",
			expected: &LvAttr{
				VolumeType:                "thin pool",
				Permissions:               "writeable",
				AllocationPolicy:          "inherited",
				FixedMinor:                "-",
				State:                     "active",
				DeviceState:               "open",
				TargetType:                "thin",
				OverrideNewBlocksWithZero: "overwrite by zero",
				VolumeHealth:              "-",
				SkipActivation:            "-",
			},
			shouldError: false,
		},
	}

	var (
		parsed *LvAttr
		err    error
	)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			parsed, err = ParseLvAttr(tc.input)
			if tc.shouldError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, parsed)

			assert.Equal(t, tc.expected.VolumeType,
				parsed.VolumeType)
			assert.Equal(t, tc.expected.Permissions,
				parsed.Permissions)
			assert.Equal(t, tc.expected.AllocationPolicy,
				parsed.AllocationPolicy)
			assert.Equal(t, tc.expected.FixedMinor,
				parsed.FixedMinor)
			assert.Equal(t, tc.expected.State,
				parsed.State)
			assert.Equal(t, tc.expected.DeviceState,
				parsed.DeviceState)
			assert.Equal(t, tc.expected.TargetType,
				parsed.TargetType)
			assert.Equal(t, tc.expected.OverrideNewBlocksWithZero,
				parsed.OverrideNewBlocksWithZero)
			assert.Equal(t, tc.expected.VolumeHealth,
				parsed.VolumeHealth)
			assert.Equal(t, tc.expected.SkipActivation,
				parsed.SkipActivation)
		})
	}
}
