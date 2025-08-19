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

var CartItemModule = fx.Module(
	"address_repository",
	fx.Provide(NewCartItemRepositorySrc),
	fx.Provide(func(p *CartItemRepositorySrc) repository.CartItemRepository { return p }),
)

type CartItemRepositorySrc struct {
	pg     *postgres.Postgres
	logger *zap.SugaredLogger
}

func NewCartItemRepositorySrc(
	pg *postgres.Postgres,
	logger *zap.SugaredLogger,
) (*CartItemRepositorySrc, error) {
	return &CartItemRepositorySrc{
		pg:     pg,
		logger: logger,
	}, nil
}

func (r *CartItemRepositorySrc) Create(ctx context.Context, cartItem *entity.CartItem) error {
	err := r.pg.Db.WithContext(ctx).Create(cartItem).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CartItemRepositorySrc) Get(ctx context.Context, query *dto.GormQuery) (*[]entity.CartItem, error) {
	var items []entity.CartItem
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

func (r *CartItemRepositorySrc) Find(ctx context.Context, query *dto.GormQuery) (*entity.CartItem, error) {
	var items entity.CartItem
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

func (r *CartItemRepositorySrc) Save(ctx context.Context, cartItem *entity.CartItem) error {
	result := r.pg.Db.WithContext(ctx).Save(cartItem)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *CartItemRepositorySrc) Delete(ctx context.Context, query *dto.GormQuery) error {
	var item entity.CartItem
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Delete(&item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
