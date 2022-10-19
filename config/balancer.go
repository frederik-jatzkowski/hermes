package config

import "fmt"

type LoadBalancer struct {
	Servers   []Server `json:"servers"`
	Algorithm string   `json:"algorithm"`
}

func (balancer LoadBalancer) validate() error {
	var err error

	// validate Algorithm
	if balancer.Algorithm != "RoundRobin" {
		return fmt.Errorf("algorithm has to be one of {'RoundRobin'} but was '%s'", balancer.Algorithm)
	}

	// validate Servers
	for _, server := range balancer.Servers {
		err = server.validate()
		if err != nil {
			return fmt.Errorf("invalid server config: %s", err)
		}
	}

	return err
}
