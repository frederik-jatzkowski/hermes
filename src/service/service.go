package service

import (
	"crypto/tls"
	"net"
	"time"

	"fleo.software/infrastructure/hermes/balancer"
	"fleo.software/infrastructure/hermes/logging"
	"fleo.software/infrastructure/hermes/logging/startup"
)

type Service struct {
	HostName         *string                `xml:"hostname,attr"`
	Balancer         *balancer.LoadBalancer `xml:"LoadBalancer"`
	KeyFile          *string                `xml:"keyfile,attr"`
	CertFile         *string                `xml:"certfile,attr"`
	cert             *tls.Certificate       `xml,"-"`
	failedLogCounter *logging.LogCounter    `xml,"-"`
}

func (s *Service) Init(collector *startup.ErrorCollector) {
	// check for all fields to be present and initialize them
	if s.HostName == nil {
		collector.Error("No hostname for Service unspecified.")
	}
	if s.Balancer == nil {
		collector.Error("No Balancer for Service specified.")
	} else {
		s.Balancer.Init(collector)
	}
	if s.KeyFile == nil {
		collector.Error("No keyfile for Service specified.")
	}
	if s.CertFile == nil {
		collector.Error("No certfile for Service specified.")
	}
	if s.CertFile != nil && s.KeyFile != nil {
		// load cert if present
		cert, err := tls.LoadX509KeyPair(*s.CertFile, *s.KeyFile)
		if err != nil {
			collector.Append(err)
		} else {
			s.cert = &cert
		}
	}
	s.failedLogCounter = logging.NewLogCounter("service: '"+*s.HostName+"', failed: %v", time.Minute, false)
}

func (s *Service) Handle(conn *net.Conn, tlsconn *tls.Conn) bool {
	if (*tlsconn).ConnectionState().ServerName != *s.HostName {
		return false
	}
	if !s.Balancer.Handle(conn) {
		s.failedLogCounter.Increment()
	}
	return true
}

func (s *Service) HandleClientHelloInfo(chi *tls.ClientHelloInfo) *tls.Certificate {
	if chi.ServerName == *s.HostName {
		return s.cert
	}
	return nil
}
