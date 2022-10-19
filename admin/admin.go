package admin

import (
	"fmt"
	"log"
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
	http.HandleFunc("/", panel.handleIndex)

	server := &http.Server{
		Addr: ":441",
	}

	log.Fatal(server.ListenAndServe())
}

func (admin *adminPanel) lock() {
	admin.accessLock.Lock()
}

func (admin *adminPanel) unlock() {
	time.Sleep(time.Second)
	admin.accessLock.Unlock()
}
