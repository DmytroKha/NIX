package domain

type Comment struct {
	Id     int64
	PostId int64
	Name   string
	Email  string
	Body   string
}

type Comments struct {
	Items []Comment
	Total uint64
	Pages uint64
}
