package ipc

import (
	"github.com/allankerr/freighter/cli/state"
)

type MessageType int

const (
	MessageUnknown MessageType = iota
	MessageInitSpec
	MessageStatusChange
)

type StatusChangePayload struct {
	Status state.Status
}
