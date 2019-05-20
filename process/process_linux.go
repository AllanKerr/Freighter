package process

import (
	"os"
	"syscall"

	"github.com/allankerr/freighter/spec"
)

type processLinux struct {
	wd   string
	path string
	argv []string
	envv []string
}

func New(config spec.Process) Process {
	return &processLinux{
		wd:   config.CWD,
		path: config.Args[0],
		argv: config.Args,
		envv: config.Env,
	}
}

func (proc *processLinux) Run() error {

	if err := os.Chdir(proc.wd); err != nil {
		return err
	}
	return syscall.Exec(proc.path, proc.argv, proc.envv)
}
