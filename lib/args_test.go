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
