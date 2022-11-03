package redirect

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"

	"github.com/frederik-jatzkowski/hermes/logs"
)

var running bool
var runningLock sync.Mutex

type redirect struct{}

func (r *redirect) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// logs.Info().Str(logs.Component, logs.Redirect).Msgf("received request for url '%s' on port 80", request.URL.Path)
	// proxy requests for /.well-known/* to certbot on port 442
	if strings.HasPrefix(request.URL.Path, "/.well-known") || strings.HasPrefix(request.URL.Path, ".well-known") {
		proxyUrl, err := url.Parse("http://localhost:442")
		if err != nil {
			logs.Error().Str(logs.Component, logs.Redirect).Msgf("could not resolve url for proxy: %s", err)

			return
		}

		proxy := httputil.NewSingleHostReverseProxy(proxyUrl)

		proxy.ServeHTTP(response, request)

		return
	}

	// otherwise redirect
	http.Redirect(
		response,
		request,
		"https://"+request.Host+request.RequestURI,
		http.StatusMovedPermanently,
	)
}

var server http.Server

func Start() {
	runningLock.Lock()
	defer runningLock.Unlock()

	if running {
		return
	}

	logs.Info().Str(logs.Component, logs.Redirect).Msg("starting redirect")

	server.Close()

	server = http.Server{
		Addr:    ":80",
		Handler: &redirect{},
	}

	go func() {
		err := server.ListenAndServe()
		if err != http.ErrServerClosed {
			logs.Error().Str(logs.Component, logs.Redirect).Msgf("error while starting redirect: %s", err)
		}
	}()

	running = true
}
func Stop() {
	runningLock.Lock()
	defer runningLock.Unlock()

	if !running {
		return
	}

	logs.Info().Str(logs.Component, logs.Redirect).Msg("stopping redirect")

	server.Close()

	logs.Info().Str(logs.Component, logs.Redirect).Msg("successfully stopped redirect")

	running = false
}
