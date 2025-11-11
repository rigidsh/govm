package asm

func INT(intNumber Argument) Command {
	result := make([]byte, 0)
	switch intNumber.(type) {
	case impl8:
		result = append(result, 0xCD)
	default:
		panic("Not supported operation")
	}

	result = append(result, intNumber.bytes()...)

	return result
}
