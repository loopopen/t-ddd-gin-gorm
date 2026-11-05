package outbound

import (
	"github.com/google/wire"
	"github.com/loopopen/t-ddd-gin-gorm/internal/adapters/outbound/pubs"
)

var ProviderSet = wire.NewSet(pubs.NewGap, pubs.NewPub)
