package server

import (
	"net"
)

func newPipe(clientConn *net.Conn, serverConn *net.Conn, server *Server) *pipe {
	return &pipe{
		clientConn: clientConn,
		serverConn: serverConn,
		server:     server,
	}
}

type pipe struct {
	clientConn *net.Conn
	serverConn *net.Conn
	server     *Server
}

func (p *pipe) start() {
	// after connection fails, close pipeline
	defer p.close()
	// start piplining in one direction
	go p.pipe(p.clientConn, p.serverConn)
	// start pipelining in the other direction
	p.pipe(p.serverConn, p.clientConn)
}

func (p *pipe) pipe(source *net.Conn, target *net.Conn) {
	b := make([]byte, 8)
	// loop until one connection fails
	for {
		i, errIn := (*source).Read(b)
		_, errOut := (*target).Write(b[0:i])
		if errIn != nil || errOut != nil {
			break
		}
	}
}

func (p *pipe) close() {
	// close both connections
	(*p.clientConn).Close()
	(*p.serverConn).Close()
}
