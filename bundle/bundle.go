package bundle

import "github.com/allankerr/freighter/spec"

type Bundle interface {
	GetConfig() (*spec.Spec, error)
}
