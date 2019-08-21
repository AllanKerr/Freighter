package container

import "github.com/allankerr/freighter/spec"

type RootFS interface {
	PrepareRoot(rootPropagation string) error
	CreateMounts(mounts []spec.Mount) error
	CreateDevices(devices []spec.Device) error
	PivotRoot() error
	FinalizeRoot() error
}
