package models

import "time"

type Post struct {
	ID uint32 `json:"id"`
	Parent uint32 `json:"parent"`
	Author string `json:"author"`
	Message string `json:"message"`
	IsEdited bool `json:"isEdited"`
	Forum string `json:"forum"`
	Thread uint32 `json:"thread"`
	Created time.Time `json:"created"`
}

type PostFull struct {
	Post *Post `json:"post"`
	User *User `json:"author,omitempty"`
	Forum *Forum `json:"forum,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
}
