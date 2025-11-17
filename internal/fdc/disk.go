package fdc

import (
	"bytes"
	"io"
)

type Disk interface {
	SectorReader(cylinder uint8, head uint8, sector uint8) (SectorHeader, io.Reader)
}

type SectorHeader struct {
	Cylinder   uint8
	Head       uint8
	Sector     uint8
	SectorSize uint8
}

type InMemoryRawDisk []byte

func (disk InMemoryRawDisk) SectorReader(cylinder uint8, head uint8, sector uint8) (SectorHeader, io.Reader) {
	imageOffset := Floppy144DiskType.CHRToOffset(int(cylinder), int(head), int(sector))

	return SectorHeader{
			Cylinder:   cylinder,
			Head:       head,
			Sector:     sector,
			SectorSize: 0x02,
		},

		bytes.NewBuffer(disk[imageOffset : imageOffset+Floppy144DiskType.SectorSize.Bytes()])
}
