package models

import "errors"

var (
	ErrorUserNotFound      = errors.New("user doesn't exists")
	ErrorUserAlreadyExists = errors.New("user with this id already exists")
)
