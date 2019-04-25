package tty

import (
	"os"

	"golang.org/x/sys/unix"
)

type linuxTerminal struct {
	orig *unix.Termios
	fd   int
}

func NewTerminal(file *os.File) (Terminal, error) {
	fd := int(file.Fd())
	termios, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		return nil, err
	}
	terminal := &linuxTerminal{
		orig: termios,
	}
	return terminal, nil
}

func (c *linuxTerminal) SetRaw() error {

	termios, err := c.getTermios()
	if err != nil {
		return err
	}
	termios.Iflag &^= (unix.IGNBRK | unix.BRKINT | unix.PARMRK | unix.ISTRIP | unix.INLCR | unix.IGNCR | unix.ICRNL | unix.IXON)
	termios.Oflag &^= unix.OPOST
	termios.Lflag &^= (unix.ECHO | unix.ECHONL | unix.ICANON | unix.ISIG | unix.IEXTEN)
	termios.Cflag &^= (unix.CSIZE | unix.PARENB)
	termios.Cflag &^= unix.CS8

	termios.Oflag |= unix.OPOST

	return c.setTermios(termios)
}

func (c *linuxTerminal) Reset() error {
	return c.setTermios(c.orig)
}

func (c *linuxTerminal) getTermios() (*unix.Termios, error) {
	return unix.IoctlGetTermios(c.fd, unix.TCGETS)
}

func (c *linuxTerminal) setTermios(termios *unix.Termios) error {
	return unix.IoctlSetTermios(c.fd, unix.TCSETS, termios)
}
