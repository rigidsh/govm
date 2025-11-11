package asm

import "encoding/binary"

type Argument interface {
	bytes() []byte
}

func SI() Argument {
	return si{}
}

type si struct{}

func (_ si) bytes() []byte { return []byte{} }

// register A

func AL() Argument {
	return al{}
}

type al struct{}

func (_ al) bytes() []byte { return []byte{} }

func AH() Argument {
	return ah{}
}

type ah struct{}

func (_ ah) bytes() []byte { return []byte{} }

func AX() Argument {
	return ax{}
}

type ax struct{}

func (_ ax) bytes() []byte { return []byte{} }

// register B

func BL() Argument {
	return bl{}
}

type bl struct{}

func (_ bl) bytes() []byte { return []byte{} }

func BX() Argument {
	return bx{}
}

type bx struct{}

func (_ bx) bytes() []byte { return []byte{} }

// register C
func CL() Argument {
	return cl{}
}

type cl struct{}

func (_ cl) bytes() []byte { return []byte{} }

func CH() Argument {
	return ch{}
}

type ch struct{}

func (_ ch) bytes() []byte { return []byte{} }

// register D
func DH() Argument {
	return dh{}
}

type dh struct{}

func (_ dh) bytes() []byte { return []byte{} }

func ES() Argument {
	return es{}
}

type es struct{}

func (_ es) bytes() []byte { return []byte{} }

func RelAddr(value int8) Argument {
	return relAddr(value)
}

type relAddr int8

func (v relAddr) bytes() []byte { return []byte{byte(v)} }

func BXRelSIAddr() Argument {
	return bxRelSIAddr{}
}

type bxRelSIAddr struct {
}

func (_ bxRelSIAddr) bytes() []byte { return []byte{} }

func SegAddr(segment, offset uint16) Argument {
	return segAddr{
		segment: segment,
		offset:  offset,
	}
}

type segAddr struct {
	segment uint16
	offset  uint16
}

func (v segAddr) bytes() []byte {
	result := []byte{}
	result = binary.LittleEndian.AppendUint16(result, v.offset)
	return binary.LittleEndian.AppendUint16(result, v.segment)
}

func Impl8(value uint8) Argument {
	return impl8(value)
}

type impl8 uint8

func (v impl8) bytes() []byte { return []byte{byte(v)} }

func Impl16(value uint16) Argument {
	return impl16(value)
}

type impl16 uint16

func (v impl16) bytes() []byte { return binary.LittleEndian.AppendUint16([]byte{}, uint16(v)) }
