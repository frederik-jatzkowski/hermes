package config

import "testing"

func TestServer_validate_Success(t *testing.T) {
	server := Server{
		Address: "0.0.0.0:8080",
	}

	err := server.validate()

	if err != nil {
		t.Errorf("unexpected error during balancer validation: %s", err)
	}
}

func TestServer_validate_EmptyAddress(t *testing.T) {
	server := Server{
		Address: "",
	}

	err := server.validate()

	if err == nil {
		t.Errorf("expected error during server validation, but got none")
	}
}

func TestServer_validate_InvalidAddress(t *testing.T) {
	server := Server{
		Address: "abc",
	}

	err := server.validate()

	if err == nil {
		t.Errorf("expected error during server validation, but got none")
	}
}
