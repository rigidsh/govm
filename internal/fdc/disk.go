package fdc

import (
	"bytes"
	"io"
)

type Disk interface {
	read() io.Reader
}

type InMemoryRawDisk []byte

func (disk InMemoryRawDisk) read() io.Reader {
	return bytes.NewBuffer(disk)
}
