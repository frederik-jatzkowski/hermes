package core

import (
	"crypto/tls"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/frederik-jatzkowski/hermes/certbot"
	"github.com/frederik-jatzkowski/hermes/config"
	"github.com/frederik-jatzkowski/hermes/logs"
	"github.com/frederik-jatzkowski/hermes/params"
)

type Service struct {
	hostName string
	balancer LoadBalancer
	stopping chan bool
	stopped  chan bool
	started  bool
	access   sync.Mutex
	cert     tls.Certificate
	certLock sync.RWMutex
}

func NewService(config config.Service) *Service {
	return &Service{
		hostName: config.HostName,
		balancer: NewLoadBalancer(config.Balancer),
		stopping: make(chan bool),
		stopped:  make(chan bool),
	}
}

func (service *Service) Start() error {
	var (
		err error
	)

	logs.Info().Str(logs.Component, logs.Gateway).Str(logs.HostName, service.hostName).Msg("starting service")

	service.access.Lock()
	defer service.access.Unlock()

	// obtain initial cert
	cert, err := certbot.ObtainCertificate(service.hostName)
	if err != nil {
		return fmt.Errorf("could not obtain certificate for host name '%s': %s", service.hostName, err)
	}
	service.cert = cert

	// start reload goroutine
	go service.reload()

	// start balancer
	service.balancer.Start()

	service.started = true
	logs.Info().Str(logs.Component, logs.Gateway).Str(logs.HostName, service.hostName).Msg("successfully started service")

	return err
}

func (service *Service) reload() {
	var (
		ticker = time.NewTicker(params.CertCheckInterval)
		cert   tls.Certificate
		err    error
	)

	// renew forever until stopped
	for {
		select {
		case <-ticker.C:
			err = certbot.Renew()
			if err != nil {
				logs.Error().Str(logs.Component, logs.Service).Str(logs.HostName, service.hostName).
					Err(err).Msg("could not renew certificate")
			}
			cert, err = certbot.ObtainCertificate(service.hostName)
			if err != nil {
				logs.Error().Str(logs.Component, logs.Service).Str(logs.HostName, service.hostName).
					Err(err).Msg("could not reload certificate")
			} else {
				service.certLock.Lock()
				service.cert = cert
				service.certLock.Unlock()
			}
		case <-service.stopping:
			ticker.Stop()
			service.stopped <- true

			return
		}
	}
}

func (service *Service) Stop() {

	logs.Info().Str(logs.Component, logs.Gateway).Str(logs.HostName, service.hostName).Msg("stopping service")

	service.access.Lock()
	defer service.access.Unlock()

	if service.started {
		service.stopping <- true
		<-service.stopped
	}

	service.started = false

	service.balancer.Stop()

	logs.Info().Str(logs.Component, logs.Gateway).Str(logs.HostName, service.hostName).Msg("successfully stopped service")
}

func (service *Service) Cert() *tls.Certificate {
	service.certLock.RLock()
	defer service.certLock.RUnlock()

	return &service.cert
}

func (service *Service) Handle(conn *net.Conn) {
	err := service.balancer.Handle(conn)
	if err != nil {
		logs.Error().Str(logs.Component, logs.Service).Str(logs.HostName, service.hostName).
			Err(err).Msg("failed to handle connection")
	}
	logs.Debug().Str(logs.Component, logs.Service).Str(logs.HostName, service.hostName).
		Msgf("handled conn: %s", (*conn).RemoteAddr().String())

}
