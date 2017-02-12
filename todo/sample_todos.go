package todo

var (
	Todo1 = &Todo{ID: NextTodoID(), Text: "Learn some GraphQL", Done: true, OwnerID: "2C2E7C8D"}
	Todo2 = &Todo{ID: NextTodoID(), Text: "Build an app", Done: false, OwnerID: "AC5CF9CF"}
	Todo3 = &Todo{ID: NextTodoID(), Text: "Finish that project", Done: true, OwnerID: "AC5CF9CF"}
	Todo4 = &Todo{ID: NextTodoID(), Text: "Eat breakfast", Done: true, OwnerID: "CCAE5161"}
)
