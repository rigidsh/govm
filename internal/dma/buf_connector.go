package dma

import (
	"io"
	"sync"
)

type BufConnector struct {
	buf           []byte
	readPosition  int
	writePosition int

	mutex sync.Mutex
}

func NewBufConnector(bufSize int) *BufConnector {
	return &BufConnector{
		buf:           make([]byte, bufSize),
		readPosition:  bufSize,
		writePosition: 0,
	}
}

func (connector *BufConnector) ReadFrom(reader io.Reader) error {
	connector.mutex.Lock()
	defer connector.mutex.Unlock()

	if connector.writePosition == connector.readPosition {
		return nil
	}

	if connector.writePosition < connector.readPosition {
		size, err := reader.Read(connector.buf[connector.writePosition:connector.readPosition])
		connector.writePosition = connector.writePosition + size
		if err != nil {
			return err
		}
	} else {
		size, err := reader.Read(connector.buf[connector.writePosition:])
		connector.writePosition = connector.writePosition + size
		if connector.writePosition == len(connector.buf) {
			connector.writePosition = 0
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (connector *BufConnector) Read(buf []byte) uint16 {
	connector.mutex.Lock()
	defer connector.mutex.Unlock()

	//if connector.writePosition == connector.readPosition {
	//	return 0
	//}
	//
	//if connector.writePosition > connector.readPosition {
	//	sizeToCopy := connector.writePosition - connector.readPosition
	//	if sizeToCopy > len(buf) {
	//		sizeToCopy = len(buf)
	//	}
	//	copy(buf, connector.buf[connector.readPosition:connector.readPosition+sizeToCopy])
	//	connector.readPosition = connector.readPosition + sizeToCopy
	//	return uint16(sizeToCopy)
	//} else {
	//	sizeToCopy := len(buf) - connector.readPosition
	//	if sizeToCopy > len(buf) {
	//		sizeToCopy = len(buf)
	//	}
	//	copy(buf, connector.buf[connector.readPosition:connector.readPosition+sizeToCopy])
	//	connector.readPosition = connector.readPosition + sizeToCopy
	//	if connector.readPosition == len(connector.buf) {
	//		connector.readPosition = 0
	//	}
	//}

	return 0
}

func (connector *BufConnector) Write(buf []byte) uint16 {
	return 0
}
