package dma

import (
	"io"
	"sync"
)

type BuffReadConnector struct {
	buff          []byte
	readPosition  int
	writePosition int
	capacity      int

	dreq Line

	mutex sync.Mutex
}

func NewBuffReadConnector(dreq Line, buffSize int) *BuffReadConnector {
	return &BuffReadConnector{
		buff:          make([]byte, buffSize),
		readPosition:  0,
		writePosition: 0,
		capacity:      buffSize,
		dreq:          dreq,
	}
}

func (connector *BuffReadConnector) ReadFrom(reader io.Reader) error {
	connector.mutex.Lock()
	defer connector.mutex.Unlock()
	defer func() {
		connector.dreq.Set(connector.capacity != len(connector.buff))
	}()

	if connector.writePosition == connector.readPosition && connector.capacity == 0 {
		return nil
	}

	if connector.writePosition < connector.readPosition {
		size, err := reader.Read(connector.buff[connector.writePosition:connector.readPosition])
		connector.writePosition = connector.writePosition + size
		connector.capacity = connector.capacity - size
		if err != nil {
			return err
		}
	} else {
		size, err := reader.Read(connector.buff[connector.writePosition:])
		connector.writePosition = connector.writePosition + size
		connector.capacity = connector.capacity - size
		if connector.writePosition == len(connector.buff) {
			connector.writePosition = 0
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (connector *BuffReadConnector) Read(buf []byte) uint16 {
	connector.mutex.Lock()
	defer connector.mutex.Unlock()
	defer func() {
		connector.dreq.Set(connector.capacity != len(connector.buff))
	}()

	if connector.writePosition == connector.readPosition && connector.capacity == len(connector.buff) {
		return 0
	}

	if connector.writePosition > connector.readPosition {
		sizeToCopy := connector.writePosition - connector.readPosition
		if sizeToCopy > len(buf) {
			sizeToCopy = len(buf)
		}
		copy(buf, connector.buff[connector.readPosition:connector.readPosition+sizeToCopy])
		connector.readPosition = connector.readPosition + sizeToCopy
		return uint16(sizeToCopy)
	} else {
		sizeToCopy := len(buf) - connector.readPosition
		if sizeToCopy > len(buf) {
			sizeToCopy = len(buf)
		}
		copy(buf, connector.buff[connector.readPosition:connector.readPosition+sizeToCopy])
		connector.readPosition = connector.readPosition + sizeToCopy
		if connector.readPosition == len(connector.buff) {
			connector.readPosition = 0
		}
		return uint16(sizeToCopy)
	}

	return 0
}

func (connector *BuffReadConnector) Write(buf []byte) uint16 {
	return 0
}
