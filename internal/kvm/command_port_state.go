package kvm

type commandPortState interface {
}
type waitRequestCommandPortStatus struct{}

type readRequestParamCommandPortStatus struct {
	buf      []byte
	size     uint8
	position uint8
	command  *CommandDefinition
}

type processingCommandPortStatus struct{}

type responseCommandPortStatus struct {
	buf          []byte
	readPosition uint8
}
