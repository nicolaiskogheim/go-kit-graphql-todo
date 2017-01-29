package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"

	"github.com/nicolaiskogheim/go-kit-graphql-todo/graphql"
	"github.com/nicolaiskogheim/go-kit-graphql-todo/inmem"
	"github.com/nicolaiskogheim/go-kit-graphql-todo/todo"
)

const (
	defaultPort      = "8080"
	defaultDebugPort = "1337"
)

func main() {
	var (
		port = envString("PORT", defaultPort)
		// TODO(nicolai): This will be used with a /metrics endpoint
		// debugPort = envString("DEBUG_PORT", defaultDebugPort)
		httpAddr = flag.String("http.addr", ":"+port, "HTTP listen address")

		ctx = context.Background()
	)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewContext(logger).With("ts", log.DefaultTimestampUTC)
		logger = log.NewContext(logger).With("caller", log.DefaultCaller)
	}

	var (
		todos = inmem.NewTodoRepository()
	)

	fieldKeys := []string{"method"}

	var todoService todo.Service
	{
		todoService = todo.NewService(todos)
		// TODO: logging
		todoService = todo.NewInstrumentingService(
			kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
				Namespace: "api",
				Subsystem: "todo_service",
				Name:      "request_count",
				Help:      "Number of requests received.",
			}, fieldKeys),
			kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
				Namespace: "api",
				Subsystem: "todo_service",
				Name:      "request_latency_microseconds",
				Help:      "Total duration of requests in microseconds.",
			}, fieldKeys),
			todoService,
		)
	}

	var gqls graphql.Service
	{
		schema, err := graphql.NewSchema(todoService)
		if err != nil {
			logger.Log("error", err)
			os.Exit(1)
		}
		gqls = graphql.NewService(schema)
		gqls = graphql.NewLoggingService(logger, gqls)
		gqls = graphql.NewInstrumentingService(
			kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
				Namespace: "api",
				Subsystem: "graphql_service",
				Name:      "request_count",
				Help:      "Number of requests received.",
			}, fieldKeys),
			kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
				Namespace: "api",
				Subsystem: "graphql_service",
				Name:      "request_latency_microseconds",
				Help:      "Total duration of requests in microseconds.",
			}, fieldKeys),
			gqls,
		)
	}

	httpLogger := log.NewContext(logger).With("component", "http")

	mux := http.NewServeMux()
	mux.Handle("/graphql", graphql.MakeHandler(ctx, gqls, httpLogger))

	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", stdprometheus.Handler())

	errc := make(chan error, 2)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// TODO(nicolai): narqo debug listener

	go func() {
		logger.Log("transport", "http", "address", httpAddr, "msg", "listening")
		errc <- http.ListenAndServe(*httpAddr, nil)
	}()

	logger.Log("terminated", <-errc)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

func envString(varName, fallback string) string {
	value := os.Getenv(varName)
	if value == "" {
		return fallback
	}
	return value
}
