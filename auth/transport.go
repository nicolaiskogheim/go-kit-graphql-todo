package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"

	"github.com/go-kit/kit/endpoint"
)

func MakeHandler(s Service) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	var authEndpoint endpoint.Endpoint
	{
		authEndpoint = makeAuthEndpoint(s)
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

		cookie := http.Cookie{
			Name:     "session",
			Value:    token,
			HttpOnly: true,
		}

		return cookie, err
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
