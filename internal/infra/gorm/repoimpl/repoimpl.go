package repoimpl

import (
	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/domain/repo"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewUserRepo,
	wire.Bind(new(repo.UserRepo), new(*UserRepo)),
)
