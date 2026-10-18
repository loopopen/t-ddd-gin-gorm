package po

import (
	"time"

	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/domain/entity"
	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/domain/enum"

	"gorm.io/gorm"
)

//go:generate go tool shoot map -path=../../../domain/entity -file=$GOFILE

type User struct {
	gorm.Model
	Name     string      `gorm:"not null;comment:姓名"`
	Gender   enum.Gender `gorm:"not null;comment:性别"`
	Birthday time.Time   `gorm:"not null;comment:生日"`
}

func (User) TableComment() string {
	return "用户"
}

func (u *User) readEntity(*entity.User) {
	u.CreatedAt = time.Time{}
	u.UpdatedAt = time.Time{}
	u.DeletedAt = gorm.DeletedAt{}
}
