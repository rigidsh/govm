package asm

func CMP(a1, a2 Argument) Command {

	result := make([]byte, 0)
	switch a1.(type) {
	case ax:
		switch a2.(type) {
		case impl16:
			result = append(result, 0x3D)
		default:
			panic("Not supported operation")
		}
	case si:
		switch a2.(type) {
		case impl16:
			result = append(result, 0x81, 0xFE)
		default:
			panic("Not supported operation")
		}
	default:
		panic("Not supported operation")
	}

	result = append(result, a1.bytes()...)
	result = append(result, a2.bytes()...)

	return result
}
