package certbot

import (
	"crypto/tls"
	"fmt"
	"io/fs"
	"strings"
	"sync"
	"time"

	"github.com/frederik-jatzkowski/hermes/logs"
	"github.com/frederik-jatzkowski/hermes/params"
)

func FindCertPath(hostName string) (string, string, error) {
	var (
		dirName = "/etc/letsencrypt/live/"
		dir     []fs.DirEntry
		err     error
	)

	// read directory with certbot-certificates
	dir, err = osReadDir(dirName)
	if err != nil {
		return "", "", fmt.Errorf("could not read directory '%s'", dirName)
	}

	// find first matching subdirectory
	for _, entry := range dir {
		if strings.HasPrefix(entry.Name(), hostName) && entry.IsDir() {
			path := dirName + entry.Name() + "/"
			return path + "fullchain.pem", path + "privkey.pem", err
		}
	}

	// if none found, return error
	return "", "", fmt.Errorf("no certificate directory found for server name '%s'", hostName)
}

func ObtainCertificate(hostName string) (tls.Certificate, error) {
	var (
		cert tls.Certificate
		err  error
		out  []byte
	)

	if hostName == "" {
		return tlsLoadX509KeyPair("/etc/letsencrypt/live/localhost/fullchain.pem", "/etc/letsencrypt/live/localhost/privkey.pem")
	}

	// check, if certificate already exists
	certFile, keyFile, err := FindCertPath(hostName)
	if err == nil {
		cert, err = tlsLoadX509KeyPair(certFile, keyFile)
		if err == nil {
			// if loading worked, just return certificate
			return cert, nil
		}
	}

	// if cert for hostname does not exist, obtain certificate from certbot
	command := execCommand(
		"certbot",
		"certonly",
		"--standalone",
		"-n",
		"--agree-tos",
		"-m",
		params.EmailAdress,
		"-d",
		hostName,
		"--http-01-port",
		"441",
	)
	out, err = command.Output()
	if err != nil {
		return cert, fmt.Errorf("certbot could not obtain new certificate (output: %s): %+v", string(out), err)
	}

	logs.Info().Str(logs.Component, logs.Certbot).Msgf("successfully registered new certificate for '%s'", hostName)

	// if certbot obtained certificate, find path to cert
	certFile, keyFile, err = FindCertPath(hostName)
	if err != nil {
		return cert, fmt.Errorf("could not find path to obtained certificate: %s", err)
	}

	// if path found, read cert
	cert, err = tlsLoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return cert, fmt.Errorf("could not load x.509 key pair for host name '%s'", hostName)
	}

	return cert, nil
}

var lastRenewal time.Time
var lastRenewalLock sync.Mutex

func Renew() error {
	lastRenewalLock.Lock()
	defer lastRenewalLock.Unlock()
	// if last renewal was only recently, skip renewal
	if lastRenewal.Add(time.Hour * 12).After(time.Now()) {
		return nil
	}

	command := execCommand(
		"certbot",
		"renew",
		"--quiet",
	)
	out, err := command.Output()
	if err != nil {
		return fmt.Errorf("certbot could not renew certificates (output: %s): %+v", string(out), err)
	}

	lastRenewal = time.Now()

	logs.Info().Str(logs.Component, logs.Certbot).Msgf("successfully renewed certificates")

	return nil
}
