package gateway

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"fleo.software/infrastructure/hermes/logging"
	"fleo.software/infrastructure/hermes/logging/startup"
	"fleo.software/infrastructure/hermes/service"
)

type Gateway struct {
	Address  *string           `xml:"laddress,attr"`
	Services []service.Service `xml:"Service"`
	address  *net.TCPAddr      `xml:"-"`
}

func (g *Gateway) Init(collector *startup.ErrorCollector) {
	if g.Address == nil {
		collector.Error("no raddress for Service specified.")
	} else {
		// check address
		addr, err := net.ResolveTCPAddr("tcp", *g.Address)
		g.address = addr
		collector.Append(err)
	}
	// init services
	for i := 0; i < len(g.Services); i++ {
		g.Services[i].Init(collector)
	}
}

func (g *Gateway) Listen() {
	lc1 := logging.NewLogCounter("gateway '"+*g.Address+"' accepted %v connection(s) with unknown hostnames: %v", time.Minute, true)
	cfg := &tls.Config{
		GetCertificate: func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
			var cert *tls.Certificate
			for i := 0; i < len(g.Services); i++ {
				cert = g.Services[i].HandleClientHelloInfo(chi)
				if cert != nil {
					return cert, nil
				}
			}
			lc1.IdentifierIncrement(chi.ServerName)
			return nil, errors.New("service not found")
		},
	}

	l, err := tls.Listen("tcp", *g.Address, cfg)
	if err != nil {
		fmt.Println(err)
		log.Fatalln(err)
	}
	lc2 := logging.NewLogCounter("gateway '"+*g.Address+"' accepted %v connections in the last minute", time.Minute, false)
	for {
		conn, err := l.Accept()
		go func() {
			conn.(*tls.Conn).Handshake()
			if err != nil {
				log.Println(err)
			}
			for i := 0; i < len(g.Services); i++ {
				if g.Services[i].Handle(&conn, conn.(*tls.Conn)) {
					break
				}
			}
			lc2.Increment() // after connected, increment counter
		}()

	}
}
