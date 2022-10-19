package config

import (
	"fmt"
	"testing"
)

func TestConfig_ReadConfigFile_Success(t *testing.T) {
	osReadFile = func(name string) ([]byte, error) {
		return []byte("abc"), nil
	}

	data, err := ReadConfigFile()

	if err != nil {
		t.Errorf("unexpected error during reading of config file: %s", err)
	}

	expected := len([]byte("abc"))
	if len(data) != expected {
		t.Errorf("expected '%d' bytes to be read but got '%d'", expected, len(data))
	}
}

func TestConfig_ReadConfigFile_Failure(t *testing.T) {
	osReadFile = func(name string) ([]byte, error) {
		return []byte{}, fmt.Errorf("test error")
	}

	_, err := ReadConfigFile()

	if err == nil {
		t.Errorf("expected error during reading of config file but got none")
	}

}

func TestConfig_NewConfig_Success(t *testing.T) {
	config, err := NewConfig([]byte("{\"email\":\"my@email.org\", \"logLevel\":\"trace\", \"gateways\":[]}"))

	if err != nil {
		t.Errorf("unexpected error during config building: %s", err)
	}

	expected := "trace"
	if config.LogLevel != expected {
		t.Errorf("expected '%s' as config.LogLevel but got '%s'", expected, config.LogLevel)
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
