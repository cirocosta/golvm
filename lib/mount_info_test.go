package lib

import (
	"io/ioutil"
	"os"
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

func TestParseMountsFile(t *testing.T) {
	var fileContent = []byte(`
proc /proc proc rw,nosuid,nodev,noexec,relatime 0 0
udev /dev devtmpfs rw,nosuid,relatime,size=4014860k,nr_inodes=1003715,mode=755 0 0
tmpfs /sys/fs/cgroup tmpfs ro,nosuid,nodev,noexec,mode=755 0 0
cgroup /sys/fs/cgroup/systemd cgroup rw,nosuid,nodev,noexec,relatime,xattr,release_agent=/lib/systemd/systemd-cgroups-agent,name=systemd 0 0
pstore /sys/fs/pstore pstore rw,nosuid,nodev,noexec,relatime 0 0
efivarfs /sys/firmware/efi/efivars efivarfs rw,nosuid,nodev,noexec,relatime 0 0
`)
	file, err := ioutil.TempFile("", "")
	assert.NoError(t, err)
	filename := file.Name()

	defer os.Remove(filename)

	_, err = file.Write(fileContent)
	assert.NoError(t, err)
	file.Close()

	infos, err := ParseMountsFile(filename)
	assert.NoError(t, err)
	assert.Equal(t, 6, len(infos))
	assert.Equal(t, "proc", infos[0].Device)
	assert.Equal(t, "efivarfs", infos[5].Device)
}
