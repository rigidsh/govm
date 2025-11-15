package fdc

import "io"

type diskDrive struct {
	disk            Disk
	currentCylinder uint8

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

func (drive *diskDrive) sectorReader(fromSector, toSector, fromHead, toHead uint8) (io.Reader, error) {
	return drive.disk.read(), nil
}
