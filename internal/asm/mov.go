package asm

func MOV(to, from Argument) Command {
	result := make([]byte, 0)
	switch to.(type) {
	case ah:
		switch from.(type) {
		case impl8:
			result = append(result, 0xB4)
		default:
			panic("Not supported operation")
		}
	case al:
		switch from.(type) {
		case impl8:
			result = append(result, 0xB0)
		case bxRelSIAddr:
			result = append(result, 0x8A, 0x01)
		default:
			panic("Not supported operation")
		}
	case ax:
		switch from.(type) {
		case impl16:
			result = append(result, 0xA1)
		default:
			panic("Not supported operation")
		}
	case bl:
		switch from.(type) {
		case impl8:
			result = append(result, 0xB3)
		default:
			panic("Not supported operation")
		}
	case bx:
		switch from.(type) {
		case impl16:
			result = append(result, 0xBB)
		default:
			panic("Not supported operation")
		}
	case cl:
		switch from.(type) {
		case impl8:
			result = append(result, 0xB1)
		default:
			panic("Not supported operation")
		}
	case ch:
		switch from.(type) {
		case impl8:
			result = append(result, 0xB5)
		default:
			panic("Not supported operation")
		}
	case dh:
		switch from.(type) {
		case impl8:
			result = append(result, 0xB6)
		default:
			panic("Not supported operation")
		}
	case es:
		switch from.(type) {
		case ax:
			result = append(result, 0x8E, 0xC0)
		default:
			panic("Not supported operation")
		}
	case si:
		switch from.(type) {
		case impl16:
			result = append(result, 0xBE)
		default:
			panic("Not supported operation")
		}
	default:
		panic("Not supported operation")
	}

	result = append(result, to.bytes()...)
	result = append(result, from.bytes()...)

	return result
}

//mov BX, 0x7C00
