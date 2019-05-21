package cli

type Freighter interface {
	Create(containerId string, bundlePath string) error
	Start(containerId string) error
}
