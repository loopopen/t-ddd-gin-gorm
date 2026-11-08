package http

import (
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	v1 "github.com/loopopen/t-ddd-gin-gorm/internal/adapters/http/handlers/v1"
	"github.com/loopopen/t-ddd-gin-gorm/internal/applic/service"
)

func registerUserAPI(api huma.API, userSvc *service.UserSvc) {
	huma.Register(api, huma.Operation{
		OperationID:  "query-users",
		Method:       http.MethodGet,
		Path:         "/api/v1/users",
		Summary:      "Query users.",
		Description:  "Query users by key.",
		Tags:         []string{"User"},
		Errors:       []int{http.StatusInternalServerError},
		MaxBodyBytes: maxBodyBytes,
	}, withTimeout(0*time.Second, v1.QueryUsers(userSvc)))
}
