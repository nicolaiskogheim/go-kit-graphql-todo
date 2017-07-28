package todo

import "github.com/nicolaiskogheim/go-kit-graphql-todo/user"

var (
	Todo1 = &Todo{ID: NextTodoID(), Text: "Learn some GraphQL", Done: true, OwnerID: user.User1.ID}
	Todo2 = &Todo{ID: NextTodoID(), Text: "Build an app", Done: false, OwnerID: user.User2.ID}
	Todo3 = &Todo{ID: NextTodoID(), Text: "Finish that project", Done: true, OwnerID: user.User2.ID}
	Todo4 = &Todo{ID: NextTodoID(), Text: "Eat breakfast", Done: true, OwnerID: user.User3.ID}
)
