package pubs

import (
	"context"
	"log/slog"

	"github.com/loopopen/gap"
	"github.com/loopopen/gap/broker/xkafka"
	"github.com/loopopen/gap/dashboard"
	"github.com/loopopen/gap/storage/xgorm"
	tginhumagorm "github.com/loopopen/t-ddd-gin-huma-gorm"
	"github.com/loopopen/t-ddd-gin-huma-gorm/cmd/api/config"

	"gorm.io/gorm"
)

func NewPub(
	ctx context.Context,
	g config.Gap,
	k xkafka.Options,
	db *gorm.DB,
	log *slog.Logger,
) gap.EventPublisher {
	if tginhumagorm.HAVE_NOT_BEEN_DELETED_YET {
		return nil
	}

	pub := gap.NewEventPublisher(
		gap.WithDrain(ctx, 5),
		xgorm.UseGorm(
			xgorm.DB(db),
		),
		xkafka.UseKafka(
			xkafka.Brokers(k.Brokers),
			xkafka.ConfigTopic(
			// xkafka.NumPartitions(4),
			// xkafka.ReplicationFactor(3),
			),
		),
		gap.UseDashboard(
			dashboard.LocationPath(g.Location),
		),
	)
	return pub
}
