package core

import (
	"fmt"
	"net"
	"sync"

	"github.com/frederik-jatzkowski/hermes/config"
)

type LoadBalancer interface {
	Start()
	Stop()
	Handle(conn *net.Conn) error
}

func NewLoadBalancer(config config.LoadBalancer) LoadBalancer {
	var (
		servers []*Server
		server  *Server
	)

	// construct servers
	for _, serverConfig := range config.Servers {
		server = NewServer(serverConfig)
		servers = append(servers, server)
	}

	// select and build algorithm
	switch config.Algorithm {
	default:
		return &roundRobinBalancer{
			servers: servers,
			next:    0,
			lock:    &sync.Mutex{},
		}
	}
}

type roundRobinBalancer struct {
	servers []*Server
	next    int
	lock    *sync.Mutex
}

func (balancer *roundRobinBalancer) Handle(conn *net.Conn) error {
	var (
		server *Server
		err    error
	)

	balancer.lock.Lock()
	defer balancer.lock.Unlock()

	var start = balancer.next
	for {
		server = balancer.servers[balancer.next]
		balancer.next = (balancer.next + 1) % len(balancer.servers)

		err = server.Handle(conn)

		if err == nil {
			return nil
		}
		if balancer.next == start {
			return fmt.Errorf("no server could handle the connection")
		}
	}
}

func (balancer *roundRobinBalancer) Start() {
	for _, server := range balancer.servers {
		server.Start()
	}
}

func (balancer *roundRobinBalancer) Stop() {
	for _, server := range balancer.servers {
		server.Stop()
	}
}
