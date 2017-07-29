//go:generate granate
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	kithttp "github.com/go-kit/kit/transport/http"
	stdprometheus "github.com/prometheus/client_golang/prometheus"

	"github.com/nicolaiskogheim/go-kit-graphql-todo/auth"
	"github.com/nicolaiskogheim/go-kit-graphql-todo/graphql"
	"github.com/nicolaiskogheim/go-kit-graphql-todo/inmem"
	"github.com/nicolaiskogheim/go-kit-graphql-todo/models"
	"github.com/nicolaiskogheim/go-kit-graphql-todo/schema"
	"github.com/nicolaiskogheim/go-kit-graphql-todo/session"
	"github.com/nicolaiskogheim/go-kit-graphql-todo/todo"
	"github.com/nicolaiskogheim/go-kit-graphql-todo/user"
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
	)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = &serializedLogger{Logger: logger}
		logger = log.With(logger,
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	var (
		todos    = inmem.NewTodoRepository()
		users    = inmem.NewUserRepository()
		sessions = inmem.NewSessionRepository()
	)

	fieldKeys := []string{"method"}

	var todoService todo.Service
	{
		todoService = todo.NewService(todos)
		todoService = todo.NewLoggingService(logger, todoService)
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

	var userService user.Service
	{
		userService = user.NewService(users)
		userService = user.NewLoggingService(logger, userService)
		userService = user.NewInstrumentingService(
			kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
				Namespace: "api",
				Subsystem: "user_service",
				Name:      "request_count",
				Help:      "Number of requests received.",
			}, fieldKeys),
			kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
				Namespace: "api",
				Subsystem: "user_service",
				Name:      "request_latency_microseconds",
				Help:      "Total duration of requests in microseconds.",
			}, fieldKeys),
			userService,
		)
	}

	var sessionService session.Service
	{
		sessionService = session.NewService(sessions)
	}
	_ = sessionService

	var authService auth.Service
	{
		authService = auth.NewService(sessionService, userService)
	}

	var gqls graphql.Service
	{
		root := models.Root{TodoService: todoService, UserService: userService}
		schema.Init(schema.ProviderConfig{
			Query:    root,
			Mutation: root,
		})

		gqls = graphql.NewService(schema.Schema())
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

	httpLogger := log.With(logger, "component", "http")

	mux := http.NewServeMux()
	mux.Handle("/auth", auth.MakeHandler(authService, httpLogger))
	mux.Handle("/graphql", graphql.MakeHandler(gqls, httpLogger,
		kithttp.ServerBefore(authService.Authenticate)))
	mux.Handle("/", http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Write(graphiql)
	}))

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

type serializedLogger struct {
	mtx sync.Mutex
	log.Logger
}

func (l *serializedLogger) Log(keyvals ...interface{}) error {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	return l.Logger.Log(keyvals...)
}

var graphiql = []byte(`
<!DOCTYPE html>
<html>
   <head>
      <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.2/graphiql.css" />
      <script src="https://cdnjs.cloudflare.com/ajax/libs/fetch/1.1.0/fetch.min.js"></script>
      <script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.6.1/react.min.js"></script>
      <script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.6.1/react-dom.min.js"></script>
      <script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.2/graphiql.js"></script>
   </head>
   <body style="width: 100%; height: 100%; margin: 0; overflow: hidden;">
      <div id="graphiql" style="height: 100vh;">Loading...</div>
      <script>
         function graphQLFetcher(graphQLParams) {
            graphQLParams.variables = graphQLParams.variables ? JSON.parse(graphQLParams.variables) : null;
            return fetch("/graphql", {
               method: "post",
               body: JSON.stringify(graphQLParams),
               credentials: "include",
            }).then(function (response) {
               return response.text();
            }).then(function (responseBody) {
               try {
                  return JSON.parse(responseBody);
               } catch (error) {
                  return responseBody;
               }
            });
         }
         ReactDOM.render(
            React.createElement(GraphiQL, {fetcher: graphQLFetcher}),
            document.getElementById("graphiql")
         );
      </script>
   </body>
</html>
`)
