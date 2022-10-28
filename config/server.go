package config

import (
	"fmt"
	"net"
)

type Server struct {
	Address         string      `json:"address"`
	ResolvedAddress net.TCPAddr `json:"-"`
}

func (server *Server) validate() error {
	var (
		resolvedAddress *net.TCPAddr
		err             error
	)

	// validate address
	if server.Address == "" {
		return fmt.Errorf("missing server address")
	}
	resolvedAddress, err = net.ResolveTCPAddr("tcp", server.Address)
	if err != nil {
		return fmt.Errorf("invalid server address: %s", err)
	}
	server.ResolvedAddress = *resolvedAddress

	return err
}
