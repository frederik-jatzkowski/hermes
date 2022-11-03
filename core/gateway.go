package core

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/frederik-jatzkowski/hermes/config"
	"github.com/frederik-jatzkowski/hermes/logs"
)

type Gateway struct {
	services map[string]*Service
	listener net.Listener
	address  net.TCPAddr
}

func NewGateway(config config.Gateway) *Gateway {
	var gateway = Gateway{
		address: config.ResolvedAddress,
	}

	gateway.services = make(map[string]*Service)
	for _, serviceConfig := range config.Services {
		gateway.services[serviceConfig.HostName] = NewService(serviceConfig)
	}

	return &gateway
}

func (gateway *Gateway) Start() error {
	var (
		listener net.Listener
		err      error
	)

	logs.Info().Str(logs.Component, logs.Gateway).Int(logs.Port, gateway.address.Port).Msg("starting gateway")

	// create listener
	listener, err = tls.Listen(
		"tcp",
		gateway.address.String(),
		&tls.Config{
			GetCertificate: gateway.handleChi,
		},
	)
	if err != nil {
		return fmt.Errorf("could not start listening on address '%s': %s", gateway.address.String(), err)
	}
	gateway.listener = listener

	// start services
	for _, service := range gateway.services {
		err = service.Start()
		if err != nil {
			return fmt.Errorf("could not start service on gateway with address '%s': %s", gateway.address.String(), err)
		}
	}

	// start listening
	go gateway.listen()

	logs.Info().Str(logs.Component, logs.Gateway).Int(logs.Port, gateway.address.Port).Msg("gateway successfully started")

	return err
}

func (gateway *Gateway) handleChi(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
	service, exists := gateway.services[chi.ServerName]
	if !exists {
		logs.Debug().Str(logs.Component, logs.Gateway).Str(logs.Sni, chi.ServerName).
			Msg("unexpected chi server name")

		return nil, fmt.Errorf("unexpected chi server name")
	}
	logs.Debug().Str(logs.Component, logs.Gateway).Str(logs.Sni, chi.ServerName).
		Msg("handled chi")

	return service.Cert(), nil
}

func (gateway *Gateway) listen() {
	var (
		service *Service
		exists  bool
		tlsConn *tls.Conn
	)

	for {
		logs.Debug().Str(logs.Component, logs.Gateway).Int(logs.Port, gateway.address.Port).Msg("accepting next connection")
		conn, err := gateway.listener.Accept()
		if err != nil {
			logs.Debug().Str(logs.Component, logs.Gateway).Int(logs.Port, gateway.address.Port).Msgf("error while accepting connection: %s", err)

			break
		}
		logs.Debug().Str(logs.Component, logs.Gateway).Int(logs.Port, gateway.address.Port).Msg("accepted connection")

		tlsConn = conn.(*tls.Conn)
		service, exists = gateway.services[tlsConn.ConnectionState().ServerName]
		if exists {
			go service.Handle(&conn)
		} else {
			logs.Debug().Str(logs.Component, logs.Gateway).Str(logs.Sni, tlsConn.ConnectionState().ServerName).
				Msg("unexpected connection state server name")

			conn.Close()
		}
	}
}

func (gateway *Gateway) Stop() {
	logs.Info().Str(logs.Component, logs.Gateway).Int(logs.Port, gateway.address.Port).Msg("stopping gateway")

	if gateway.listener != nil {
		gateway.listener.Close()
	}

	for _, service := range gateway.services {
		service.Stop()
	}

	logs.Info().Str(logs.Component, logs.Gateway).Int(logs.Port, gateway.address.Port).Msg("successfully stopped gateway")
}
