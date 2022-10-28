package admin

import (
	"net/http"

	"github.com/frederik-jatzkowski/hermes/logs"
	"github.com/frederik-jatzkowski/hermes/params"
)

func authorized(request *http.Request) bool {
	// retrieve and compare credentials
	user, password, ok := request.BasicAuth()

	return ok && user == params.User && password == params.Password
}

func (admin *adminPanel) authEndpoint(response http.ResponseWriter, request *http.Request) {
	// prevent brute force attacks
	admin.lock()
	defer admin.unlock()

	switch request.Method {
	case http.MethodGet:
		if authorized(request) {
			response.WriteHeader(http.StatusOK)
			logs.Info().Str(logs.Component, logs.Admin).Str(logs.ClientAddress, request.RemoteAddr).Msg("someone logged into the admin panel")

			return
		}
		response.WriteHeader(http.StatusUnauthorized)
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
	}
}
