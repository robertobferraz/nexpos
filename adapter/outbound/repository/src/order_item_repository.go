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

var OrderItemModule = fx.Module(
	"order_repository",
	fx.Provide(NewOrderItemRepositorySrc),
	fx.Provide(func(p *OrderItemRepositorySrc) repository.OrderItemRepository { return p }),
)

type OrderItemRepositorySrc struct {
	pg     *postgres.Postgres
	logger *zap.SugaredLogger
}

func NewOrderItemRepositorySrc(
	pg *postgres.Postgres,
	logger *zap.SugaredLogger,
) (*OrderItemRepositorySrc, error) {
	return &OrderItemRepositorySrc{
		pg:     pg,
		logger: logger,
	}, nil
}

func (r *OrderItemRepositorySrc) Create(ctx context.Context, orderItem *entity.OrderItem) error {
	err := r.pg.Db.WithContext(ctx).Create(orderItem).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *OrderItemRepositorySrc) Get(ctx context.Context, query *dto.GormQuery) (*[]entity.OrderItem, error) {
	var items []entity.OrderItem
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

func (r *OrderItemRepositorySrc) Find(ctx context.Context, query *dto.GormQuery) (*entity.OrderItem, error) {
	var items entity.OrderItem
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

func (r *OrderItemRepositorySrc) Save(ctx context.Context, orderItem *entity.OrderItem) error {
	result := r.pg.Db.WithContext(ctx).Save(orderItem)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *OrderItemRepositorySrc) Delete(ctx context.Context, query *dto.GormQuery) error {
	var item entity.OrderItem
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Delete(&item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
