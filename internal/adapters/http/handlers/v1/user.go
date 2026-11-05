package v1

import (
	"context"
	"time"

	"github.com/loopopen/t-ddd-gin-gorm/internal/adapters/http/dto"
	"github.com/loopopen/t-ddd-gin-gorm/internal/adapters/http/handlers"
	"github.com/loopopen/t-ddd-gin-gorm/internal/applic/query"
	"github.com/loopopen/t-ddd-gin-gorm/internal/applic/result"
	"github.com/loopopen/t-ddd-gin-gorm/internal/applic/service"
)

type QueryUsersInput struct {
	Key string `query:"key" required:"true" minLength:"1" doc:"Search key for user name"`
}

type QueryUsersOutput struct {
	Body dto.Response[[]result.User]
}

func QueryUsers(svc *service.UserSvc) handlers.Handler[QueryUsersInput, QueryUsersOutput] {
	return func(ctx context.Context, input *QueryUsersInput) (*QueryUsersOutput, error) {
		users, err := svc.Query(ctx, query.UserQuery{Key: input.Key})
		if err != nil {
			return nil, err
		}
		select {
		case <-time.After(2 * time.Second):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
		return &QueryUsersOutput{Body: dto.OK(users)}, nil
	}
}
