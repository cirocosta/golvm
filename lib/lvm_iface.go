package lib

type LvmIface interface {
	PhysicalVolumes() (output []byte, err error)
	VolumeGroups() (output []byte, err error)
	LogicalVolumes() (output []byte, err error)
}
