package fs

type RootFS interface {
	PrepareRoot() error
	AddSystemCommands() error
	PivotRoot() error
}
