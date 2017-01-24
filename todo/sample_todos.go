package todo

var (
	Todo1 = &Todo{ID: NextTodoID(), Text: "Learn some GraphQL", Done: true}
	Todo2 = &Todo{ID: NextTodoID(), Text: "Build an app", Done: false}
	Todo3 = &Todo{ID: NextTodoID(), Text: "Finish that project", Done: true}
	Todo4 = &Todo{ID: NextTodoID(), Text: "Eat breakfast", Done: true}
)
