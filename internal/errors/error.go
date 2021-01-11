package errors

import (
	"errors"
	"net/http"
)

const (
	ErrInternal = iota
	ErrNotFoundUser
	ErrEmailExists
	ErrForumDoesNotExists
	ErrParent
	ConstParentNotFound = "00404"
	ConstUserNotFound = "23503"
)

var Errors = map[int]error {
	ErrInternal:	errors.New("internal server error"),
	ErrParent:	errors.New("parent post was created in another thread"),
	ErrNotFoundUser:	errors.New("can't find user"),
	ErrEmailExists:	errors.New("can't find user"),
	ErrForumDoesNotExists:	errors.New("can't find user"),
}

var StatusCode = map[int]int {
	ErrInternal:	http.StatusInternalServerError,
	ErrParent:	http.StatusConflict,
	ErrNotFoundUser:	http.StatusNotFound,
	ErrEmailExists:	http.StatusConflict,
	ErrForumDoesNotExists: http.StatusNotFound,
}