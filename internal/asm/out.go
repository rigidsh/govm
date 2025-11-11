package asm

func OUT(port, from Argument) Command {
	result := make([]byte, 0)
	switch from.(type) {
	case al:
		switch port.(type) {
		case impl8:
			result = append(result, 0xE6)
		default:
			panic("Not supported operation")
		}
	}

	result = append(result, port.bytes()...)
	result = append(result, from.bytes()...)

	return result
}
