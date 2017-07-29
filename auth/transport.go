package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

func MakeHandler(s Service, logger log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	var authEndpoint endpoint.Endpoint
	{
		authLogger := log.With(logger, "method", "auth")

		authEndpoint = makeAuthEndpoint(s)
		authEndpoint = makeLoggingAuthEndpoint(authLogger)(authEndpoint)
	}

	authHandler := kithttp.NewServer(
		authEndpoint,
		decodeAuthRequest,
		encodeAuthResponse,
		opts...,
	)

	return authHandler
}

func makeAuthEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*http.Request)
		token, err := s.Login(req)

		if err != nil {
			return nil, err
		}

		cookie := http.Cookie{
			Name:     "session",
			Value:    token.ToString(),
			HttpOnly: true,
		}

		return cookie, nil
	}
}

func makeLoggingAuthEndpoint(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				logger.Log("error", err, "took", time.Since(begin))
			}(time.Now())
			return next(ctx, request)
		}
	}
}

func decodeAuthRequest(_ context.Context, req *http.Request) (interface{}, error) {
	return req, nil
}

func encodeAuthResponse(ctx context.Context, rw http.ResponseWriter, response interface{}) error {

	cookie, ok := response.(http.Cookie)
	if ok == false {
		return errors.New("Bad response")
	}

	http.SetCookie(rw, &cookie)
	rw.WriteHeader(http.StatusOK)

	return nil

}

func encodeError(_ context.Context, err error, rw http.ResponseWriter) {
	rw.WriteHeader(http.StatusUnauthorized)
	rw.Header().Set("Content-Type", "application/json; charset= utf-8")

	json.NewEncoder(rw).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
