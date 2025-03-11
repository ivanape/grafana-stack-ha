package main

import (
	"broker/obs"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_, span := globalTracer.Start(r.Context(), "broker")
	span.SetAttributes(
		attribute.KeyValue{
			Key:   attribute.Key("app"),
			Value: attribute.StringValue("example"),
		},
	)
	defer span.End()

	obs.LogInfoWithSpan(logger, span, r.Context(), "Hit the broker")

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		logger.Error("Error reading JSON: ", err)
		app.errorJSON(w, err)
		return
	}

	context, span := globalTracer.Start(r.Context(), "handle-submision")
	span.SetAttributes(
		attribute.KeyValue{
			Key:   attribute.Key("app"),
			Value: attribute.StringValue("example"),
		},
	)
	defer span.End()

	obs.LogInfoWithSpan(logger, span, context, "Received request: ", requestPayload)

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, r, requestPayload.Auth, context)
	default:
		app.errorJSON(w, errors.New("unknown action"))

	}
}

func (app *Config) authenticate(w http.ResponseWriter, r *http.Request, a AuthPayload, context context.Context) {
	_, authSpan := globalTracer.Start(context, "authenticate")
	authSpan.SetAttributes(
		attribute.KeyValue{
			Key:   attribute.Key("app"),
			Value: attribute.StringValue("example"),
		},
	)
	authSpan.AddEvent("Calling to /authenticate sevice")
	defer authSpan.End()

	jsonData, _ := json.MarshalIndent(a, "", "\t")
	request, err := http.NewRequestWithContext(context, "POST", "http://auth:80/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	response, err := client.Do(request)
	if err != nil {
		obs.LogErrorWithSpan(logger, authSpan, r.Context(), "Error calling auth service: ", err)
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// make sure we get back the correct status code

	if response.StatusCode == http.StatusUnauthorized {
		obs.LogErrorWithSpan(logger, authSpan, r.Context(), "Unauthorized response from auth service: ", response.StatusCode)
		app.errorJSON(w, errors.New(strconv.Itoa(response.StatusCode)))
		return
	} else if response.StatusCode != http.StatusAccepted {
		obs.LogErrorWithSpan(logger, authSpan, r.Context(), "Unexpected response from auth service: ", response.StatusCode)
		app.errorJSON(w, errors.New(strconv.Itoa(response.StatusCode)))
		return
	}

	// Create a variable we'll read response.Body into

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		obs.LogErrorWithSpan(logger, authSpan, r.Context(), "Error decoding response from auth service: ", err)
		app.errorJSON(w, err)
		return
	}

	if jsonFromService.Error {
		obs.LogErrorWithSpan(logger, authSpan, r.Context(), "Error from auth service: ", jsonFromService.Message)
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	obs.LogInfoWithSpan(logger, authSpan, r.Context(), "User authenticated: ", jsonFromService.Data)
	app.writeJSON(w, http.StatusAccepted, payload)

}
