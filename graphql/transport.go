package graphql

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/graphql-go/handler"
	// kittracing "github.com/go-kit/kit/tracing/opentracing"
	// "github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
)

func MakeHandler(ctx context.Context, gqs Service, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	var graphqlEndpoint endpoint.Endpoint
	{
		graphqlLogger := log.NewContext(logger).With("method", "graphql")

		graphqlEndpoint = makeGraphqlEndpoint(gqs)
		// graphqlEndpoint = kittracing.TracerServer(tracer, "graphql")(graphqlEndpoint)
		graphqlEndpoint = makeLoggingGraphqlEndpoint(graphqlLogger)(graphqlEndpoint)
	}

	graphqlHandler := kithttp.NewServer(
		ctx,
		graphqlEndpoint,
		decodeGraphqlRequest,
		encodeResponse,
		opts...,
	)

	return graphqlHandler
}

// TODO(nicolai): alias graphqlRequest = handler.RequestOptions ?
//				  Maybe we need Go 1.9 for this?
// type graphqlRequest struct {
// 	query string
// }

func makeGraphqlEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		// TODO(nicolai): Is this the right place to do this?
		req := request.(*handler.RequestOptions)
		res := s.Do(ctx, req)
		return res, nil
	}
}

// TODO(nicolai): put in logging.go ? Or is this endpoint stuff?
// TODO(nicolai): Log more?
func makeLoggingGraphqlEndpoint(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				logger.Log("error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)
		}
	}
}

// TODO(nicolai): This is unused.
var errBadRequest = errors.New("bad request")

// TODO(nicolai): Don't return 500 Internal Server Error on bad request
// TODO(nicolai): Figure out how narqo gets "Decode: bad request"
func decodeGraphqlRequest(_ context.Context, req *http.Request) (interface{}, error) {
	return handler.NewRequestOptions(req), nil
}

func encodeResponse(ctx context.Context, rw http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), rw)
		return nil
	}
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(rw).Encode(response)
}

type errorer interface {
	error() error
}

func encodeError(_ context.Context, err error, rw http.ResponseWriter) {
	// TODO(nicolai): We should be able to be more granular here.
	// Everything isn't an internal server error.
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Header().Set("Content-Type", "application/json; charset= utf-8")
	json.NewEncoder(rw).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
