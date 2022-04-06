package balancer

import (
	"net"

	"fleo.software/infrastructure/hermes/balancer/algorithms"
	"fleo.software/infrastructure/hermes/logs"
	"fleo.software/infrastructure/hermes/server"
)

type LoadBalancer struct {
	Servers   []server.Server              `xml:"Server"`
	Algorithm *string                      `xml:"algorithm,attr"`
	algorithm algorithms.BalancerAlgorithm `xml:"-"`
	Ok        bool                         `xml:"-"`
}

func (b *LoadBalancer) Init() {
	b.Ok = true
	// init algorithm
	algo, err := algorithms.ResolveAlgorithm(b.Algorithm, &b.Servers)
	if err != nil {
		logs.LaunchPrint(err, "4101")
		b.Ok = false
	} else {
		b.algorithm = algo
	}
	// init servers
	for i := 0; i < len(b.Servers); i++ {
		b.Servers[i].Init()
	}
	if !b.Ok {
		logs.BothPrint("invalid load balancer could not start operating", "4001")
	}
}

func (b *LoadBalancer) Handle(conn *net.Conn) bool {
	if b.Ok {
		s := b.algorithm.Next() // get next server
		if s != nil {
			s.Handle(conn)
			return true
		}
		(*conn).Close()
	}
	return false
}
