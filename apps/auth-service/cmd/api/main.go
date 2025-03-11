package main

import (
	"authentication/models"
	"authentication/obs"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/yukitsune/lokirus"
	"go.opentelemetry.io/otel/trace"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var counts int64

var logger *logrus.Logger

var _ trace.Tracer

type Config struct {
	DB     *sql.DB
	Models models.Models
}

func main() {
	var err error
	_, err = obs.NewTracer()
	if err != nil {
		logger.Panic(err)
	}
	metricsConfg, err := obs.NewMetricConfig()
	if err != nil {
		logger.Panic(err)
	}

	profiler, err := obs.NewProfiler()
	if err != nil {
		logger.Panic(err)
	}

	defer profiler.Stop()

	logger = logrus.New()
	// Configure the Loki hook
	opts := lokirus.NewLokiHookOptions().
		// Grafana doesn't have a "panic" level, but it does have a "critical" level
		// https://grafana.com/docs/grafana/latest/explore/logs-integration/
		WithLevelMap(lokirus.LevelMap{logrus.PanicLevel: "critical"}).
		WithFormatter(&logrus.JSONFormatter{}).
		WithStaticLabels(obs.DefaultServiceTags)

	lokiWebHookUrl := os.Getenv("LOKI_WEBHOOK_URL")

	hook := lokirus.NewLokiHookWithOpts(
		lokiWebHookUrl,
		opts,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.FatalLevel)

	logger.Hooks.Add(hook)

	logger.Println("Starting authentication service")

	conn := connectToDB()
	if conn == nil {
		logger.Panic("Can't connect to Postgres!")
	}

	// set up config
	app := Config{
		DB:     conn,
		Models: models.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT")),
		Handler: app.routes(metricsConfg),
	}

	err1 := srv.ListenAndServe()

	if err1 != nil {
		logger.Panic(err1)
	}

}

// openDB is responsible for creating connection to database
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// connectToDB is using openDB function to create new connection to database and handle errors
func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			logger.Println("Postgres not yet ready ...")
			counts++
		} else {
			logger.Println("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			logger.Println(err)
			return nil
		}

		logger.Println("Backing off for two seconds.....")
		time.Sleep(2 * time.Second)
		continue
	}
}
