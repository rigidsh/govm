package dma

import "fmt"

type Line interface {
	Get() bool
	Set(value bool)
}

type ObservableLine struct {
	value   bool
	posEdge chan interface{}
}

func (line *ObservableLine) Set(value bool) {
	notify := line.value != value && value
	line.value = value
	if notify {
		select {
		case line.posEdge <- nil:
		default:
			fmt.Println("No active listener")
		}

	}
}

func (line *ObservableLine) Get() bool {
	return line.value
}

func (line *ObservableLine) PosEdge() chan interface{} {
	return line.posEdge
}

func NewLine() *ObservableLine {
	return &ObservableLine{
		value:   false,
		posEdge: make(chan interface{}, 0),
	}
}
