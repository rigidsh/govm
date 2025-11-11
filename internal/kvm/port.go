package kvm

import (
	"io"
	"unsafe"
)

type IOPort interface {
	OnWrite(data []byte)
	OnRead(size uint8) []byte
}

func CreateRegister16Port(register *uint16) IOPort {
	return &registerPort{
		register:     unsafe.Pointer(register),
		registerSize: 2,
	}
}

func CreateRegister8Port(register *uint8) IOPort {
	return &registerPort{
		register:     unsafe.Pointer(register),
		registerSize: 1,
	}
}

type registerPort struct {
	register     unsafe.Pointer
	registerSize uint8
}

func ROPort(port IOPort) IOPort {
	return &accessControlPort{
		IOPort:   port,
		writable: false,
		readable: true,
	}
}
func WOPort(port IOPort) IOPort {

	return &accessControlPort{
		IOPort:   port,
		writable: true,
		readable: false,
	}
}

type accessControlPort struct {
	IOPort

	writable bool
	readable bool
}

func (port accessControlPort) OnWrite(data []byte) {
	if !port.writable {
		return
	}

	port.IOPort.OnWrite(data)
}

func (port accessControlPort) OnRead(size uint8) []byte {
	if !port.readable {
		return make([]byte, size)
	}

	return port.IOPort.OnRead(size)
}

type oneByteOperationsPort struct {
	IOPort
}

func OneByteOperationsPort(IOPort IOPort) IOPort {
	return &oneByteOperationsPort{IOPort: IOPort}
}

func (port *oneByteOperationsPort) OnWrite(data []byte) {
	port.IOPort.OnWrite([]byte{data[len(data)-1]})
}

func (port *oneByteOperationsPort) OnRead(size uint8) []byte {
	return port.IOPort.OnRead(1)
}

func (port *registerPort) OnWrite(data []byte) {
	copyBytes := uint8(len(data))
	if copyBytes > port.registerSize {
		copyBytes = port.registerSize
	}

	copy(port.registerBytes(), data[:copyBytes])
}

func (port *registerPort) registerBytes() []byte {
	return unsafe.Slice((*byte)(port.register), port.registerSize)
}

func (port *registerPort) OnRead(size uint8) []byte {
	copyBytes := size
	if copyBytes > port.registerSize {
		copyBytes = port.registerSize
	}

	result := make([]byte, size)

	copy(result, port.registerBytes()[:copyBytes])
	return result
}

func CreateIOPort(in io.Reader, out io.Writer) IOPort {
	return &ioPort{
		in:  in,
		out: out,
	}
}

type ioPort struct {
	in  io.Reader
	out io.Writer
}

func (port *ioPort) OnWrite(data []byte) {
	port.out.Write(data)
}

func (port *ioPort) OnRead(size uint8) []byte {
	buf := make([]byte, size)
	port.in.Read(buf)

	return buf
}

func CompositePort(write, read IOPort) IOPort {
	return &compositePort{
		write: write,
		read:  read,
	}
}

type compositePort struct {
	write IOPort
	read  IOPort
}

func (port *compositePort) OnWrite(data []byte) {
	port.write.OnWrite(data)
}

func (port *compositePort) OnRead(size uint8) []byte {
	return port.read.OnRead(size)
}

func NopPort() IOPort {
	return &nopPort{}
}

type nopPort struct{}

func (_ nopPort) OnWrite(data []byte) {
}

func (_ nopPort) OnRead(size uint8) []byte {
	return make([]byte, size)
}

func CallbackPort(callback func(write bool, data []byte) []byte) IOPort {
	return &callbackPort{
		callback: callback,
	}
}

type callbackPort struct {
	callback func(write bool, data []byte) []byte
}

func (port *callbackPort) OnWrite(data []byte) {
	port.callback(true, data)
}

func (port *callbackPort) OnRead(size uint8) []byte {
	return port.callback(false, nil)
}

type FlipFlop struct {
	isSet bool
	buf   byte
}

func NewFlipFlop() *FlipFlop {
	return &FlipFlop{}
}

func (f *FlipFlop) Set(value byte) {
	f.isSet = true
	f.buf = value
}

func (f *FlipFlop) Get() (byte, bool) {
	return f.buf, f.isSet
}

func (f *FlipFlop) Reset() {
	f.isSet = false
}

func FlipFlopPort(port IOPort, flipFlop *FlipFlop) IOPort {
	return &flipFlopPort{
		IOPort:   port,
		flipFlop: flipFlop,
	}
}

type flipFlopPort struct {
	IOPort
	flipFlop *FlipFlop
}

func (port *flipFlopPort) OnWrite(data []byte) {
	if len(data) != 1 {
		return
	}
	if flipFlopValue, flipFlopSet := port.flipFlop.Get(); flipFlopSet {
		port.IOPort.OnWrite([]byte{flipFlopValue, data[0]})
		port.flipFlop.Reset()
	} else {
		port.flipFlop.Set(data[0])
	}
}

func (port *flipFlopPort) OnRead(size uint8) []byte {
	if flipFlopValue, flipFlopSet := port.flipFlop.Get(); flipFlopSet {
		port.flipFlop.Reset()
		return []byte{flipFlopValue}
	} else {
		data := port.IOPort.OnRead(2)
		port.flipFlop.Set(data[1])
		return []byte{data[0]}
	}
}
