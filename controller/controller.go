package controller

import (
	"os"
	"os/exec"

	"github.com/allankerr/freighter/log"
)

type Controller struct {
}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Run() error {

	remoteLogger := log.InitParentLogger()

	log.Info("parent log test")

	cmd := exec.Command("/proc/self/exe", "init")
	cmd.ExtraFiles = []*os.File{remoteLogger.Child()}
	cmd.Env = []string{}
	cmd.Start()

	return cmd.Wait()
}
