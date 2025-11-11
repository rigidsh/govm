package dma

func mmtBuffConnector(dma *DMA) *ChannelConnector {
	return &ChannelConnector{
		Write: func(buf []byte) uint16 {
			dma.vm.Memory().Write(dma.channels[1].phyAddress(), buf)
			dma.channels[1].currentAddress = dma.channels[1].currentAddress + uint16(len(buf))
			return uint16(len(buf))
		},
	}
}
