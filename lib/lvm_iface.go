package lib

type LvmIface interface {
	PhysicalVolumes() (output string err error)
	VolumeGroups() (output string,   err error)
	LogicalVolumes() (output string, err error)
}

