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

type LvAttr struct {
	IsThin bool
}
