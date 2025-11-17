package fdc

type SectorSize uint8

func (s SectorSize) Bytes() int {
	return 256 * int(s)
}

type DiskType struct {
	NumberOfHeads      int
	NumberOfCylinders  int
	SectorsPerCylinder int
	SectorSize         SectorSize
}

func (diskType DiskType) CHRToOffset(cylinder, head, sector int) int {
	return (sector - 1) + (head * diskType.SectorsPerCylinder) + cylinder*diskType.NumberOfHeads*diskType.SectorsPerCylinder
}

var Floppy144DiskType = DiskType{
	NumberOfHeads:      2,
	NumberOfCylinders:  80,
	SectorsPerCylinder: 18,
	SectorSize:         0x02,
}
