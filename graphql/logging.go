package graphql

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"golang.org/x/net/context"
)

type loggingService struct {
	logger log.Logger
	Service
}

func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) Do(ctx context.Context, request interface{}) (res *graphql.Result) {
	req := request.(*handler.RequestOptions)
	defer func(begin time.Time) {
		var err error
		if len(res.Errors) > 0 {
			err = fmt.Errorf("request error: %v", res.Errors)
		}

		// TODO(nicolai): Can/should we do anything with errors?
		variables, _ := json.Marshal(req.Variables)

		s.logger.Log(
			"method", "do",
			"took", time.Since(begin),
			"error", err,
			"operationName", req.OperationName,
			"variables", variables,
			"query", req.Query,
		)
	}(time.Now())
	res = s.Service.Do(ctx, request)
	return
}
