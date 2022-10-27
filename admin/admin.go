package admin

import (
	cryptoRand "crypto/rand"
	"fmt"
	"log"
	"math/big"
	mathRand "math/rand"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/frederik-jatzkowski/hermes/config"
	"github.com/frederik-jatzkowski/hermes/core"
	"github.com/frederik-jatzkowski/hermes/params"
)

var panel *adminPanel

type adminPanel struct {
	tlsCore    *core.Core
	accessLock sync.Mutex
}

func Start() {
	var (
		adminConfig config.Config
		err         error
	)

	// singleton
	if panel != nil {
		return
	}

	adminConfig = config.Config{
		Gateways: []config.Gateway{
			{
				ResolvedAddress: net.TCPAddr{Port: 440},
				Services: []config.Service{
					{
						HostName: params.AdminHost,
						Balancer: config.LoadBalancer{
							Algorithm: "RoundRobin",
							Servers: []config.Server{
								{
									ResolvedAddress: net.TCPAddr{
										Port: 441,
										IP:   net.IPv4(127, 0, 0, 1),
									},
								},
							},
						},
					},
					{
						HostName: "localhost",
						Balancer: config.LoadBalancer{
							Algorithm: "RoundRobin",
							Servers: []config.Server{
								{
									ResolvedAddress: net.TCPAddr{
										Port: 441,
										IP:   net.IPv4(127, 0, 0, 1),
									},
								},
							},
						},
					},
				},
			},
		},
	}

	panel = &adminPanel{
		tlsCore: core.NewCore(adminConfig),
	}

	// start admin core
	err = panel.tlsCore.Start()
	if err != nil {
		log.Fatal(fmt.Errorf("could not start admin panel: %s", err))
	}

	// define endpoints
	http.HandleFunc("/config", panel.handleConfig)
	http.HandleFunc("/auth", panel.handleAuth)
	http.Handle("/", http.FileServer(http.Dir("/opt/hermes/static")))

	server := &http.Server{
		Addr: ":441",
	}

	log.Fatal(server.ListenAndServe())
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
