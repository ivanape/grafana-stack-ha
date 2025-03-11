package main

import (
	"broker/obs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/riandyrn/otelchi"
	otelchimetric "github.com/riandyrn/otelchi/metric"
)

func (app *Config) routes(metricConfig otelchimetric.BaseConfig) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(
		otelchi.Middleware(obs.DefaultServiceTags["service"], otelchi.WithChiRoutes(r)),
		otelchimetric.NewRequestDurationMillis(metricConfig),
		otelchimetric.NewRequestInFlight(metricConfig),
		otelchimetric.NewResponseSizeBytes(metricConfig),
		middleware.Heartbeat("/ping"),
	)

	r.Post("/", app.Broker)
	r.Post("/handle", app.HandleSubmission)

	return r
}
