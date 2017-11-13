package lib

// PhysicalInfo represents the result of the retrieval
// of a single entry from the execution of the 'vgs'
// command.
//
//	e.g:	PV:VG:Fmt:Attr:PSize:PFree
//
type PhysicalInfo struct {
	PhysicalVolume   string
	VolumeGroup      string
	Attr             string
	Fmt              string
	PhysicalSize     float64
	PhysicalFreeSize float64
}
