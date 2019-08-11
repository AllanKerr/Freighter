package state

type StateStore interface {
	Create(containerID string) (*State, error)
	Load(containerID string) (*State, error)
	Save(state *State) error
}
