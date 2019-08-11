package state

import "github.com/allankerr/freighter/spec"

type Status int

const (
	Creating Status = iota
	Created
	Running
	Stopped
)

func (s Status) String() string {
	switch s {
	case Creating:
		return "creating"
	case Created:
		return "created"
	case Running:
		return "running"
	case Stopped:
		return "stopped"
	default:
		return "unknown"
	}
}

type State struct {
	data stateData
}

func newDefaultState(ID string) *State {
	data := stateData{
		OCIVersion: spec.Version,
		ID:         ID,
		status:     Creating.String(),
	}
	return &State{data}
}

func newState(data stateData) *State {
	return &State{data}
}

type stateData struct {
	OCIVersion  string
	ID          string
	status      string
	pid         uint64
	bundle      string
	annotations map[string]string
}
