package dma

type line struct {
	value  bool
	change chan interface{}
}

func (line *line) set(value bool) {
	notify := line.value != value && value
	line.value = value
	if notify {
		line.change <- nil
	}
}

func newLine() *line {
	return &line{
		value:  false,
		change: make(chan interface{}, 1),
	}
}
