package repository

import "errors"

var (
	ErrorAlreadyExists = errors.New("an entry with this ID already exists")
	ErrorNotFound      = errors.New("an entry with this id was not found")
)
