package user

var (
	User1 = &User{ID: NextUserID(), Name: "Jon", Email: "jon@jon.com", Password: "123"}
	User2 = &User{ID: NextUserID(), Name: "Jim", Email: "jim@jim.com", Password: "123"}
	User3 = &User{ID: NextUserID(), Name: "Jane", Email: "jane@jane.com", Password: "123"}
	User4 = &User{ID: NextUserID(), Name: "Joane", Email: "joane@joane.com", Password: "123"}
)
