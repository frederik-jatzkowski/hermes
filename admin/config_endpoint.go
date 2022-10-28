package admin

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/frederik-jatzkowski/hermes/config"
	"github.com/frederik-jatzkowski/hermes/core"
	"github.com/frederik-jatzkowski/hermes/logs"
)

func (admin *adminPanel) configEndpoint(response http.ResponseWriter, request *http.Request) {
	// prevent brute force attacks
	admin.lock()
	defer admin.unlock()

	// check authorization
	if !authorized(request) {
		response.WriteHeader(http.StatusUnauthorized)

		return
	}

	// route by method
	switch request.Method {
	case http.MethodGet:
		admin.configEndpointGet(response, request)
	case http.MethodPost:
		admin.configEndpointPost(response, request)
	default:
		response.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (admin *adminPanel) configEndpointGet(response http.ResponseWriter, request *http.Request) {
	// read config history
	history, err := config.ConfigHistory()
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)

		return
	}

	// parse config history
	historyData, err := json.Marshal(history)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)

		return
	}
	logs.Info().Str(logs.Component, logs.Admin).Str(logs.ClientAddress, request.RemoteAddr).Msg("served configuration history")

	response.WriteHeader(http.StatusOK)
	response.Write(historyData)
}

type configPostResponseBody struct {
	Ok         bool     `json:"ok"`
	Exceptions []string `json:"exceptions"`
}

func (body configPostResponseBody) Marshal() []byte {
	data, _ := json.Marshal(body)
	return data
}

func (admin *adminPanel) configEndpointPost(response http.ResponseWriter, request *http.Request) {
	responseBody := configPostResponseBody{
		Ok:         true,
		Exceptions: []string{},
	}

	// read body
	body, err := io.ReadAll(request.Body)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		responseBody.Ok = false
		responseBody.Exceptions = append(responseBody.Exceptions, err.Error())
		response.Write(responseBody.Marshal())

		return
	}

	// parse and validate conf
	conf, err := config.NewConfig(body)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		responseBody.Ok = false
		responseBody.Exceptions = append(responseBody.Exceptions, err.Error())
		response.Write(responseBody.Marshal())

		return
	}

	logs.Info().Str(logs.Component, logs.Admin).Str(logs.ClientAddress, request.RemoteAddr).Msg("starting to apply new configuration")

	// stop running core
	admin.sysCore.Stop()

	// build new core
	newCore := core.NewCore(conf)

	// try to start new core for config
	err = newCore.Start()
	if err != nil {
		// prepare response body
		err = fmt.Errorf("rolling back, could not apply new config: %s", err)
		responseBody.Ok = false
		responseBody.Exceptions = append(responseBody.Exceptions, err.Error())
		logs.Error().Str(logs.Component, logs.Admin).Str(logs.ClientAddress, request.RemoteAddr).Msg(err.Error())

		// try to rollback to old core
		err = admin.sysCore.Start()
		if err != nil {
			err = fmt.Errorf("could not rollback to old configuration: %s", err)
			responseBody.Exceptions = append(responseBody.Exceptions, err.Error())
			logs.Error().Str(logs.Component, logs.Admin).Str(logs.ClientAddress, request.RemoteAddr).Msg(err.Error())
		}

		response.WriteHeader(http.StatusOK)
		response.Write(responseBody.Marshal())

		return
	}

	logs.Info().Str(logs.Component, logs.Admin).Str(logs.ClientAddress, request.RemoteAddr).Msg("successfully started new core")

	// try to persist valid config
	err = config.AppendConfig(conf)
	if err != nil {
		// prepare response body
		err = fmt.Errorf("rolling back, could not persist new config: %s", err)
		responseBody.Ok = false
		responseBody.Exceptions = append(responseBody.Exceptions, err.Error())
		logs.Error().Str(logs.Component, logs.Admin).Str(logs.ClientAddress, request.RemoteAddr).Msg(err.Error())

		// try to rollback to old core
		err = admin.sysCore.Start()
		if err != nil {
			err = fmt.Errorf("could not rollback to old configuration: %s", err)
			responseBody.Exceptions = append(responseBody.Exceptions, err.Error())
			logs.Error().Str(logs.Component, logs.Admin).Str(logs.ClientAddress, request.RemoteAddr).Msg(err.Error())
		}

		response.WriteHeader(http.StatusOK)
		response.Write(responseBody.Marshal())

		return
	}

	logs.Info().Str(logs.Component, logs.Admin).Str(logs.ClientAddress, request.RemoteAddr).Msg("successfully applied new configuration")

	// successfully applied new config
	admin.sysCore = newCore
	response.WriteHeader(http.StatusOK)
	response.Write(responseBody.Marshal())
}
