package dma

type mmtBuffConnector struct {
	dma *DMA
}

func (_ *mmtBuffConnector) Read(buf []byte) uint16 {
	return 0
}

func (connector *mmtBuffConnector) Write(buf []byte) uint16 {
	connector.dma.vm.Memory().Write(connector.dma.channels[1].phyAddress(), buf)
	connector.dma.channels[1].currentAddress = connector.dma.channels[1].currentAddress + uint16(len(buf))
	return uint16(len(buf))
}
