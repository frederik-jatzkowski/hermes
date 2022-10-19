package certbot

import (
	"crypto/tls"
	"io/fs"
	"os"
	"os/exec"
)

var osReadDir func(name string) ([]fs.DirEntry, error) = os.ReadDir

var tlsLoadX509KeyPair func(certFile string, keyFile string) (tls.Certificate, error) = tls.LoadX509KeyPair

var execCommand func(name string, arg ...string) *exec.Cmd = exec.Command
