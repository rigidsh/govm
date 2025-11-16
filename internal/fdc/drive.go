package fdc

import (
	"bytes"
	"io"
)

type diskDrive struct {
	disk            Disk
	currentCylinder uint8
	currentHead     uint8

	sectorSize uint8
	gapLength  uint8
	dataLength uint8
}

func (drive *diskDrive) setSettings(sectorSize, gapLength, dataLength uint8) {
	drive.sectorSize = sectorSize
	drive.gapLength = gapLength
	drive.dataLength = dataLength
}

func (drive *diskDrive) seek(cylinder uint8) {
	drive.currentCylinder = cylinder
}

func (drive *diskDrive) sectorReader(sector uint8) (io.Reader, error) {
	buf := make([]byte, 512)
	for i := 0; i < 512; i++ {
		buf[i] = 0xFF
	}

	return bytes.NewBuffer(buf), nil
}
