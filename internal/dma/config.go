package dma

var MasterPortConfig = PortConfig{
	BaseAddressPort:   [4]uint16{0x00, 0x02, 0x04, 0x06},
	BaseCounterPort:   [4]uint16{0x01, 0x03, 0x05, 0x07},
	PagePort:          [4]uint16{0x87, 0x83, 0x81, 0x82},
	CommandPort:       0x08,
	RequestPort:       0x09,
	SingleMaskPort:    0x0A,
	ModePort:          0x0B,
	FlipFlopResetPort: 0x0C,
	MasterRestPort:    0x0D,
	StatusPort:        0x0F,
}

type PortConfig struct {
	BaseAddressPort   [4]uint16 //Tested RW
	BaseCounterPort   [4]uint16 //Tested RW
	PagePort          [4]uint16 //no test
	CommandPort       uint16
	RequestPort       uint16
	SingleMaskPort    uint16
	ModePort          uint16
	FlipFlopResetPort uint16 // Tested W. not support R
	MasterRestPort    uint16
	StatusPort        uint16
}
