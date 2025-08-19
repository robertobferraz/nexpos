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

var CartModule = fx.Module(
	"address_repository",
	fx.Provide(NewCartRepositorySrc),
	fx.Provide(func(p *CartRepositorySrc) repository.CartRepository { return p }),
)

type CartRepositorySrc struct {
	pg     *postgres.Postgres
	logger *zap.SugaredLogger
}

func NewCartRepositorySrc(
	pg *postgres.Postgres,
	logger *zap.SugaredLogger,
) (*CartRepositorySrc, error) {
	return &CartRepositorySrc{
		pg:     pg,
		logger: logger,
	}, nil
}

func (r *CartRepositorySrc) Create(ctx context.Context, cart *entity.Cart) error {
	err := r.pg.Db.WithContext(ctx).Create(cart).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CartRepositorySrc) Get(ctx context.Context, query *dto.GormQuery) (*[]entity.Cart, error) {
	var items []entity.Cart
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

func (r *CartRepositorySrc) Find(ctx context.Context, query *dto.GormQuery) (*entity.Cart, error) {
	var items entity.Cart
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

func (r *CartRepositorySrc) Save(ctx context.Context, cart *entity.Cart) error {
	result := r.pg.Db.WithContext(ctx).Save(cart)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *CartRepositorySrc) Delete(ctx context.Context, query *dto.GormQuery) error {
	var item entity.Cart
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Delete(&item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
