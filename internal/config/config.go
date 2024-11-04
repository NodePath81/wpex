package config

import (
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Bind          string   `yaml:"bind"`
	Port          uint     `yaml:"port"`
	BroadcastRate uint     `yaml:"broadcast_rate"`
	Allows        []string `yaml:"allows"`
}

// Read config from file
func ReadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return readConfigFromReader(file)
}

func readConfigFromReader(reader io.Reader) (*Config, error) {
	var config Config
	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
