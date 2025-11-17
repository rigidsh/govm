package dma

import (
	"io"
	"sync"
)

type BufRReadConnector struct {
	buf           []byte
	readPosition  int
	writePosition int
	capacity      int

	dreq Line

	mutex sync.Mutex
}

func NewBufReadConnector(dreq Line, bufSize int) *BufRReadConnector {
	return &BufRReadConnector{
		buf:           make([]byte, bufSize),
		readPosition:  0,
		writePosition: 0,
		capacity:      bufSize,
		dreq:          dreq,
	}
}

func (connector *BufRReadConnector) ReadFrom(reader io.Reader) error {
	connector.mutex.Lock()
	defer connector.mutex.Unlock()
	defer func() {
		connector.dreq.Set(connector.capacity != len(connector.buf))
	}()

	if connector.writePosition == connector.readPosition && connector.capacity == 0 {
		return nil
	}

	if connector.writePosition < connector.readPosition {
		size, err := reader.Read(connector.buf[connector.writePosition:connector.readPosition])
		connector.writePosition = connector.writePosition + size
		connector.capacity = connector.capacity - size
		if err != nil {
			return err
		}
	} else {
		size, err := reader.Read(connector.buf[connector.writePosition:])
		connector.writePosition = connector.writePosition + size
		connector.capacity = connector.capacity - size
		if connector.writePosition == len(connector.buf) {
			connector.writePosition = 0
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (connector *BufRReadConnector) Read(buf []byte) uint16 {
	connector.mutex.Lock()
	defer connector.mutex.Unlock()
	defer func() {
		connector.dreq.Set(connector.capacity != len(connector.buf))
	}()

	if connector.writePosition == connector.readPosition && connector.capacity == len(connector.buf) {
		return 0
	}

	if connector.writePosition > connector.readPosition {
		sizeToCopy := connector.writePosition - connector.readPosition
		if sizeToCopy > len(buf) {
			sizeToCopy = len(buf)
		}
		copy(buf, connector.buf[connector.readPosition:connector.readPosition+sizeToCopy])
		connector.readPosition = connector.readPosition + sizeToCopy
		return uint16(sizeToCopy)
	} else {
		sizeToCopy := len(buf) - connector.readPosition
		if sizeToCopy > len(buf) {
			sizeToCopy = len(buf)
		}
		copy(buf, connector.buf[connector.readPosition:connector.readPosition+sizeToCopy])
		connector.readPosition = connector.readPosition + sizeToCopy
		if connector.readPosition == len(connector.buf) {
			connector.readPosition = 0
		}
		return uint16(sizeToCopy)
	}

	return 0
}

func (connector *BufRReadConnector) Write(buf []byte) uint16 {
	return 0
}
