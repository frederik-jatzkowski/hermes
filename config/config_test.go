package config

import (
	"testing"
)

func TestConfig_NewConfig_Success(t *testing.T) {
	_, err := NewConfig([]byte("{\"email\":\"my@email.org\", \"logLevel\":\"trace\", \"gateways\":[]}"))

	if err != nil {
		t.Errorf("unexpected error during config building: %s", err)
	}
}

func TestConfig_NewConfig_ParseFailure(t *testing.T) {
	_, err := NewConfig([]byte("<Config></gifnoC>"))

	if err == nil {
		t.Errorf("expected error during config building, but got none")
	}
}

func TestConfig_NewConfig_ValidateFailure(t *testing.T) {
	_, err := NewConfig([]byte("{}"))

	if err == nil {
		t.Errorf("expected error during config building, but got none")
	}
}

func TestConfig_validate_InvalidLogLevel(t *testing.T) {
	_, err := NewConfig([]byte("{\"email\":\"my@email.org\", \"logLevel\":\"tarce\", \"gateways\":[]}"))

	if err == nil {
		t.Errorf("expected error during config building, but got none")
	}
}

func TestConfig_validate_InvalidGateway(t *testing.T) {
	_, err := NewConfig([]byte("{\"email\":\"my@email.org\", \"logLevel\":\"trace\", \"gateways\":[{}]}"))

	if err == nil {
		t.Errorf("expected error during config building, but got none")
	}
}
