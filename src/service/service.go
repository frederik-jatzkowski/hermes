package service

import (
	"crypto/tls"
	"net"

	"fleo.software/infrastructure/hermes/balancer"
	"fleo.software/infrastructure/hermes/certbot"
	"fleo.software/infrastructure/hermes/logs"
)

type Service struct {
	ServerName *string                `xml:"servername,attr"`
	Balancer   *balancer.LoadBalancer `xml:"LoadBalancer"`
	KeyFile    *string                `xml:"keyfile,attr"`
	CertFile   *string                `xml:"certfile,attr"`
	Ok         bool                   `xml:"-"`
	cert       *tls.Certificate
}

func (s *Service) Init() {
	s.Ok = true // assert correct
	// check for all fields to be present and initialize them
	if s.ServerName == nil {
		logs.LaunchPrint("invalid service: missing server name", "3101")
		s.Ok = false // fatal
	}
	if s.Balancer == nil {
		logs.LaunchPrint("invalid service: missing load balancer", "3201")
		s.Ok = false // fatal
	}
	if s.Ok {
		cert, err := certbot.ObtainCertificate(*s.ServerName) // obtain certificate
		if err != nil {
			logs.LaunchPrint(err, "3301")
			s.Ok = false // fatal
		} else {
			s.cert = cert
		}
	}
	if !s.Ok {
		logs.BothPrint("invalid service '"+*s.ServerName+"' could not start operating", "3001") // log invalid service
	} else {
		s.Balancer.Init()
	}
}

func (s *Service) Handle(conn *net.Conn, tlsconn *tls.Conn) bool {
	// check if it can handle
	if s.Ok && (*tlsconn).ConnectionState().ServerName == *s.ServerName {
		// hand over to balancer
		if !s.Balancer.Handle(conn) {
			logs.Enumerator("failed services").Add(*s.ServerName)
		}
		return true
	}
	return false
}

func (s *Service) HandleClientHelloInfo(chi *tls.ClientHelloInfo) *tls.Certificate {
	if s.Ok {
		if chi.ServerName == *s.ServerName {
			return s.cert // return this services cert, if ServerNames match
		}
	}
	return nil
}
