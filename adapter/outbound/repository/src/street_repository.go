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

var StreetModule = fx.Module(
	"street_repository",
	fx.Provide(NewStreetRepositorySrc),
	fx.Provide(func(p *StreetRepositorySrc) repository.StreetRepository { return p }),
)

type StreetRepositorySrc struct {
	pg     *postgres.Postgres
	logger *zap.SugaredLogger
}

func NewStreetRepositorySrc(
	pg *postgres.Postgres,
	logger *zap.SugaredLogger,
) (*StreetRepositorySrc, error) {
	return &StreetRepositorySrc{
		pg:     pg,
		logger: logger,
	}, nil
}

func (r *StreetRepositorySrc) Create(ctx context.Context, street *entity.Street) error {
	err := r.pg.Db.WithContext(ctx).Create(street).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *StreetRepositorySrc) Get(ctx context.Context, query *dto.GormQuery) (*[]entity.Street, error) {
	var items []entity.Street
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

func (r *StreetRepositorySrc) Find(ctx context.Context, query *dto.GormQuery) (*entity.Street, error) {
	var items entity.Street
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

func (r *StreetRepositorySrc) Save(ctx context.Context, street *entity.Street) error {
	result := r.pg.Db.WithContext(ctx).Save(street)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *StreetRepositorySrc) Delete(ctx context.Context, query *dto.GormQuery) error {
	var item entity.Street
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Delete(&item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
