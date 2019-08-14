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
	data          stateData
	directoryPath string
}

func newDefaultState(ID string, directoryPath string) *State {
	data := stateData{
		OCIVersion: spec.Version,
		ID:         ID,
		status:     Creating.String(),
	}
	return &State{data, directoryPath}
}

func newState(data stateData, directoryPath string) *State {
	return &State{data, directoryPath}
}

func (s *State) ID() string {
	return s.data.ID
}

func (s *State) DirectoryPath() string {
	return s.directoryPath
}

type stateData struct {
	OCIVersion  string
	ID          string
	status      string
	pid         uint64
	bundle      string
	annotations map[string]string
}
