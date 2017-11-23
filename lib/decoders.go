package lib

import (
	"encoding/json"

	"github.com/pkg/errors"
)

// DecodePhysicalVolumesReponse takes a JSON response from
// the execition of the 'pvs' command and returns a slice of
// PhysicalVolume structs.
func DecodePhysicalVolumesResponse(response []byte) (infos []*PhysicalVolume, err error) {
	if response == nil {
		err = errors.Errorf("response can't be nil")
		return
	}

	if len(response) == 0 {
		err = errors.Errorf("can't decode empty response")
		return
	}

	var report = new(PhysicalVolumesReport)
	err = json.Unmarshal(response, report)
	if err != nil {
		err = errors.Wrapf(err, "errored decoding pvs response")
		return
	}

	if len(report.Report) != 1 {
		err = errors.Errorf(
			"unexpected number of responses decoded - %s",
			response)
		return
	}

	infos = report.Report[0].Pv
	return
}

// DecodeVolumeGroupsRepsponse takes a JSON response from
// the execition of the 'vgs' command and returns a slice of
// VolumeGroup structs.
func DecodeVolumeGroupsResponse(response []byte) (infos []*VolumeGroup, err error) {
	if response == nil {
		err = errors.Errorf("response can't be nil")
		return
	}

	if len(response) == 0 {
		err = errors.Errorf("can't decode empty response")
		return
	}

	var report = new(VolumeGroupsReport)
	err = json.Unmarshal(response, report)
	if err != nil {
		err = errors.Wrapf(err, "errored decoding vgs response")
		return
	}

	if len(report.Report) != 1 {
		err = errors.Errorf(
			"unexpected number of responses decoded - %s",
			response)
		return
	}

	infos = report.Report[0].Vg
	return
}

// DecodeLogicalVolumesResponse takes a JSON response from
// the execition of the 'lvs' command and returns a slice of
// VolumeGroup structs.
func DecodeLogicalVolumesResponse(response []byte) (infos []*LogicalVolume, err error) {
	if response == nil {
		err = errors.Errorf("response can't be nil")
		return
	}

	if len(response) == 0 {
		err = errors.Errorf("can't decode empty response")
		return
	}

	var report = new(LogicalVolumesReport)
	err = json.Unmarshal(response, report)
	if err != nil {
		err = errors.Wrapf(err, "errored decoding lvs response")
		return
	}

	if len(report.Report) != 1 {
		err = errors.Errorf(
			"unexpected number of responses decoded - %s",
			response)
		return
	}

	infos = report.Report[0].Lv
	return
}
