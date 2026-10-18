package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-contrib/cors"
	gintimeout "github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/loopopen/gap"
	"github.com/loopopen/t-ddd-gin-huma-gorm/cmd/api/config"
	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/adapters/http/dto"
	v1 "github.com/loopopen/t-ddd-gin-huma-gorm/internal/adapters/http/handlers/v1"
	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/applic/service"
)

const maxBodyBytes = 4 * 1024 * 1024

var ProviderSet = wire.NewSet(NewApp)

func NewApp(
	e *config.Env,
	c *config.Config,
	userSvc *service.UserSvc,
	pub gap.EventPublisher,
	logger *slog.Logger,
) *gin.Engine {
	if e.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.New()
	app.Use(ErrorMiddleware(logger))
	app.Use(bodyLimit(maxBodyBytes))
	app.Use(cors.New(corsConfig(c.CORS)))
	if !e.IsProd() && c.Timeout > 0 {
		app.Use(requestTimeout(time.Duration(c.Timeout) * time.Millisecond))
	}

	registerAPI(app, e, c, userSvc, logger)

	if pub != nil {
		app.Any("/dashboard/*path", gin.WrapH(gap.NewDashboardHandler(pub)))
		gap.Subscribe(gap.From(pub), gap.Inject(userSvc))
	}

	return app
}

func OpenAPISpec(e *config.Env, c *config.Config, userSvc *service.UserSvc, logger *slog.Logger) *huma.OpenAPI {
	gin.SetMode(gin.TestMode)
	return registerAPI(gin.New(), e, c, userSvc, logger).OpenAPI()
}

func registerAPI(app *gin.Engine, e *config.Env, c *config.Config, userSvc *service.UserSvc, logger *slog.Logger) huma.API {
	docConfig := huma.DefaultConfig("API Document-"+e.Name, version(e.CommitSHA))
	docConfig.Info.Contact = &huma.Contact{Name: "Owner", Email: "user@example.com"}
	docConfig.DocsPath = defaultString(c.OpenAPI.DocsPath, "/docs")
	docConfig.OpenAPIPath = defaultString(c.OpenAPI.SpecPath, "/openapi")
	if c.OpenAPI.ServerURL != "" {
		docConfig.Servers = []*huma.Server{{URL: strings.TrimRight(c.OpenAPI.ServerURL, "/")}}
	}

	api := humagin.New(app, docConfig)
	huma.NewError = NewHumaError
	handleError := HandlerError(logger)

	huma.Register(api, huma.Operation{
		OperationID:  "query-users",
		Method:       http.MethodGet,
		Path:         "/api/v1/users",
		Summary:      "Query users.",
		Description:  "Query users by key.",
		Tags:         []string{"User"},
		Errors:       []int{http.StatusInternalServerError},
		MaxBodyBytes: maxBodyBytes,
	}, v1.QueryUsers(userSvc, handleError))

	huma.Register(api, huma.Operation{
		OperationID:  "query-users-timeout-demo",
		Method:       http.MethodGet,
		Path:         "/api/v1/users/unsafe-timeout",
		Summary:      "Query users. (Timeout Demo)",
		Description:  "Query users by key and demonstrate a one-second request timeout.",
		Tags:         []string{"User"},
		Errors:       []int{http.StatusRequestTimeout, http.StatusInternalServerError},
		MaxBodyBytes: maxBodyBytes,
	}, v1.QueryUsersWithTimeout(userSvc, handleError, time.Second))

	return api
}

func bodyLimit(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
		c.Next()
	}
}

func requestTimeout(limit time.Duration) gin.HandlerFunc {
	return gintimeout.New(
		gintimeout.WithTimeout(limit),
		gintimeout.WithResponse(func(c *gin.Context) {
			c.JSON(http.StatusRequestTimeout, dto.NewError(
				http.StatusRequestTimeout,
				"FRAMEWORK_ERROR",
				"request timed out",
				"",
				nil,
			))
		}),
	)
}

func corsConfig(c config.CORS) cors.Config {
	conf := cors.DefaultConfig()
	conf.AllowOrigins = []string{"*"}
	if len(c.AllowOrigins) > 0 {
		conf.AllowOrigins = c.AllowOrigins
	}
	if len(c.AllowHeaders) > 0 {
		conf.AllowHeaders = c.AllowHeaders
	}
	return conf
}

func version(commitSHA string) string {
	if commitSHA == "" {
		return "1.0"
	}
	return fmt.Sprintf("1.0-%s", commitSHA)
}

func defaultString(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
