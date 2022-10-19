package admin

import "net/http"

func (admin *adminPanel) handleConfig(response http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		admin.handleConfigGet(response, request)
	case http.MethodPost:
		admin.handleConfigPost(response, request)
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (admin *adminPanel) handleConfigGet(response http.ResponseWriter, request *http.Request) {
	if !authorized(request) {
		response.WriteHeader(http.StatusMethodNotAllowed)

		return
	}
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("Hallo Get"))
}

func (admin *adminPanel) handleConfigPost(response http.ResponseWriter, request *http.Request) {
	if !authorized(request) {
		response.WriteHeader(http.StatusMethodNotAllowed)

		return
	}
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("Hallo Post"))
}
