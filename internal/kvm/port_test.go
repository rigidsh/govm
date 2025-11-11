package kvm

import (
	"testing"
)

func TestFlipFlopPort_OnWrite(t *testing.T) {

	register := uint16(0x0000)
	flipFlop := NewFlipFlop()
	port := FlipFlopPort(CreateRegister16Port(&register), flipFlop)

	tests := []struct {
		name                string
		inputData           [][]byte
		expectedWrittenData uint16
	}{
		{
			name:                "Simple write by 8",
			inputData:           [][]byte{{0x34}, {0x12}},
			expectedWrittenData: 0x1234,
		},
		{
			name:                "Simple write by 16",
			inputData:           [][]byte{{0x34, 0x12}},
			expectedWrittenData: 0x1234,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flipFlop.Reset()

			for _, data := range tt.inputData {
				port.OnWrite(data)
			}

			if register != tt.expectedWrittenData {
				t.Errorf("OnWrite(%v) записал данные %v, ожидалось %v",
					tt.inputData, register, tt.expectedWrittenData)
			}
		})
	}
}
