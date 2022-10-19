package admin

import (
	"net/http"

	"github.com/frederik-jatzkowski/hermes/params"
)

func authorized(request *http.Request) bool {
	user, password, ok := request.BasicAuth()

	return ok && user == params.User && password == params.Password
}
