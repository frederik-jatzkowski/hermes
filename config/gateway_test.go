package config

import (
	"testing"
)

func TestGateway_validate_Success(t *testing.T) {
	gateway := Gateway{
		Services: []Service{},
		Address:  "0.0.0.0:8080",
	}

	err := gateway.validate()

	if err != nil {
		t.Errorf("unexpected error during gateway validation: %s", err)
	}
}

func TestGateway_validate_EmptyAddress(t *testing.T) {
	gateway := Gateway{
		Services: []Service{},
		Address:  "",
	}

	err := gateway.validate()

	if err == nil {
		t.Errorf("expected error during gateway validation but got none")
	}
}

func TestGateway_validate_InvalidAddress(t *testing.T) {
	gateway := Gateway{
		Services: []Service{},
		Address:  "abc",
	}

	err := gateway.validate()

	if err == nil {
		t.Errorf("expected error during gateway validation but got none")
	}
}

func TestGateway_validate_InvalidService(t *testing.T) {
	gateway := Gateway{
		Services: []Service{{}},
		Address:  "0.0.0.0:8080",
	}

	err := gateway.validate()

	if err == nil {
		t.Errorf("expected error during config building, but got none")
	}
}
