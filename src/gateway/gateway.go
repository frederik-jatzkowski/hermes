package gateway

import (
	"crypto/tls"
	"errors"
	"log"
	"net"

	"fleo.software/infrastructure/hermes/logs"
	"fleo.software/infrastructure/hermes/service"
)

type Gateway struct {
	Address  *string           `xml:"laddress,attr"`
	Services []service.Service `xml:"Service"`
	Ok       bool              `xml:"-"`
	address  *net.TCPAddr      `xml:"-"`
}

func (g *Gateway) Init() {
	g.Ok = true // assume correct
	if g.Address == nil {
		logs.LaunchPrint("invalid gateway: missing local adddress", "2101")
		g.Ok = false // fatal
	} else {
		addr, err := net.ResolveTCPAddr("tcp", *g.Address) // check address
		if err != nil {
			logs.LaunchPrint(err, "2201")
			g.Ok = false // fatal
		}
		g.address = addr
	}
	if !g.Ok { // if gateway invalid, log that
		logs.BothPrint("invalid gateway '"+*g.Address+"' could not start operating", "2001")
	} else {
		for i := 0; i < len(g.Services); i++ {
			g.Services[i].Init() // init services
		}
	}
}

func (g *Gateway) Listen() {
	if !g.Ok {
		return // if invalid gateway, dont listen
	}
	// set cfg for this gateway
	cfg := &tls.Config{
		GetCertificate: func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
			for i := 0; i < len(g.Services); i++ {
				cert := g.Services[i].HandleClientHelloInfo(chi) // find service to provide cert
				if cert != nil {
					return cert, nil
				}
			}
			// if no service found
			logs.Enumerator("unknown hostnames").Add(chi.ServerName) // log unknown server name in ClientHelloInfo
			return nil, errors.New("no suitable service found")      // return error
		},
	}
	// create listener
	l, err := tls.Listen("tcp", *g.Address, cfg)
	if err != nil {
		logs.BothPrint(err, "2301")
		logs.BothPrint("invalid gateway '"+*g.Address+"' could not start operating", "2301")
		return // if listener fails, gateway cannot start
	}
	for {
		// accept forever
		conn, err := l.Accept()
		if err != nil {
			log.Println(err) // log failures during accept
		}
		go func() {
			// perform handshake to get ServerName and log failures
			err := conn.(*tls.Conn).Handshake()
			if err != nil {
				logs.ContinuousPrint(err, "2401")
				conn.Close() // if handshake error, close connection
				return
			}
			for i := 0; i < len(g.Services); i++ {
				if g.Services[i].Handle(&conn, conn.(*tls.Conn)) {
					break // find service for handling
				}
			}
			// count connections
			logs.Counter("inbound connections").Increment()
		}()

	}
}
