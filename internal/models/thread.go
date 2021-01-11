package models

import "time"

type Thread struct {
	ID uint32 `json:"id"`
	Title string `json:"title"`
	Nickname string `json:"author"`
	Forum string `json:"forum"`
	Message string `json:"message"`
	Votes int32 `json:"votes"`
	Slug string `json:"slug"`
	Created time.Time `json:"created"`
}
