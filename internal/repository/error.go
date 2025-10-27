package repository

import "errors"

var (
	ErrorAlreadyExists = errors.New("An entry with this ID already exists.")
	ErrorNotFound      = errors.New("An entry with this id was not found.")
)
