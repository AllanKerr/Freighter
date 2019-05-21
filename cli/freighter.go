package cli

import (
	"os"
)

type freighterLinux struct {
	rootPath string
}

const defaultRootPath = "/run/freighter"

func NewDefaultConfiguration() (Freighter, error) {
	return New(defaultRootPath)
}

func New(rootPath string) (Freighter, error) {

	if err := os.MkdirAll(rootPath, os.ModePerm); err != nil {
		return nil, err
	}
	return &freighterLinux{
		rootPath,
	}, nil
}

func (f *freighterLinux) Create(containerId string, bundlePath string) error {

	return nil
}

func (f *freighterLinux) Start(containerId string) error {

	return nil
}
