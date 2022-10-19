package config

import "testing"

func TestLoadBalancer_validate_Success(t *testing.T) {
	balancer := LoadBalancer{
		Algorithm: "RoundRobin",
	}

	err := balancer.validate()

	if err != nil {
		t.Errorf("unexpected error during balancer validation: %s", err)
	}
}

func TestLoadBalancer_validate_InvalidAlgorithm(t *testing.T) {
	balancer := LoadBalancer{
		Algorithm: "SquareRobin",
	}

	err := balancer.validate()

	if err == nil {
		t.Errorf("expected error during balancer validation, but got none")
	}
}

func TestLoadBalancer_validate_InvalidServer(t *testing.T) {
	balancer := LoadBalancer{
		Algorithm: "RoundRobin",
		Servers:   []Server{{Address: "abc"}},
	}

	err := balancer.validate()

	if err == nil {
		t.Errorf("expected error during balancer validation, but got none")
	}
}
