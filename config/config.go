package config

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	Gateways []Gateway `json:"gateways"`
	Redirect bool      `json:"redirect"`
	LogLevel string    `json:"logLevel"`
}

func ReadConfigFile() ([]byte, error) {
	var (
		data []byte
		err  error
	)

	// read config file
	data, err = osReadFile("/var/hermes/config.xml")
	if err != nil {
		return data, fmt.Errorf("could not read config file: %s", err)
	}

	return data, err
}

func NewConfig(data []byte) (Config, error) {
	var (
		config Config
		err    error
	)

	// parse config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("could not parse config: %s", err)
	}

	// validate config
	err = config.validate()
	if err != nil {
		return config, fmt.Errorf("error in config: %s", err)
	}

	return config, err
}

func (config Config) validate() error {
	var err error

	// validate Gateways
	for _, gateway := range config.Gateways {
		err = gateway.validate()
		if err != nil {
			return fmt.Errorf("invalid gateway config: %s", err)
		}
	}

	return err
}
