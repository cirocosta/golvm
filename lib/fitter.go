package lib

import (
	"github.com/pkg/errors"
)

// PickBestVolumeGroup picks the volume group that best
// accomodates space of a given size.
// 'size' specifies the size to be accomodated - if 0, any
// volume with free space fits it.
func PickBestVolumeGroup(size float64, vols []*VolumeGroup) (bestVol *VolumeGroup, err error) {
	if vols == nil {
		err = errors.Errorf("can't pick best volume from nil list of vols")
		return
	}

	for _, vol := range vols {
		if vol.Free > size {
			if bestVol == nil {
				bestVol = vol
				continue
			}

			if bestVol.Free > vol.Free {
				continue
			}

			bestVol = vol
		}
	}

	return
}
