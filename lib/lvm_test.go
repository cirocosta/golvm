package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)


func TestBuildGetDeviceFormatArgs(t *testing.T) {
	var testCases = []struct {
		desc        string
		device      string
		expected    []string
		shouldError bool
	}{
		{
			desc:        "fail with empty device",
			device:      "",
			expected:    []string{},
			shouldError: true,
		},
		{
			desc:   "works with device",
			device: "/dev/device",
			expected: []string{
				"--noheadings",
				"--discard",
				"--output=FSTYPE",
				"/dev/device",
			},
			shouldError: false,
		},
	}

	var (
		err  error
		args []string
	)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			args, err = BuildGetDeviceFormatArgs(tc.device)
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

func TestBuildMakeFsArgs(t *testing.T) {
	var testCases = []struct {
		desc        string
		fsType      string
		device      string
		expected    []string
		shouldError bool
	}{
		{
			desc:        "fail with empty fstype",
			fsType:      "",
			device:      "",
			expected:    []string{},
			shouldError: true,
		},
		{
			desc:        "fail with empty device",
			fsType:      "ext4",
			device:      "",
			expected:    []string{},
			shouldError: true,
		},
		{
			desc:        "fail with unknown fstype",
			fsType:      "unknown",
			device:      "/dev/device",
			expected:    []string{},
			shouldError: true,
		},
		{
			desc:   "works with ext4",
			fsType: "ext4",
			device: "/dev/device",
			expected: []string{
				"-t",
				"ext4",
				"/dev/device",
			},
			shouldError: false,
		},
		{
			desc:   "works with xfs",
			fsType: "xfs",
			device: "/dev/device",
			expected: []string{
				"-t",
				"xfs",
				"/dev/device",
			},
			shouldError: false,
		},
	}

	var (
		err  error
		args []string
	)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			args, err = BuildMakeFsArgs(tc.fsType, tc.device)
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
		args []string
	)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			args, err = BuildLogicalVolumeRemovalArgs(tc.cfg)
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
			desc: "fails with only name and vg",
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
		{
			desc: "thin vol works with name, thinpool, size and vg",
			cfg: &LvCreationConfig{
				Name:        "name",
				VolumeGroup: "volumegroup",
				Size:        "22M",
				ThinPool:    "tp",
			},
			expected: []string{
				"--setactivationskip", "n",
				"--name", "name",
				"--virtualsize", "22M",
				"--thin",
				"volumegroup/tp",
			},
			shouldError: false,
		},
	}

	var (
		err  error
		args []string
	)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			args, err = BuildLogicalVolumeCretionArgs(tc.cfg)
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
		{
			desc:  "proper thin parse",
			input: "Vwi-a-tz--",
			expected: &LvAttr{
				VolumeType:                "thin volume",
				Permissions:               "writeable",
				AllocationPolicy:          "inherited",
				FixedMinor:                "-",
				State:                     "active",
				DeviceState:               "-",
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
