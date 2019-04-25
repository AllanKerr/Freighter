package tty

type Terminal interface {
	SetRaw() error
	Reset() error
}

type TerminalMaster interface {
	Terminal
	GetSlavePath() string
	UnlockSlave() error
	Listen()
}
