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

var OrderModule = fx.Module(
	"order_repository",
	fx.Provide(NewOrderRepositorySrc),
	fx.Provide(func(p *OrderRepositorySrc) repository.OrderRepository { return p }),
)

type OrderRepositorySrc struct {
	pg     *postgres.Postgres
	logger *zap.SugaredLogger
}

func NewOrderRepositorySrc(
	pg *postgres.Postgres,
	logger *zap.SugaredLogger,
) (*OrderRepositorySrc, error) {
	return &OrderRepositorySrc{
		pg:     pg,
		logger: logger,
	}, nil
}

func (r *OrderRepositorySrc) Create(ctx context.Context, order *entity.Order) error {
	err := r.pg.Db.WithContext(ctx).Create(order).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderRepositorySrc) Get(ctx context.Context, query *dto.GormQuery) (*[]entity.Order, error) {
	var items []entity.Order
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Find(&items)
	if result.Error != nil {
		if errors.Is(gorm.ErrRecordNotFound, result.Error) {
			return &items, nil
		} else {
			return nil, result.Error
		}
	}
	return &items, nil
}

func (r *OrderRepositorySrc) Find(ctx context.Context, query *dto.GormQuery) (*entity.Order, error) {
	var items entity.Order
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Find(&items)
	if result.Error != nil {
		if errors.Is(gorm.ErrRecordNotFound, result.Error) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}

	if items.ID == nil {
		return nil, nil
	}

	return &items, nil
}

func (r *OrderRepositorySrc) Save(ctx context.Context, order *entity.Order) error {
	result := r.pg.Db.WithContext(ctx).Save(order)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *OrderRepositorySrc) Delete(ctx context.Context, query *dto.GormQuery) error {
	var item entity.Order
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Delete(&item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
