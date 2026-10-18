package repo

import (
	"context"

	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/domain/entity"
)

type UserRepo interface {
	Repo[uint, entity.User, *entity.User]
	Bind(txer Txer) (UserRepo, error)

	//todo:
	Query(ctx context.Context, key string) ([]*entity.User, error)
}
