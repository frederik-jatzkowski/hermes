package core

import (
	"github.com/frederik-jatzkowski/hermes/config"
	"github.com/frederik-jatzkowski/hermes/logs"
)

type Core struct {
	gateways []*Gateway
}

func NewCore(config config.Config) *Core {
	var (
		core Core
	)

	for _, gatewayConfig := range config.Gateways {
		core.gateways = append(core.gateways, NewGateway(gatewayConfig))
	}

	return &core
}

func (core *Core) Start() error {
	var err error

	logs.Info().Str(logs.Component, logs.Core).Msg("starting core")

	for _, gateway := range core.gateways {
		err = gateway.Start()
		if err != nil {
			core.Stop()

			return err
		}
	}

	logs.Info().Str(logs.Component, logs.Core).Msg("successfully started core")

	return err
}

func (core *Core) Stop() {
	logs.Info().Str(logs.Component, logs.Core).Msg("stopping core")

	for _, gateway := range core.gateways {
		gateway.Stop()
	}

	logs.Info().Str(logs.Component, logs.Core).Msg("successfully stopped core")
}
