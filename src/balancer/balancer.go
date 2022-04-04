package balancer

import (
	"net"

	"fleo.software/infrastructure/hermes/balancer/algorithms"
	"fleo.software/infrastructure/hermes/logging/startup"
	"fleo.software/infrastructure/hermes/server"
)

type LoadBalancer struct {
	Servers   []server.Server              `xml:"Server"`
	Algorithm *string                      `xml:"algorithm,attr"`
	algorithm algorithms.BalancerAlgorithm `xml:"-"`
}

func (lb *LoadBalancer) Init(collector *startup.ErrorCollector) {
	// initialize balancing algorithm
	if lb.Algorithm == nil {
		algo := "RoundRobin"
		lb.Algorithm = &algo
	}
	algo, err := algorithms.ResolveAlgorithm(lb.Algorithm, &lb.Servers)
	if err == nil {
		lb.algorithm = algo
	} else {
		collector.Append(err)
	}

	// init servers
	for i := 0; i < len(lb.Servers); i++ {
		lb.Servers[i].Init(collector)
	}
}

func (lb *LoadBalancer) Handle(conn *net.Conn) bool {
	s := lb.algorithm.Next()
	if s != nil {
		s.Handle(conn)
		return true
	}
	(*conn).Close()
	return false
}
