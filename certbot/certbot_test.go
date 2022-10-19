package certbot

import (
	"crypto/tls"
	"fmt"
	"io/fs"
	"os/exec"
	"strings"
	"testing"
)

type mockDirEntry struct {
	name     string
	isDir    bool
	fileMode fs.FileMode
	info     fs.FileInfo
	err      error
}

func (entry mockDirEntry) Name() string {
	return entry.name
}
func (entry mockDirEntry) IsDir() bool {
	return entry.isDir
}
func (entry mockDirEntry) Type() fs.FileMode {
	return entry.fileMode
}
func (entry mockDirEntry) Info() (fs.FileInfo, error) {
	return entry.info, entry.err
}

func TestFindCertPath_Success(t *testing.T) {
	var hostName = "hermes.org"

	osReadDir = func(name string) ([]fs.DirEntry, error) {
		return []fs.DirEntry{
			mockDirEntry{name: hostName, isDir: true},
		}, nil
	}

	certFile, keyFile, err := FindCertPath(hostName)

	if err != nil {
		t.Errorf("unexpected error during finding of cert path: %s", err)
	}

	if !strings.Contains(certFile, hostName) {
		t.Errorf("expected certFile '%s' to contain '%s'", certFile, hostName)
	}

	if !strings.Contains(keyFile, hostName) {
		t.Errorf("expected keyFile '%s' to contain '%s'", keyFile, hostName)
	}
}

func TestFindCertPath_ReadFailed(t *testing.T) {
	var hostName = "hermes.org"

	osReadDir = func(name string) ([]fs.DirEntry, error) {
		return []fs.DirEntry{}, fmt.Errorf("test error")
	}

	certFile, keyFile, err := FindCertPath(hostName)

	if err == nil {
		t.Errorf("expected error during finding of cert path but got '%s' and '%s' instead", certFile, keyFile)
	}
}

func TestFindCertPath_NotFound(t *testing.T) {
	var hostName = "hermes.org"

	osReadDir = func(name string) ([]fs.DirEntry, error) {
		return []fs.DirEntry{}, nil
	}

	certFile, keyFile, err := FindCertPath(hostName)

	if err == nil {
		t.Errorf("expected error during finding of cert path but got '%s' and '%s' instead", certFile, keyFile)
	}
}

func TestObtainCertificate_Success(t *testing.T) {
	var hostName = "hermes.org"

	osReadDir = func(name string) ([]fs.DirEntry, error) {
		return []fs.DirEntry{
			mockDirEntry{name: hostName, isDir: true},
		}, nil
	}

	tlsLoadX509KeyPair = func(certFile, keyFile string) (tls.Certificate, error) {
		return tls.Certificate{}, nil
	}

	_, err := ObtainCertificate(hostName)

	if err != nil {
		t.Errorf("unexpected error while obtaining certificate: %s", err)
	}
}

func TestObtainCertificate_Certbot_Success(t *testing.T) {
	var hostName = "hermes.org"

	counter := 0
	osReadDir = func(name string) ([]fs.DirEntry, error) {
		if counter == 0 {
			counter++
			return []fs.DirEntry{}, nil
		} else {
			return []fs.DirEntry{
				mockDirEntry{name: hostName, isDir: true},
			}, nil
		}
	}

	tlsLoadX509KeyPair = func(certFile, keyFile string) (tls.Certificate, error) {
		return tls.Certificate{}, nil
	}

	execCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", "")
	}

	_, err := ObtainCertificate(hostName)

	if err != nil {
		t.Errorf("unexpected error while obtaining certificate: %s", err)
	}
}

func TestObtainCertificate_Certbot_Failure(t *testing.T) {
	var hostName = "hermes.org"

	counter := 0
	osReadDir = func(name string) ([]fs.DirEntry, error) {
		if counter == 0 {
			counter++
			return []fs.DirEntry{}, nil
		} else {
			return []fs.DirEntry{
				mockDirEntry{name: hostName, isDir: true},
			}, nil
		}
	}

	tlsLoadX509KeyPair = func(certFile, keyFile string) (tls.Certificate, error) {
		return tls.Certificate{}, nil
	}

	execCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("exit", "1")
	}

	_, err := ObtainCertificate(hostName)

	if err == nil {
		t.Errorf("expected error while obtaining certificate but got none")
	}
}

func TestObtainCertificate_CertbotWroteNoKeyPair(t *testing.T) {
	var hostName = "hermes.org"

	osReadDir = func(name string) ([]fs.DirEntry, error) {
		return []fs.DirEntry{}, nil
	}

	tlsLoadX509KeyPair = func(certFile, keyFile string) (tls.Certificate, error) {
		return tls.Certificate{}, nil
	}

	execCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", "")
	}

	_, err := ObtainCertificate(hostName)

	if err == nil {
		t.Errorf("expected error while obtaining certificate but got none")
	}
}

func TestObtainCertificate_InvalidKeyPairFiles(t *testing.T) {
	var hostName = "hermes.org"

	counter := 0
	osReadDir = func(name string) ([]fs.DirEntry, error) {
		if counter == 0 {
			counter++
			return []fs.DirEntry{}, nil
		} else {
			return []fs.DirEntry{
				mockDirEntry{name: hostName, isDir: true},
			}, nil
		}
	}

	tlsLoadX509KeyPair = func(certFile, keyFile string) (tls.Certificate, error) {
		return tls.Certificate{}, fmt.Errorf("test error")
	}

	execCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", "")
	}

	_, err := ObtainCertificate(hostName)

	if err == nil {
		t.Errorf("expected error while obtaining certificate but got none")
	}
}
