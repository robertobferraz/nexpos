package src

import (
	"context"
	"errors"
	"github.com/robertobff/food-service/adapter/outbound/database/postgres"
	"github.com/robertobff/food-service/domain/dto"
	"github.com/robertobff/food-service/domain/entity"
	"github.com/robertobff/food-service/domain/repository"
	"gorm.io/gorm"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var ItemModule = fx.Module(
	"item_repository",
	fx.Provide(NewItemRepositorySrc),
	fx.Provide(func(p *ItemRepositorySrc) repository.ItemRepository { return p }),
)

type ItemRepositorySrc struct {
	pg     *postgres.Postgres
	logger *zap.SugaredLogger
}

func NewItemRepositorySrc(
	pg *postgres.Postgres,
	logger *zap.SugaredLogger,
) (*ItemRepositorySrc, error) {
	return &ItemRepositorySrc{
		pg:     pg,
		logger: logger,
	}, nil
}

func (r *ItemRepositorySrc) Create(ctx context.Context, item *entity.Item) error {
	err := r.pg.Db.WithContext(ctx).Create(item).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ItemRepositorySrc) Get(ctx context.Context, query *dto.GormQuery) (*[]entity.Item, error) {
	var items []entity.Item
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

func (r *ItemRepositorySrc) Find(ctx context.Context, query *dto.GormQuery) (*entity.Item, error) {
	var items entity.Item
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

func (r *ItemRepositorySrc) Save(ctx context.Context, item *entity.Item) error {
	result := r.pg.Db.WithContext(ctx).Save(item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ItemRepositorySrc) Delete(ctx context.Context, query *dto.GormQuery) error {
	var item entity.Item
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Delete(&item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
