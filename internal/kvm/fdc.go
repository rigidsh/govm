package kvm

func CreateFloppy144Drive() (*FDD, error) {
	return &FDD{
		image:            []byte{},
		sectorSize:       512,
		sectorsPerTrack:  18,
		headsPerCylinder: 2,
	}, nil
}

type FDD struct {
	image            []byte
	sectorSize       uint16
	sectorsPerTrack  uint16
	headsPerCylinder uint16

	currentCylinder uint16
}

type DiskHeadSelector uint8

func (d DiskHeadSelector) Disk() uint8 {
	return uint8(d & 0b00000011)
}

func (d DiskHeadSelector) Head() uint8 {
	return uint8(d&0b00001100) >> 2
}

//func (image *FDD) ReadCHS(c, h, s uint16) ([]byte, error) {
//
//}

func (image *FDD) Recalibrate() {
	image.currentCylinder = 0
}

func (image *FDD) Seek(cylinder uint16) {
	image.currentCylinder = cylinder
}

func (image *FDD) ReadData() {
}

func CreateFDC(vm *VM) (*FDC, error) {
	fdc := &FDC{
		vm:    vm,
		disks: make([]*FDD, 0),
	}

	return fdc, nil
}

type FDC struct {
	vm    *VM
	disks []*FDD
}

func (fdc *FDC) Recalibrate(selector DiskHeadSelector) {
	fdc.disks[selector.Disk()].Recalibrate()
}

func (fdc *FDC) Seek(selector DiskHeadSelector, cylinder uint16) {
	fdc.disks[selector.Disk()].Seek(cylinder)
}
