package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config describes the configuration of the service.
type Config struct {
	Application application `yaml:"application"`
}

type application struct {
	Authors []author `yaml:"authors"`
	Name    string   `yaml:"name"`
	Usage   string   `yaml:"usage"`
	Version string   `yaml:"version"`
}

type author struct {
	Name  string `yaml:"name"`
	Email string `yaml:"email"`
}

// New creates and returns a configuration object for the service.
func New(configFile string) (Config, error) {
	var config Config

	yamlBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		return Config{}, err
	}

	err = yaml.Unmarshal(yamlBytes, &config)
	if err != nil {
		return Config{}, err
	}

	return config, err
}
