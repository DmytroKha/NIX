package domain

type User struct {
	Id       int64
	Email    string
	Password string
	Name     string
}

type Users struct {
	Items []User
	Total uint64
	Pages uint64
}
