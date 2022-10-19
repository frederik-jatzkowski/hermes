package core

import (
	"github.com/frederik-jatzkowski/hermes/config"
)

type Core struct {
	gateways []Gateway
}

func NewCore(config config.Config) *Core {
	var (
		core Core
	)

	for _, gatewayConfig := range config.Gateways {
		core.gateways = append(core.gateways, *NewGateway(gatewayConfig))
	}

	return &core
}

func (core Core) Start() error {
	var err error

	for _, gateway := range core.gateways {
		err = gateway.Start()
		if err != nil {
			return err
		}
	}

	return err
}

func (core Core) Stop() {
	for _, gateway := range core.gateways {
		gateway.Stop()
	}
}
