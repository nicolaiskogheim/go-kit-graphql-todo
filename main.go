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

	var gqls graphql.Service
	{
		var todoService todo.Service
		todoService = todo.NewService(todos)
		schema, err := graphql.NewSchema(todoService)
		if err != nil {
			logger.Log("error", err)
			os.Exit(1)
		}
		gqls = graphql.NewService(schema)
		gqls = graphql.NewLoggingService(logger, gqls)
	}

	httpLogger := log.NewContext(logger).With("component", "http")

	mux := http.NewServeMux()
	mux.Handle("/graphql", graphql.MakeHandler(ctx, gqls, httpLogger))

	errc := make(chan error, 2)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// TODO(nicolai): narqo debug listener

	go func() {
		logger.Log("transport", "http", "address", httpAddr, "msg", "listening")
		errc <- http.ListenAndServe(*httpAddr, mux)
	}()

	logger.Log("terminated", <-errc)
}

func envString(varName, fallback string) string {
	value := os.Getenv(varName)
	if value == "" {
		return fallback
	}
	return value
}
