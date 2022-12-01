package domain

type Post struct {
	Id     int64
	UserId int64
	Title  string
	Body   string
}

type Posts struct {
	Items []Post
	Total uint64
	Pages uint64
}
