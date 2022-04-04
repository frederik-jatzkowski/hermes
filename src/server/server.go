package server

import (
	"crypto/tls"
	"log"
	"net"
	"sync"
	"time"

	"fleo.software/infrastructure/hermes/logging/startup"
)

const RETRY_INTERVAL = 5 * time.Second
const HEALTHCHECK_INTERVAL = 100 * time.Millisecond

type dial func(address *string) (net.Conn, error)

type Server struct {
	Address    *string      `xml:"raddress,attr"`
	Secure     bool         `xml:"secure,attr"`
	TCPAddress *net.TCPAddr `xml:"-"`
	online     bool         `xml:"-"`
	checkConn  *net.Conn    `xml:"-"`
	mutex      sync.Mutex   `xml:"-"`
	dial       dial
}

func (s *Server) Init(collector *startup.ErrorCollector) {
	// check for all fields to be present
	if s.Address == nil {
		collector.Error("Attribute 'address' of Server unspecified.")
	} else {
		addr, err := net.ResolveTCPAddr("tcp", *s.Address)
		collector.Append(err)
		s.TCPAddress = addr
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
	// init pool
	go s.monitor()
}
func (s *Server) CanHandle() bool {
	s.mutex.Lock()
	online := s.online
	s.mutex.Unlock()
	return online
}
func (s *Server) Handle(clientConn *net.Conn) {
	serverConn, err := s.tryConn()
	if err == nil {
		go newPipe(clientConn, serverConn, s).start()
	} else {
		(*clientConn).Close()
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
				log.Printf("server '%v' went offline", *s.Address)
			} else {
				log.Printf("server '%v' came online", *s.Address)
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
