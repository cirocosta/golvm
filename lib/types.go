package lib

// PhysicalInfo represents the result of the retrieval
// of a single entry from the execution of the 'vgs'
// command.
//
//	e.g:	PV:VG:Fmt:Attr:PSize:PFree
//

//
// {
//  "report": [
//    {
//      "vg": [
//	{
//	  "vg_name": "myvg",
//	  "pv_count": "1",
//	  "lv_count": "0",
//	  "snap_count": "0",
//	  "vg_attr": "wz--n-",
//	  "vg_size": "48,00",
//	  "vg_free": "48,00"
//	}
//      ]
//    }
//  ]
// }

type PhysicalVolumesReport struct {
	Report []struct {
		Pv []*PhysicalInfo `json:"pv"`
	} `json:"report"`
}

type PhysicalInfo struct {
	PhysicalVolume   string  `json:"pv_name"`
	VolumeGroup      string  `json:"vg_name"`
	Attr             string  `json:"pv_attr"`
	Fmt              string  `json:"pv_fmt"`
	PhysicalSize     float64 `json:"pv_size,string"`
	PhysicalSizeFree float64 `json:"pv_free,string"`
}

//   VG:#PV:#LV:#SN:Attr:VSize:VFree
//   myvg:1:0:0:wz--n-:48,00:48,00

type VolumeGroup struct {
	Name            string
	PhysicalVolumes int
	LogicalVolumes  int
}

type VolumeGroupsReport struct{}
type LogicalVolumesReport struct{}
