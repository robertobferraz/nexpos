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

var UserOrdersModule = fx.Module(
	"user_orders_repository",
	fx.Provide(NewUserOrdersRepositorySrc),
	fx.Provide(func(p *UserOrdersRepositorySrc) repository.UserOrdersRepository { return p }),
)

type UserOrdersRepositorySrc struct {
	pg     *postgres.Postgres
	logger *zap.SugaredLogger
}

func NewUserOrdersRepositorySrc(
	pg *postgres.Postgres,
	logger *zap.SugaredLogger,
) (*UserOrdersRepositorySrc, error) {
	return &UserOrdersRepositorySrc{
		pg:     pg,
		logger: logger,
	}, nil
}

func (r *UserOrdersRepositorySrc) Create(ctx context.Context, userOrders *entity.UserOrders) error {
	err := r.pg.Db.WithContext(ctx).Create(userOrders).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *UserOrdersRepositorySrc) Get(ctx context.Context, query *dto.GormQuery) (*[]entity.UserOrders, error) {
	var items []entity.UserOrders
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

func (r *UserOrdersRepositorySrc) Find(ctx context.Context, query *dto.GormQuery) (*entity.UserOrders, error) {
	var items entity.UserOrders
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Find(&items)
	if result.Error != nil {
		if errors.Is(gorm.ErrRecordNotFound, result.Error) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}
	return &items, nil
}

func (r *UserOrdersRepositorySrc) Save(ctx context.Context, userOrders *entity.UserOrders) error {
	result := r.pg.Db.WithContext(ctx).Save(userOrders)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *UserOrdersRepositorySrc) Delete(ctx context.Context, query *dto.GormQuery) error {
	var item entity.UserOrders
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Delete(&item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
