package log

import "os"

// RemoteLogger listens for input from the child file and logs it.
type RemoteLogger interface {
	Child() *os.File
	listen()
}
