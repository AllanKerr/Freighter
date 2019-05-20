package uts

import "golang.org/x/sys/unix"

func SetHostname(hostname string) error {
	return unix.Sethostname([]byte(hostname))
}
