package repo

import (
	"context"

	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/domain/entity"
)

type Repo[Tid entity.ID, T any, Te entity.EntityConstraint[Tid, T]] interface {
	Exists(ctx context.Context, id Tid) (bool, error)

	FindByID(ctx context.Context, id Tid) (Te, error)

	Create(ctx context.Context, e Te) (Te, error)

	Save(ctx context.Context, e Te) error

	RemoveByID(ctx context.Context, id Tid) error
}
