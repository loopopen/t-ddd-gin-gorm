package service

import (
	"context"
	"log/slog"

	"github.com/loopopen/pkg/slogx"
	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/applic/query"
	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/applic/result"
	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/domain/event"
	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/domain/repo"

	"github.com/loopopen/gap"
	"gorm.io/gorm"
)

type UserSvc struct {
	logger   *slog.Logger
	db       *gorm.DB
	pub      gap.EventPublisher
	userRepo repo.UserRepo
}

func NewUserSvc(logger *slog.Logger, db *gorm.DB, pub gap.EventPublisher, quantRepo repo.UserRepo) *UserSvc {
	svc := &UserSvc{
		logger:   slogx.Within(logger, &UserSvc{}),
		db:       db,
		pub:      pub,
		userRepo: quantRepo,
	}
	return svc
}

func (svc *UserSvc) Query(ctx context.Context, q query.UserQuery) ([]result.User, error) {
	users, err := svc.userRepo.Query(ctx, q.Key)
	if err != nil {
		return nil, err
	}
	var rs []result.User
	for _, q := range users {
		rs = append(rs, *new(result.User).FromEntity(q))
	}
	return rs, nil
}

//go:generate go tool gapc -file=$GOFILE

// @subscribe
func (svc *UserSvc) HandlerUserCreated() gap.Handler[*event.UserCreated] {
	return func(ctx context.Context, msg *event.UserCreated, headers map[string]string) error {
		panic("unimplemented")
	}
}
