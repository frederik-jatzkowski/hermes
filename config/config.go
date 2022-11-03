package config

import (
	"encoding/json"
	"fmt"
)

type Config struct {
	Unix     int64     `json:"unix"`
	Gateways []Gateway `json:"gateways"`
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

func (config *Config) validate() error {
	var err error

	// validate Gateways
	for i, gateway := range config.Gateways {
		err = gateway.validate()
		if err != nil {
			return fmt.Errorf("invalid gateway config: %s", err)
		}
		// reassign to save changes
		config.Gateways[i] = gateway
	}

	return err
}
