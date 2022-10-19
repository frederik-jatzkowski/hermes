package core

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/frederik-jatzkowski/hermes/config"
	"github.com/frederik-jatzkowski/hermes/logs"
	"github.com/frederik-jatzkowski/hermes/params"
)

type Server struct {
	address    net.TCPAddr
	healthy    bool
	healthLock sync.RWMutex
	conns      map[*net.TCPConn]bool
	connsLock  sync.Mutex
	stopLock   sync.RWMutex // will be hold during the closing process to stop new connections being build
}

func NewServer(config config.Server) *Server {
	return &Server{
		address: config.ResolvedAddress,
		healthy: false,
		conns:   make(map[*net.TCPConn]bool),
	}
}

func (server *Server) Handle(clientConn *net.Conn) error {
	var (
		serverConn *net.TCPConn
		err        error
	)

	// make sure, server is not beeing stopped right now
	stopping := !server.stopLock.TryRLock()
	if stopping {
		// if closing, close incoming
		(*clientConn).Close()

		// do not connect to server
		return nil
	}

	// check health info
	server.healthLock.RLock()
	if !server.healthy {
		server.healthLock.RUnlock()
		server.stopLock.RUnlock()

		return fmt.Errorf("server not available")
	}
	server.healthLock.RUnlock()

	// open connection
	serverConn, err = net.DialTCP("tcp", nil, &server.address)
	if err != nil {
		server.stopLock.RUnlock()

		return fmt.Errorf("error while connecting to server: %s", err)
	}

	// start transmission
	go server.transmit(clientConn, serverConn)

	return err
}

func (server *Server) transmit(clientConn *net.Conn, serverConn *net.TCPConn) {
	logs.Debug().Str(logs.Component, logs.Server).Str(logs.ServerAddress, server.address.String()).
		Str(logs.ClientAddress, (*clientConn).RemoteAddr().String()).Msg("transmission started")

	// transmission
	go server.closeCopy(*clientConn, serverConn)
	server.closeCopy(serverConn, *clientConn)

	// remove connection from server
	server.connsLock.Lock()
	delete(server.conns, serverConn)
	server.connsLock.Unlock()

	// end of lifecycle
	logs.Debug().Str(logs.Component, logs.Server).Str(logs.ServerAddress, server.address.String()).
		Str(logs.ClientAddress, (*clientConn).RemoteAddr().String()).Msg("transmission stopped")
}

func (server *Server) closeCopy(src net.Conn, dst net.Conn) {
	// copy, until one connection fails or is closed
	bytes, err := io.Copy(src, dst)
	if err != nil {
		logs.Error().Str(logs.Component, logs.Server).Msgf("transmission failed after %d bytes: %s", bytes, err)
	}

	// close destination
	dst.Close()
}

func (server *Server) Start() {
	go server.monitor()
}

func (server *Server) monitor() {
	var (
		err error
	)

	// check forever until server is being stopped
	for {
		// if stopping, seize monitoring
		stopping := !server.stopLock.TryRLock()
		if stopping {
			break
		}

		// try to connect
		_, err = net.DialTCP("tcp", nil, &server.address)
		if err != nil {
			// discover, that server is unavailable
			server.healthLock.Lock()
			if server.healthy {
				logs.Error().Str(logs.Component, logs.Server).Str(logs.ServerAddress, server.address.String()).Msg("server unavailable")
			}
			server.healthy = false
			server.healthLock.Unlock()
		} else {
			// discover, that server is available
			server.healthLock.Lock()
			if !server.healthy {
				logs.Info().Str(logs.Component, logs.Server).Str(logs.ServerAddress, server.address.String()).Msg("server available")
			}
			server.healthy = true
			server.healthLock.Unlock()
		}
		server.stopLock.RUnlock()

		// wait for the specified duration
		time.Sleep(params.HealthCheckInterval)
	}
}

func (server *Server) Stop() {
	// prevent further connections to be handled
	server.stopLock.Lock()

	// close all existing connections
	server.connsLock.Lock()
	for conn := range server.conns {
		conn.Close()
	}
	server.connsLock.Unlock()

	logs.Debug().Str(logs.Component, logs.Server).Str(logs.ServerAddress, server.address.String()).Msg("server stopped")
}
