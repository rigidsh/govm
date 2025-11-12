package dma

type Line struct {
	value   bool
	posEdge chan interface{}
}

func (line *Line) Set(value bool) {
	notify := line.value != value && value
	line.value = value
	if notify {
		select {
		case line.posEdge <- nil:
		}

	}
}

func (line *Line) Get() bool {
	return line.value
}

func (line *Line) PosEdge() chan interface{} {
	return line.posEdge
}

func NewLine() *Line {
	return &Line{
		value:   false,
		posEdge: make(chan interface{}, 0),
	}
}
