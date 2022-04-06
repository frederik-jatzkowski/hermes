package certbot

import (
	"crypto/tls"
	"errors"
	"os"
	"os/exec"
	"strings"

	"fleo.software/infrastructure/hermes/logs"
)

var installed bool = false
var mail string = ""

func findCertificate(hostname string) (string, error) {
	// read directory with certbot-certificates
	dir, err := os.ReadDir("/etc/letsencrypt/live/")
	if err != nil {
		return "", err
	}
	// find first matching subdirectory
	for i := range dir {
		if strings.HasPrefix(dir[i].Name(), hostname) && dir[i].IsDir() {
			return "/etc/letsencrypt/live/" + dir[i].Name() + "/", nil
		}
	}
	// if none found, return error
	return "", errors.New("no certificate directory found for hostname '" + hostname + "'")
}

func ObtainCertificate(servername string) (*tls.Certificate, error) {
	// install certbot
	if !installed {
		err := exec.Command("apt-get", "install", "certbot", "-y").Run()
		if err != nil {
			logs.LaunchPrint(err, "0101")
			return nil, errors.New("certbot is not installed or installable on the system")
		}
		logs.LaunchPrint("certbot is installed on the system", "0100")
		installed = true
	}
	// check, if certificate already exists
	path, err := findCertificate(servername)
	if err == nil {
		cert, err := tls.LoadX509KeyPair(path+"fullchain.pem", path+"privkey.pem")
		if err == nil {
			return &cert, err
		}
	}
	// if cert for hostname does not exist, obtain certificate from certbot
	command := exec.Command("certbot", "certonly", "--standalone", "-n", "--agree-tos", "-m", mail, "-d", servername)
	out, err := command.Output()
	if err != nil {
		logs.LaunchPrint(string(out), "")
		return nil, err
	}
	logs.LaunchPrint("successfully obtained new certificate for servername '"+servername+"'", "3300")
	// if certbot is successfull, find path to cert
	path, err = findCertificate(servername)
	if err != nil {
		return nil, err
	}
	// if path found, read cert
	cert, err := tls.LoadX509KeyPair(path+"fullchain.pem", path+"privkey.pem")
	if err != nil {
		return nil, err
	}
	return &cert, nil

}

func SetEmail(email string) {
	mail = email
}
