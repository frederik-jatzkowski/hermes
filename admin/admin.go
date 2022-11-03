package admin

import (
	cryptoRand "crypto/rand"
	"crypto/tls"
	"fmt"
	"math/big"
	mathRand "math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/frederik-jatzkowski/hermes/certbot"
	"github.com/frederik-jatzkowski/hermes/config"
	"github.com/frederik-jatzkowski/hermes/core"
	"github.com/frederik-jatzkowski/hermes/logs"
	"github.com/frederik-jatzkowski/hermes/params"
)

var panel *adminPanel

type adminPanel struct {
	sysCore    *core.Core
	server     *http.Server
	accessLock sync.Mutex
}

func Start() error {
	var (
		err error
	)

	logs.Info().Str(logs.Component, logs.Admin).Msg("starting admin panel")

	// singleton check
	if panel != nil {
		return fmt.Errorf("another admin panel already exists")
	}
	panel = &adminPanel{}

	// check retrieval of certificate for admin panel
	_, err = certbot.ObtainCertificate(params.AdminHost)
	if err != nil {
		return fmt.Errorf("failed to obtain certificate for admin panel: %s", err)
	}

	// start admin server
	go panel.serve()

	// start system core
	err = panel.startSysCore()
	if err != nil {
		logs.Error().Msgf("could not start system core on startup: %s", err)
	}

	logs.Info().Str(logs.Component, logs.Admin).Msg("successfully started admin panel")

	return nil
}

func Stop() {
	logs.Info().Str(logs.Component, logs.Admin).Msg("stopping admin panel")

	panel.accessLock.Lock()
	defer panel.accessLock.Unlock()

	panel.server.Close()

	logs.Info().Str(logs.Component, logs.Admin).Msg("successfully stopped admin panel")

	panel.sysCore.Stop()
}

func (admin *adminPanel) startSysCore() error {
	// get active system config
	configs, err := config.ConfigHistory()
	if err != nil {
		return fmt.Errorf("could not retrieve config history: %s", err)
	}
	if len(configs) == 0 {
		return fmt.Errorf("config history is empty")
	}
	config := configs[len(configs)-1]

	// build system from config
	admin.sysCore = core.NewCore(config)

	// start system from config
	return admin.sysCore.Start()
}

func (admin *adminPanel) serve() {
	// define serveMux
	mux := http.NewServeMux()
	// define server
	panel.server = &http.Server{
		Addr:    ":440",
		Handler: mux,
		TLSConfig: &tls.Config{
			GetCertificate: func(chi *tls.ClientHelloInfo) (*tls.Certificate, error) {
				var (
					err  error
					cert tls.Certificate
				)

				cert, err = certbot.ObtainCertificate(params.AdminHost)
				if err != nil {
					logs.Error().Str(logs.Component, logs.Admin).Msgf("failed to obtain certificate for admin panel: %s", err)

					return nil, err
				}

				return &cert, err
			},
		},
	}

	// define endpoints
	mux.HandleFunc("/config", panel.configEndpoint)
	mux.HandleFunc("/auth", panel.authEndpoint)
	mux.Handle("/", http.FileServer(http.Dir("/opt/hermes/static")))

	// serve admin panel until server is actively closed
	var err error
	for err != http.ErrServerClosed {
		err := admin.server.ListenAndServeTLS("", "")
		if err != http.ErrServerClosed {
			logs.Error().Msgf("admin panel server terminated unexpectedly: %s", err)
		}
	}
}

func (admin *adminPanel) lock() {
	// lock admin panel
	admin.accessLock.Lock()

	// prevent timing attacks
	const MAX_WAIT = 256
	wait, err := cryptoRand.Int(cryptoRand.Reader, big.NewInt(MAX_WAIT))
	if err != nil {
		wait = big.NewInt(int64(mathRand.Intn(MAX_WAIT)))
	}
	time.Sleep(time.Millisecond * time.Duration(wait.Int64()))
}

func (admin *adminPanel) unlock() {
	go func() {
		// prevent brute force attacks
		time.Sleep(time.Second)

		// release admin panel after waiting period
		admin.accessLock.Unlock()
	}()
}
