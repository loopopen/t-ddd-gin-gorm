package v1

import (
	"context"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/adapters/http/dto"
	timeoutx "github.com/loopopen/t-ddd-gin-huma-gorm/internal/adapters/http/timeout"
	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/applic/query"
	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/applic/result"
	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/applic/service"
)

type QueryUsersInput struct {
	Key string `query:"key" required:"true" minLength:"1" doc:"Search key for user name"`
}

type QueryUsersOutput struct {
	Body dto.Response[[]result.User]
}

func QueryUsers(svc *service.UserSvc, handleError func(error) huma.StatusError) func(context.Context, *QueryUsersInput) (*QueryUsersOutput, error) {
	return func(ctx context.Context, input *QueryUsersInput) (*QueryUsersOutput, error) {
		users, err := svc.Query(ctx, query.UserQuery{Key: input.Key})
		if err != nil {
			return nil, handleError(err)
		}
		return &QueryUsersOutput{Body: dto.OK(users)}, nil
	}
}

func QueryUsersWithTimeout(svc *service.UserSvc, handleError func(error) huma.StatusError, limit time.Duration) func(context.Context, *QueryUsersInput) (*QueryUsersOutput, error) {
	return func(ctx context.Context, input *QueryUsersInput) (*QueryUsersOutput, error) {
		users, err := timeoutx.Run(ctx, limit, func(ctx context.Context) ([]result.User, error) {
			users, err := svc.Query(ctx, query.UserQuery{Key: input.Key})
			if err != nil {
				return nil, err
			}

			select {
			case <-time.After(10 * time.Second):
				return users, nil
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		})
		if err != nil {
			if timeoutx.IsTimeout(err) {
				return nil, huma.Error408RequestTimeout("request timed out")
			}
			return nil, handleError(err)
		}
		return &QueryUsersOutput{Body: dto.OK(users)}, nil
	}
}
