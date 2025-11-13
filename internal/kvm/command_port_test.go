package kvm

import (
	"testing"
)

func TestCommandPort(t *testing.T) {

	port := NewCommandPort()
	port.RegisterCommand(0x42, &CommandDefinition{
		ArgumentSize: 2,
		Command: func(argument []byte, resultCallback func([]byte)) {
			resultCallback(argument)
		},
	})

	tests := []struct {
		name           string
		inputData      []byte
		expectedResult []byte
	}{
		{
			name:           "Echo 2 byte argument",
			inputData:      []byte{0x42, 0x01, 0x02},
			expectedResult: []byte{0x01, 0x02},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			for _, data := range tt.inputData {
				port.OnWrite([]byte{data})
			}
			for _, expectedByte := range tt.expectedResult {
				actualByte := port.OnRead(1)[0]
				if actualByte != expectedByte {
					t.Errorf("Incorrect response. Expected: %X, Actual: %X", expectedByte, actualByte)
				}
			}
		})
	}
}
