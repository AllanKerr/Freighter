package cli

import (
	"os"

	"github.com/allankerr/freighter/bundle"
	"github.com/allankerr/freighter/cli/state"
	"github.com/allankerr/freighter/log"
)

type freighterLinux struct {
	rootPath   string
	stateStore state.StateStore
}

const defaultRootPath = "/run/freighter"

func NewDefaultConfiguration() (Freighter, error) {
	return New(defaultRootPath)
}

func New(rootPath string) (Freighter, error) {
	log.InitParentLogger()

	if err := os.MkdirAll(rootPath, os.ModePerm); err != nil {
		return nil, err
	}
	stateStore := state.NewFileStateStore(rootPath)
	return &freighterLinux{
		rootPath,
		stateStore,
	}, nil
}

func (f *freighterLinux) Create(containerID string, bundlePath string) error {

	bundle, err := bundle.New(bundlePath)
	if err != nil {
		return err
	}
	config, err := bundle.GetConfig()
	if err != nil {
		return err
	}
	state, err := f.stateStore.Create(containerID)
	if err != nil {
		return err
	}
	if err := createContainer(state, config); err != nil {
		return err
	}
	return nil
}

func (f *freighterLinux) Start(containerID string) error {

	return nil
}
