package lib

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

type LogicalVolumesReport struct{}
