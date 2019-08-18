package container

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func findFileDescriptor(name string) (uintptr, error) {

	val, err := findEnvironmentVariable(name)
	if err != nil {
		return 0, err
	}
	fd, err := strconv.ParseUint(val, 10, 64)
	return uintptr(fd), err
}

func findEnvironmentVariable(name string) (string, error) {

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 && parts[0] == name {
			return parts[1], nil
		}
	}
	return "", fmt.Errorf("Unable to find environment variable: %s", name)
}
