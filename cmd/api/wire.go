//go:build wireinject

package main

import (
	"log/slog"

	"github.com/google/wire"
	"github.com/loopopen/gap/broker/xkafka"
	"github.com/loopopen/t-ddd-gin-gorm/cmd/api/config"
	"github.com/loopopen/t-ddd-gin-gorm/internal/adapters/http"
	"github.com/loopopen/t-ddd-gin-gorm/internal/adapters/outbound"
	"github.com/loopopen/t-ddd-gin-gorm/internal/applic/service"
	"github.com/loopopen/t-ddd-gin-gorm/internal/infra/conf"
	"github.com/loopopen/t-ddd-gin-gorm/internal/infra/gorm"
	"github.com/loopopen/t-ddd-gin-gorm/internal/infra/gorm/repoimpl"
)

func wireApp(e *config.Env, c *config.Config, g config.Gap, k xkafka.Options, orm conf.ORM, log *slog.Logger) (*http.App, error) {
	panic(wire.Build(
		repoimpl.ProviderSet,
		service.ProviderSet,
		outbound.ProviderSet,
		gorm.NewGormDB,
		http.ProviderSet,
	))
}
