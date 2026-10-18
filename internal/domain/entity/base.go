package entity

import (
	"errors"
)

//go:generate go tool shoot new -getset -sign=@ -file=$GOFILE

var _ Entity[uint] = &Base[uint]{}

type Base[Tid ID] struct {
	//@get
	id Tid
}

func (e *Base[Tid]) Equals(other Entity[Tid]) bool {
	if other == nil {
		return false
	}
	return e.id == other.Id()
}

func (e Base[Tid]) Validate() error {
	return errors.New("unimplemented")
}
