package lib

import (
	"github.com/rs/zerolog"
)

// Lvm encapsulates a series of methods for
// dealing with LVM management.
// It's mostly stateless except for a logger.
type Lvm struct {
	logger zerolog.Logger
}

// LvmConfig provides the configuration details for
// the Lvm helper.
type LvmConfig struct{}

// LvCreationConfig is a simplified configuration
// struct to be passed to logical volume creation
// methods.
type LvCreationConfig struct {
	Name        string
	Size        string
	Snapshot    string
	KeyFile     string
	ThinPool    string
	VolumeGroup string
	FsType      string
}

type PhysicalVolumesReport struct {
	Report []struct {
		Pv []*PhysicalVolume `json:"pv"`
	} `json:"report"`
}

type VolumeGroupsReport struct {
	Report []struct {
		Vg []*VolumeGroup `json:"vg"`
	} `json:"report"`
}

type LogicalVolumesReport struct {
	Report []struct {
		Lv []*LogicalVolume `json:"lv"`
	} `json:"report"`
}

type PhysicalVolume struct {
	PhysicalVolume   string  `json:"pv_name"`
	VolumeGroup      string  `json:"vg_name"`
	Attr             string  `json:"pv_attr"`
	Fmt              string  `json:"pv_fmt"`
	PhysicalSize     float64 `json:"pv_size,string"`
	PhysicalSizeFree float64 `json:"pv_free,string"`
}

type VolumeGroup struct {
	Attr      string  `json:"vg_attr"`
	Name      string  `json:"vg_name"`
	Free      float64 `json:"vg_free,string"`
	Size      float64 `json:"vg_size,string"`
	LvCount   uint64  `json:"lv_count,string"`
	PvCount   uint64  `json:"pv_count,string"`
	SnapCount uint64  `json:"snap_count,string"`
}

type LogicalVolume struct {
	ConvertLv       string  `json:"convert_lv"`
	CopyPercent     string  `json:"copy_percent"`
	DataPercent     string  `json:"data_percent"`
	LvAttr          string  `json:"lv_attr"`
	LvName          string  `json:"lv_name"`
	LvSize          float64 `json:"lv_size,string"`
	MetadataPercent string  `json:"metadata_percent"`
	MirrorLog       string  `json:"mirror_log"`
	MovePv          string  `json:"move_pv"`
	Origin          string  `json:"origin"`
	PoolLv          string  `json:"pool_lv"`
	VgName          string  `json:"vg_name"`
}

var (
	lvAttrVolumeTypeMap = map[string]string{
		"-": "-",
		"C": "cache",
		"m": "mirrored",
		"M": "mirrored without initial sync",
		"o": "origin",
		"O": "origin without merging snapshot",
		"r": "raid",
		"R": "raid without initial sync",
		"s": "snapshot merging",
		"S": "snapshot",
		"p": "pvmove",
		"v": "virtual",
		"V": "thin volume",
		"i": "mirror or raid image",
		"I": "mirror our raid image out-of-sync",
		"t": "thin pool",
		"T": "thin poool data",
		"e": "metadata or pool metadata sparse",
	}
	lvAttrPermissionsMap = map[string]string{
		"-": "-",
		"w": "writeable",
		"r": "read-only",
		"R": "read-only activation of non-read-only volume",
	}

	lvAttrAllocationPolicyMap = map[string]string{
		"-": "-",
		"a": "anywhere",
		"c": "contiguous",
		"i": "inherited",
		"l": "cling",
		"n": "normal",
	}

	lvAttrFixedMinorMap = map[string]string{
		"-": "-",
		"m": "fixed minor",
	}

	lvAttrStateMap = map[string]string{
		"-": "-",
		"a": "active",
		"h": "historical",
		"s": "suspended",
		"I": "invalid snapshot",
		"S": "suspended snapshot",
		"m": "snapshot merge failed",
		"M": "suspended snapshot merge failed",
		"d": "device present without tables",
		"i": "mapped device present with inactive table",
		"c": "thin-pool check needed",
		"C": "suspended thin-pool check needed",
		"X": "unknown",
	}

	lvAttrDeviceStateMap = map[string]string{
		"-": "-",
		"o": "open",
		"X": "unknown",
	}

	lvAttrTargetTypeMap = map[string]string{
		"-": "-",
		"C": "cache",
		"m": "mirror",
		"r": "raid",
		"s": "snapshot",
		"t": "thin",
		"u": "unknown",
		"v": "virtual",
	}

	lvAttrOverrideNewBlocksWithZeroMap = map[string]string{
		"-": "-",
		"z": "overwrite by zero",
	}

	lvAttrVolumeHealthMap = map[string]string{
		"-": "-",
		"p": "partial",
		"X": "unknown",
		"r": "refresh needed",
		"m": "mismatches exist",
		"w": "writemostly",
	}

	lvAttrSkipActivationMap = map[string]string{
		"-": "-",
		"k": "skip activation",
	}

	lvAttrMapper = []map[string]string{
		lvAttrVolumeTypeMap,
		lvAttrPermissionsMap,
		lvAttrAllocationPolicyMap,
		lvAttrFixedMinorMap,
		lvAttrStateMap,
		lvAttrDeviceStateMap,
		lvAttrTargetTypeMap,
		lvAttrOverrideNewBlocksWithZeroMap,
		lvAttrVolumeHealthMap,
		lvAttrSkipActivationMap,
	}
)

type LvRemovalConfig struct {
	LvName string
	VgName string
}

type LvAttr struct {
	VolumeType                string
	Permissions               string
	AllocationPolicy          string
	FixedMinor                string
	State                     string
	DeviceState               string
	TargetType                string
	OverrideNewBlocksWithZero string
	VolumeHealth              string
	SkipActivation            string
}
