package kvm

type CommandDefinition struct {
	ArgumentSize uint8
	Command      func(argument []byte)
}

func NewCommandDefinition(argumentSize uint8, command func(argument []byte)) *CommandDefinition {
	return &CommandDefinition{ArgumentSize: argumentSize, Command: command}
}

type CommandPort struct {
	state    commandPortState
	commands map[byte]*CommandDefinition
}

func NewCommandPort() *CommandPort {
	return &CommandPort{
		state:    &waitRequestCommandPortStatus{},
		commands: make(map[byte]*CommandDefinition),
	}
}

func (port *CommandPort) RegisterCommand(command byte, definition *CommandDefinition) {
	port.commands[command] = definition
}

func (port *CommandPort) WriteResult(result []byte) {
	switch port.state.(type) {
	case *processingCommandPortStatus:
		port.state = &responseCommandPortStatus{
			buf:          result,
			readPosition: 0,
		}
	}
}

func (port *CommandPort) processCommandData(data byte) {
	switch port.state.(type) {
	case *waitRequestCommandPortStatus:
		if command, ok := port.commands[data]; ok {
			if command.ArgumentSize == 0 {
				command.Command([]byte{})
				port.state = &processingCommandPortStatus{}
			} else {
				port.state = &readRequestParamCommandPortStatus{
					buf:      make([]byte, command.ArgumentSize),
					size:     command.ArgumentSize,
					command:  command,
					position: 0,
				}
			}
		}
	case *readRequestParamCommandPortStatus:
		paramReader := port.state.(*readRequestParamCommandPortStatus)
		paramReader.buf[paramReader.position] = data
		paramReader.position = paramReader.position + 1
		if paramReader.position == paramReader.size {
			port.state = &processingCommandPortStatus{}
			paramReader.command.Command(paramReader.buf)
		}
	}
}

func (port *CommandPort) processRead() byte {
	switch port.state.(type) {
	case *responseCommandPortStatus:
		response := port.state.(*responseCommandPortStatus)
		if response.readPosition < uint8(len(response.buf)) {
			result := response.buf[response.readPosition]
			response.readPosition = response.readPosition + 1
			if response.readPosition == uint8(len(response.buf)) {
				port.state = &waitRequestCommandPortStatus{}
			}
			return result
		}
	}

	return 0
}

func (port *CommandPort) OnWrite(data []byte) {
	for _, b := range data {
		port.processCommandData(b)
	}
}

func (port *CommandPort) OnRead(size uint8) []byte {
	result := make([]byte, size)
	for i := uint8(0); i < size; i++ {
		result[i] = port.processRead()
	}
	return result
}
