package pubs

import (
	"github.com/loopopen/gap"
	"github.com/loopopen/gap/broker/xkafka"
	"github.com/loopopen/gap/dashboard"
	"github.com/loopopen/gap/storage/xgorm"
	tginhumagorm "github.com/loopopen/t-ddd-gin-gorm"
	"github.com/loopopen/t-ddd-gin-gorm/cmd/api/config"

	"gorm.io/gorm"
)

func NewGap(g config.Gap, db *gorm.DB) *gap.Gap {
	if tginhumagorm.HAVE_NOT_BEEN_DELETED_YET {
		return nil
	}

	return gap.New(
		xgorm.UseGorm(
			xgorm.DB(db),
		),
		gap.UseDashboard(
			dashboard.LocationPath(g.Location),
		),
	)
}

func NewPub(g *gap.Gap, k xkafka.Options) gap.EventPublisher {
	if g == nil {
		return nil
	}

	return g.NewEventPublisher(
		xkafka.UseKafka(
			xkafka.Brokers(k.Brokers),
			xkafka.UserName(k.UserName),
			xkafka.Password(k.Password),
			xkafka.ConfigTopic(
			// xkafka.NumPartitions(4),
			// xkafka.ReplicationFactor(3),
			),
		),
	)
}
