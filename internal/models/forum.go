package models

type Forum struct {
	Slug string `json:"slug"`
	Title string `json:"title"`
	User string `json:"user"`
	Posts uint32 `json:"posts"`
	Threads uint32 `json:"threads"`
}
