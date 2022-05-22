package server

import (
	"crypto/tls"
	"io"
	"net"
	"sync"
	"time"

	"fleo.software/infrastructure/hermes/logs"
)

const RETRY_INTERVAL = 5 * time.Second
const HEALTHCHECK_INTERVAL = 100 * time.Millisecond

type dial func(address *string) (net.Conn, error)

type Server struct {
	Address    *string      `xml:"raddress,attr"`
	Secure     bool         `xml:"tls,attr"`
	tcpAddress *net.TCPAddr `xml:"-"`
	Ok         bool         `xml:"-"`
	online     bool         `xml:"-"`
	checkConn  *net.Conn    `xml:"-"`
	mutex      sync.Mutex   `xml:"-"`
	dial       dial         `xml:"-"`
}

func (s *Server) Init() {
	s.Ok = true // assume correct
	// check for all fields to be present
	if s.Address == nil {
		logs.LaunchPrint("invalid server: missing remote address", "5101")
		s.Ok = false // fatal
	} else {
		addr, err := net.ResolveTCPAddr("tcp", *s.Address)
		if err != nil {
			logs.LaunchPrint(err, "5201")
			s.Ok = false // fatal
		}
		s.tcpAddress = addr
	}
	// init dial function
	if s.Secure {
		s.dial = func(address *string) (net.Conn, error) {
			return tls.Dial("tcp", *address, nil)
		}
	} else {
		s.dial = func(address *string) (net.Conn, error) {
			return net.Dial("tcp", *address)
		}
	}

	if s.Ok {
		go s.monitor() // monitor server
	} else {
		logs.BothPrint("invalid server: '"+*s.Address+"' could not start operating", "5001")
	}
}
func (s *Server) CanHandle() bool {
	s.mutex.Lock()
	online := s.online
	s.mutex.Unlock()
	return online
}
func (s *Server) Handle(clientConn *net.Conn) {
	defer (*clientConn).Close()
	serverConn, err := s.tryConn()
	if err == nil {
		//defer (*serverConn).Close()
		go func() {
			defer (*serverConn).Close()
			io.Copy(*clientConn, *serverConn)
		}()
		//go io.Copy(*clientConn, *serverConn)
		io.Copy(*serverConn, *clientConn)
	}
}

// utility functions
func (s *Server) monitor() {
	// health checks
	var online bool
	for {
		s.mutex.Lock()
		online = s.online
		s.mutex.Unlock()
		conn, err := s.tryConn()
		// check state changes
		if online != (err == nil) {
			if online {
				logs.ContinuousPrint("server '"+*s.Address+"' went offline", "5011")
			} else {
				logs.ContinuousPrint("server '"+*s.Address+"' came online", "5010")
			}
		}
		if err == nil {
			(*conn).Close()
			time.Sleep(HEALTHCHECK_INTERVAL)
		} else {
			time.Sleep(RETRY_INTERVAL)
		}
	}
}
func (s *Server) tryConn() (*net.Conn, error) {
	conn, err := s.dial(s.Address)
	s.mutex.Lock()
	s.online = err == nil
	s.mutex.Unlock()
	return &conn, err
}
