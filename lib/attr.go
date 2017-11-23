package lib

import (
	"github.com/pkg/errors"
	"strings"
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
