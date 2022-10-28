package redirect

import (
	"net/http"

	"github.com/frederik-jatzkowski/hermes/logs"
)

type redirect struct{}

func (r *redirect) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	http.Redirect(
		response,
		request,
		"https://"+request.Host+request.RequestURI,
		http.StatusMovedPermanently,
	)
}

var server http.Server

func Start() {
	logs.Info().Str(logs.Component, logs.Redirect).Msg("starting redirect")

	server.Close()

	server = http.Server{
		Handler: &redirect{},
	}

	go server.ListenAndServe()

	logs.Info().Str(logs.Component, logs.Redirect).Msg("successfully started redirect")
}
func Stop() {
	logs.Info().Str(logs.Component, logs.Redirect).Msg("stopping redirect")

	server.Close()

	logs.Info().Str(logs.Component, logs.Redirect).Msg("successfully stopped redirect")
}
