package config

import "fmt"

type Service struct {
	HostName string       `json:"hostName"`
	Balancer LoadBalancer `json:"balancer"`
}

func (service Service) validate() error {
	var (
		err error
	)

	// TODO: validate server name

	// validate load balancer
	err = service.Balancer.validate()
	if err != nil {
		return fmt.Errorf("invalid load balancer config: %s", err)
	}

	return err
}
