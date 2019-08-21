package container

import "golang.org/x/sys/unix"

func setHostname(hostname string) error {
	return unix.Sethostname([]byte(hostname))
}
