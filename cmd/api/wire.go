//go:build wireinject

package main

import (
	"context"
	"log/slog"

	"github.com/loopopen/gap/broker/xkafka"
	"github.com/loopopen/t-ddd-gin-huma-gorm/cmd/api/config"
	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/adapters/http"
	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/adapters/outbound"
	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/applic/service"
	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/infra/conf"
	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/infra/gorm"
	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/infra/gorm/repoimpl"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func wireApp(ctx context.Context, e *config.Env, c *config.Config, g config.Gap, k xkafka.Options, orm conf.ORM, log *slog.Logger) (*gin.Engine, error) {
	panic(wire.Build(
		repoimpl.ProviderSet,
		service.ProviderSet,
		outbound.ProviderSet,
		gorm.NewGormDB,
		http.ProviderSet,
	))
}
