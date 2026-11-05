package result

import (
	"github.com/loopopen/pkg/timex"
	"github.com/loopopen/t-ddd-gin-gorm/internal/domain/entity"
	"github.com/loopopen/t-ddd-gin-gorm/internal/infra"
)

//go:generate go tool shoot new -json -file=$GOFILE
//go:generate go tool shoot map -path=../../domain/entity -way=<-  -file=$GOFILE -i

type User struct {
	infra.Mapper `json:"-"`
	ID           uint       `json:"id"`
	Name         string     `json:"name"`
	Birthday     timex.Date `json:"birthday" format:"date" example:"1987-11-29"`
	Age          int        `json:"age"`
}

func (q *User) readEntity(e *entity.User) {
	q.Age = e.Age()
}
