package repoimpl

import (
	"context"

	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/domain/entity"
	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/domain/repo"
	"github.com/loopopen/t-ddd-gin-huma-gorm/internal/infra"
	"github.com/loopopen/t-ddd-gin-huma-gorm/pkg/schema/errx"
	"gorm.io/gorm"
)

type TE[Tid entity.ID, T any] = entity.EntityConstraint[Tid, T]
type TPO[Te any, Tpo any] = interface {
	infra.ConvertiblePO[Te, Tpo]
}

type Base[Tid entity.ID, T any, Te TE[Tid, T], Tpo TPO[Te, Tpo]] struct {
	db        *gorm.DB
	assembler infra.Assembler[Te, Tpo]
	zeroID    Tid
	zeroPO    Tpo //only used for .Model(zeroPO)
}

func NewBase[Tid entity.ID, T any, Te TE[Tid, T], Tpo TPO[Te, Tpo]](db *gorm.DB) *Base[Tid, T, Te, Tpo] {
	r := &Base[Tid, T, Te, Tpo]{
		db:        db,
		assembler: infra.DefAssembler[Te, Tpo]{},
	}
	var _ repo.Repo[Tid, T, Te] = r
	return r
}

func (r *Base[Tid, T, Te, Tpo]) Exists(ctx context.Context, id Tid) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(r.zeroPO).
		Where("id = ?", id).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 1, nil
}

func (r *Base[Tid, T, Te, Tpo]) FindByID(ctx context.Context, id Tid) (Te, error) {
	var ent Te
	po := r.assembler.ToPO(ent)
	db := r.db.WithContext(ctx).
		Find(&po, "id = ?", id)
	if db.Error != nil {
		return ent, db.Error
	}
	return po.ToEntity(), nil
}

func (r *Base[Tid, T, Te, Tpo]) Create(ctx context.Context, ent Te) (Te, error) {
	if ent == nil {
		return nil, errx.ArgumentIsNil("ent")
	}
	po := r.assembler.ToPO(ent)
	err := r.db.WithContext(ctx).
		Create(po).Error
	if err != nil {
		return nil, err
	}
	ent = r.assembler.ToEntity(po)
	return ent, nil
}

func (r *Base[Tid, T, Te, Tpo]) Save(ctx context.Context, ent Te) error {
	if ent == nil {
		return errx.ArgumentIsNil("ent")
	}
	po := r.assembler.ToPO(ent)
	err := r.db.WithContext(ctx).
		Save(po).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Base[Tid, T, Te, Tpo]) RemoveByID(ctx context.Context, id Tid) error {
	return r.db.WithContext(ctx).
		Model(r.zeroPO).
		Delete("id = ?", id).Error
}
