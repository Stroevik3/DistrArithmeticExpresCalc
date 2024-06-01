package store

import "errors"

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrRecordDubleId  = errors.New("an entry with the same id already exists")
	ErrStackIsEmpty   = errors.New("stack is empty")
)
