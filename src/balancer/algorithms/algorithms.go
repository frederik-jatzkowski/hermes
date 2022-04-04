package algorithms

import (
	"errors"

	"fleo.software/infrastructure/hermes/server"
)

type BalancerAlgorithm interface {
	Next() *server.Server
}

func ResolveAlgorithm(name *string, servers *[]server.Server) (BalancerAlgorithm, error) {
	if len(*servers) == 0 {
		return nil, errors.New("not a single Server specified for LoadBalancer")
	} else if *name == "RoundRobin" || name == nil {
		return NewRoundRobin(servers), nil
	} else {
		return nil, errors.New("unknown load balancing alogrithm: '" + *name + "'. Available Algorithms: 'RoundRobin', 'LeastConnections'")
	}
}
