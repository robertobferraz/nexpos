package src

import (
	"context"
	"errors"

	"github.com/robertobff/nexpos/adapter/outbound/database/postgres"
	"github.com/robertobff/nexpos/domain/dto"
	"github.com/robertobff/nexpos/domain/entity"
	"github.com/robertobff/nexpos/domain/repository"
	"gorm.io/gorm"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var DiscountModule = fx.Module(
	"discount_repository",
	fx.Provide(NewDiscountRepositorySrc),
	fx.Provide(func(p *DiscountRepositorySrc) repository.DiscountRepository { return p }),
)

type DiscountRepositorySrc struct {
	pg     *postgres.Postgres
	logger *zap.SugaredLogger
}

func NewDiscountRepositorySrc(
	pg *postgres.Postgres,
	logger *zap.SugaredLogger,
) (*DiscountRepositorySrc, error) {
	return &DiscountRepositorySrc{
		pg:     pg,
		logger: logger,
	}, nil
}

func (r *DiscountRepositorySrc) Create(ctx context.Context, discount *entity.Discount) error {
	err := r.pg.Db.WithContext(ctx).Create(discount).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *DiscountRepositorySrc) Get(ctx context.Context, query *dto.GormQuery) (*[]entity.Discount, error) {
	var discount []entity.Discount
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Find(&discount)
	if result.Error != nil {
		if errors.Is(gorm.ErrRecordNotFound, result.Error) {
			return &discount, nil
		} else {
			return nil, result.Error
		}
	}
	return &discount, nil
}

func (r *DiscountRepositorySrc) Find(ctx context.Context, query *dto.GormQuery) (*entity.Discount, error) {
	var discount entity.Discount
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Find(&discount)
	if result.Error != nil {
		if errors.Is(gorm.ErrRecordNotFound, result.Error) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}
	return &discount, nil
}

func (r *DiscountRepositorySrc) Save(ctx context.Context, discount *entity.Discount) error {
	result := r.pg.Db.WithContext(ctx).Save(discount)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *DiscountRepositorySrc) Delete(ctx context.Context, query *dto.GormQuery) error {
	var discount entity.Discount
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Delete(&discount)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
