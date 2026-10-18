package entity

import (
	"time"

	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/domain/enum"
)

//go:generate go tool shoot new -getset -file=$GOFILE

type User struct {
	Base[uint]
	name     string
	birthday time.Time
	gender   enum.Gender
}

func (u *User) Age() int {
	return time.Now().Year() - u.birthday.Year()
}
