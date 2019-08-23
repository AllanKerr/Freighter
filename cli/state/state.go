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
	store         StateStore
	data          stateData
	directoryPath string
}

func newDefaultState(store StateStore, ID string, directoryPath string) *State {
	data := stateData{
		OCIVersion: spec.Version,
		ID:         ID,
		Status:     Creating.String(),
	}
	return newState(store, data, directoryPath)
}

func newState(store StateStore, data stateData, directoryPath string) *State {

	return &State{store, data, directoryPath}
}

func (s *State) ID() string {
	return s.data.ID
}

func (s *State) DirectoryPath() string {
	return s.directoryPath
}

func (s *State) SetPID(PID uint64) {
	s.data.PID = PID
}

func (s *State) PID() uint64 {
	return s.data.PID
}

func (s *State) SetStatus(status Status) {
	s.data.Status = status.String()
}

func (s *State) Save() error {
	return s.store.Save(s)
}

type stateData struct {
	OCIVersion  string            `json:"ociVersion"`
	ID          string            `json:"id"`
	Status      string            `json:"status"`
	PID         uint64            `json:"pid"`
	Bundle      string            `json:"bundle"`
	Annotations map[string]string `json:"annotations,omitempty"`
}
