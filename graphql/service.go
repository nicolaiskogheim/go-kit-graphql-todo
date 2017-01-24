package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"golang.org/x/net/context"
)

type Service interface {
	Do(ctx context.Context, request interface{}) *graphql.Result
}

type service struct {
	schema graphql.Schema
}

func NewService(schema graphql.Schema) Service {
	return &service{
		schema: schema,
	}
}

func (s *service) Do(ctx context.Context, request interface{}) *graphql.Result {
	// TODO(nicolai): type alias RequestOptions to GraphqlRequest?
	options := request.(*handler.RequestOptions)
	params := graphql.Params{
		Context:        ctx,
		OperationName:  options.OperationName,
		RequestString:  options.Query,
		Schema:         s.schema,
		VariableValues: options.Variables,
	}
	return graphql.Do(params)
}
