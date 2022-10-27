package admin

import (
	"net/http"

	"github.com/frederik-jatzkowski/hermes/params"
)

func authorized(request *http.Request) bool {
	// retrieve and compare credentials
	user, password, ok := request.BasicAuth()

	return ok && user == params.User && password == params.Password
}

func (admin *adminPanel) handleAuth(response http.ResponseWriter, request *http.Request) {
	// prevent brute force attacks
	admin.lock()
	defer admin.unlock()

	switch request.Method {
	case http.MethodGet:
		if authorized(request) {
			response.WriteHeader(http.StatusOK)

			return
		}
		response.WriteHeader(http.StatusMethodNotAllowed)
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
	}
}
