package config

import (
	"fmt"
	"net"
)

type Gateway struct {
	Services        []Service   `json:"services"`
	Address         string      `json:"address"`
	ResolvedAddress net.TCPAddr `json:"-"`
}

func (gateway *Gateway) validate() error {
	var (
		resolvedAddress *net.TCPAddr
		err             error
	)

	// validate address
	if gateway.Address == "" {
		return fmt.Errorf("missing gateway address")
	}
	resolvedAddress, err = net.ResolveTCPAddr("tcp", gateway.Address)
	if err != nil {
		return fmt.Errorf("invalid gateway address: %s", err)
	}
	gateway.ResolvedAddress = *resolvedAddress

	// certbot needs port 442
	if resolvedAddress.Port == 442 {
		return fmt.Errorf("port 442 is reserved for certbot")
	}

	// port 442 should not be used over tls anyway, it is standard for unencrypted http
	if resolvedAddress.Port == 80 {
		return fmt.Errorf("port 80 is reserved for http-to-https-redirects")
	}

	// the admin panel needs port 440
	if resolvedAddress.Port == 440 {
		return fmt.Errorf("port '%d' is reserved for the admin panel", resolvedAddress.Port)
	}

	// the admin panel needs port 441
	if resolvedAddress.Port == 441 {
		return fmt.Errorf("port '%d' is reserved for the admin panel", resolvedAddress.Port)
	}

	// validate Services
	for i, service := range gateway.Services {
		err = service.validate()
		if err != nil {
			return fmt.Errorf("invalid service config: %s", err)
		}
		gateway.Services[i] = service
	}

	return err
}
