package bundle

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/allankerr/freighter/spec"
)

type bundleLinux struct {
	bundlePath string
}

const configFileName = "config.json"

func New(bundlePath string) (Bundle, error) {

	configPath := getConfigPath(bundlePath)
	if hasConfig, err := hasConfig(configPath); err != nil {
		return nil, err
	} else if !hasConfig {
		return nil, fmt.Errorf("Missing config file: %s", configPath)
	}

	return &bundleLinux{
		bundlePath: bundlePath,
	}, nil
}

func (bundle *bundleLinux) GetConfig() (*spec.Spec, error) {

	config := &spec.Spec{}
	configPath := bundle.getConfigPath()
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(configFile, config); err != nil {
		return nil, err
	}
	return config, nil
}

func (bundle *bundleLinux) getConfigPath() string {
	return getConfigPath(bundle.bundlePath)
}

func hasConfig(configPath string) (bool, error) {
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			err = nil
		}
		return false, err
	}
	return true, nil
}

func getConfigPath(bundlePath string) string {
	return path.Join(bundlePath, configFileName)
}
