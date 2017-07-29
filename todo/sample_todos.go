package todo

import "github.com/nicolaiskogheim/go-kit-graphql-todo/user"

var (
	Todo1 = &Todo{ID: "7A421DFE", Text: "Learn some GraphQL", Done: true, OwnerID: user.User1.ID}
	Todo2 = &Todo{ID: "3C2120A0", Text: "Build an app", Done: false, OwnerID: user.User2.ID}
	Todo3 = &Todo{ID: "AF1BD873", Text: "Finish that project", Done: true, OwnerID: user.User2.ID}
	Todo4 = &Todo{ID: "BAA13A54", Text: "Eat breakfast", Done: true, OwnerID: user.User3.ID}
)
