package models

type Area struct {
	Name string
	Url  string
}

type Pagination struct {
	Next     *string
	Previous *string
}
