package config

import "testing"

func TestService_validate_Success(t *testing.T) {
	service := Service{
		HostName: "myname",
		Balancer: LoadBalancer{
			Algorithm: "RoundRobin",
		},
	}

	err := service.validate()

	if err != nil {
		t.Errorf("unexpected error during service validation: %s", err)
	}
}

func TestService_validate_InvalidLoadBalancer(t *testing.T) {
	service := Service{
		HostName: "myname",
		Balancer: LoadBalancer{},
	}

	err := service.validate()

	if err == nil {
		t.Errorf("expected error during service validation, but got none")
	}
}
