package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseMountInfoLine(t *testing.T) {
	var testCases = []struct {
		desc        string
		input       string
		expected    *MountInfo
		shouldError bool
	}{
		{
			desc:        "empty should err",
			input:       "",
			expected:    nil,
			shouldError: true,
		},
		{
			desc:        "too few should err",
			input:       "ahuah auh aaa",
			expected:    nil,
			shouldError: true,
		},
		{
			desc:  "parse accordingly w/ virtual device",
			input: "udev /proc/timer_stats devtmpfs rw,nosuid,mode=755 0 0",
			expected: &MountInfo{
				Device:   "udev",
				Location: "/proc/timer_stats",
				Format:   "devtmpfs",
				Options:  "rw,nosuid,mode=755",
			},
			shouldError: false,
		},
		{
			desc:  "parse accordingly w/ regular device",
			input: "/dev/mapper/volgroup2-abc /mnt/abc ext4 rw,relatime,data=ordered 0 0",
			expected: &MountInfo{
				Device:   "/dev/mapper/volgroup2-abc",
				Location: "/mnt/abc",
				Format:   "ext4",
				Options:  "rw,relatime,data=ordered",
			},
			shouldError: false,
		},
		{
			desc:  "parse accordingly w/ regula and leading spaces",
			input: "\t\n    /device /location fmt opts 0 0\t   ",
			expected: &MountInfo{
				Device:   "/device",
				Location: "/location",
				Format:   "fmt",
				Options:  "opts",
			},
			shouldError: false,
		},
		{
			desc:  "parse accordingly w/ any space between parameters",
			input: "\t\n    /device\t\n     /location \t fmt opts 0 0\t   ",
			expected: &MountInfo{
				Device:   "/device",
				Location: "/location",
				Format:   "fmt",
				Options:  "opts",
			},
			shouldError: false,
		},
	}

	var (
		err    error
		actual *MountInfo
	)

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actual, err = ParseMountLine(tc.input)
			if tc.shouldError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, actual)

			assert.Equal(t, tc.expected.Device, actual.Device)
			assert.Equal(t, tc.expected.Location, actual.Location)
			assert.Equal(t, tc.expected.Format, actual.Format)
			assert.Equal(t, tc.expected.Options, actual.Options)
		})
	}
}
