package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/loopopen/gap"
	"github.com/loopopen/t-ddd-gin-gorm/cmd/api/config"
	"github.com/loopopen/t-ddd-gin-gorm/internal/adapters/http/dto"
	"github.com/loopopen/t-ddd-gin-gorm/internal/adapters/http/handlers"
	"github.com/loopopen/t-ddd-gin-gorm/internal/applic/service"
	"github.com/loopopen/t-ddd-gin-gorm/pkg/schema/errx"
)

const maxBodyBytes = 4 * 1024 * 1024
const gapShutdownTimeout = 5 * time.Second

var ProviderSet = wire.NewSet(NewApp)

// App groups the HTTP server and its Gap runtime so they can be
// started and stopped together.
type App struct {
	Engine *gin.Engine
	gap    *gap.Gap
	logger *slog.Logger
}

func NewApp(
	e *config.Env,
	c *config.Config,
	userSvc *service.UserSvc,
	g *gap.Gap,
	pub gap.EventPublisher,
	logger *slog.Logger,
) *App {
	if e.IsProd() {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()
	engine.Use(ErrorHandler(logger))
	engine.Use(bodyLimit(maxBodyBytes))
	engine.Use(cors.New(corsConfig(c.CORS)))
	if e.IsProd() && c.Timeout > 0 {
		engine.Use(deadline(time.Duration(c.Timeout) * time.Millisecond))
	}

	registerAPI(engine, e, c, userSvc, logger)

	if g != nil {
		engine.Any("/dashboard/*path", gin.WrapH(g.DashboardHandler()))
		g.Subscribe(gap.From(pub), gap.Inject(userSvc))
	}

	return &App{
		Engine: engine,
		gap:    g,
		logger: logger,
	}
}

func (a *App) ShutdownGap() error {
	if a.gap != nil {
		drainCtx, cancel := context.WithTimeout(context.Background(), gapShutdownTimeout)
		defer cancel()
		if err := a.gap.Shutdown(drainCtx); err != nil {
			return err
		}
	}
	return nil
}

func OpenAPISpec(e *config.Env, c *config.Config, userSvc *service.UserSvc, logger *slog.Logger) *huma.OpenAPI {
	gin.SetMode(gin.TestMode)
	return registerAPI(gin.New(), e, c, userSvc, logger).OpenAPI()
}

func registerAPI(app *gin.Engine, e *config.Env, c *config.Config, userSvc *service.UserSvc, logger *slog.Logger) huma.API {
	hc := huma.DefaultConfig("API Document-"+e.Name, version(e.CommitSHA))
	hc.Info.Contact = &huma.Contact{Name: "Owner", Email: "user@example.com"}
	hc.DocsPath = defaultString(c.OpenAPI.DocsPath, "/docs")
	hc.OpenAPIPath = defaultString(c.OpenAPI.SpecPath, "/openapi")
	if c.OpenAPI.ServerURL != "" {
		hc.Servers = []*huma.Server{{URL: strings.TrimRight(c.OpenAPI.ServerURL, "/")}}
	}

	api := humagin.New(app, hc)
	huma.NewError = newHumaError(logger)

	registerUserAPI(api, userSvc)

	return api
}

func bodyLimit(limit int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limit)
		c.Next()
	}
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

func deadline(t time.Duration) gin.HandlerFunc {
	return timeout.New(
		timeout.WithTimeout(t),
		timeout.WithResponse(func(c *gin.Context) {
			c.JSON(http.StatusGatewayTimeout, dto.Err(errx.ErrServiceTimeout))
		}),
	)
}

func withTimeout[I, O any](t time.Duration, handler handlers.Handler[I, O]) handlers.Handler[I, O] {
	if t <= 0 {
		return handler
	}

	return func(ctx context.Context, input *I) (*O, error) {
		ctx, cancel := context.WithTimeout(ctx, t)
		defer cancel()

		type result struct {
			output *O
			err    error
		}
		results := make(chan result, 1)
		go func() {
			output, err := handler(ctx, input)
			results <- result{output: output, err: err}
		}()

		select {
		case result := <-results:
			return result.output, result.err
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return nil, errx.ErrServiceTimeout
			}
			return nil, ctx.Err()
		}
	}
}
