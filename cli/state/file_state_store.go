package state

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
)

const stateFileName = "state.json"

type fileStateStore struct {
	rootPath string
	mux      sync.Mutex
}

func NewFileStateStore(rootPath string) StateStore {
	return &fileStateStore{
		rootPath: rootPath,
	}
}

func (f *fileStateStore) Create(containerID string) (*State, error) {
	f.mux.Lock()
	defer f.mux.Unlock()
	return f.createNewState(containerID)
}

func (f *fileStateStore) createNewState(containerID string) (*State, error) {
	dirPath := f.buildStateDirectoryPath(containerID)
	if _, err := os.Stat(dirPath); err == nil || !os.IsNotExist(err) {
		if err == nil {
			err = fmt.Errorf("%s already exists", containerID)
		}
		return nil, err
	}
	if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		return nil, err
	}
	state := newDefaultState(containerID, dirPath)
	if err := f.save(state); err != nil {
		return nil, err
	}
	return state, nil
}

func (f *fileStateStore) Load(containerID string) (*State, error) {
	f.mux.Lock()
	defer f.mux.Unlock()
	return f.load(containerID)
}

func (f *fileStateStore) load(containerID string) (*State, error) {
	statePath := f.buildStatePath(containerID)
	if _, err := os.Stat(statePath); err != nil {
		return nil, err
	}
	stateFile, err := ioutil.ReadFile(statePath)
	if err != nil {
		return nil, err
	}
	data := stateData{}
	if err := json.Unmarshal(stateFile, &data); err != nil {
		return nil, err
	}
	dirPath := f.buildStateDirectoryPath(containerID)
	state := newState(data, dirPath)
	return state, nil
}

func (f *fileStateStore) Save(state *State) error {
	f.mux.Lock()
	defer f.mux.Unlock()
	return nil
}

func (f *fileStateStore) save(state *State) error {
	json, err := json.Marshal(state.data)
	if err != nil {
		return err
	}
	statePath := f.buildStatePath(state.data.ID)
	return ioutil.WriteFile(statePath, json, os.ModePerm)
}

func (f *fileStateStore) buildStateDirectoryPath(containerID string) string {
	return path.Join(f.rootPath, containerID)
}

func (f *fileStateStore) buildStatePath(containerID string) string {
	dirPath := f.buildStateDirectoryPath(containerID)
	return path.Join(dirPath, stateFileName)
}
