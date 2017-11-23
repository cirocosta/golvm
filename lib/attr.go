package lib

import (
	"github.com/pkg/errors"
	"strings"
)

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

// ParseLvAttr takes an 'attr' string from 'lvs' command and
// parses it so that it can be consumed via the LvAttr struct.
// In case of any unexpected tokens or malformed attr, fails
// with an error.
func ParseLvAttr(attr string) (parsedAttr *LvAttr, err error) {
	if attr == "" {
		err = errors.Errorf("attr must not be empty")
		return
	}

	if len(attr) != 10 {
		err = errors.Errorf(
			"malformed attr '%s' - must be 10chars",
			attr)
		return
	}

	var (
		mapping map[string]string
		present bool
		val     string
		chars   = strings.Split(attr, "")
	)

	parsedAttr = new(LvAttr)
	for ndx, character := range chars {
		mapping = lvAttrMapper[ndx]
		val, present = mapping[character]
		if !present {
			err = errors.Errorf(
				"unexpected character '%' for lv attr '%d'",
				character, ndx)
			return
		}

		switch ndx {
		case 0:
			parsedAttr.VolumeType = val
		case 1:
			parsedAttr.Permissions = val
		case 2:
			parsedAttr.AllocationPolicy = val
		case 3:
			parsedAttr.FixedMinor = val
		case 4:
			parsedAttr.State = val
		case 5:
			parsedAttr.DeviceState = val
		case 6:
			parsedAttr.TargetType = val
		case 7:
			parsedAttr.OverrideNewBlocksWithZero = val
		case 8:
			parsedAttr.VolumeHealth = val
		case 9:
			parsedAttr.SkipActivation = val
		}
	}

	return
}
