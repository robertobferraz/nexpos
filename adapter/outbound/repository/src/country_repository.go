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

var CountryModule = fx.Module(
	"country_repository",
	fx.Provide(NewCountryRepositorySrc),
	fx.Provide(func(p *CountryRepositorySrc) repository.CountryRepository { return p }),
)

type CountryRepositorySrc struct {
	pg     *postgres.Postgres
	logger *zap.SugaredLogger
}

func NewCountryRepositorySrc(
	pg *postgres.Postgres,
	logger *zap.SugaredLogger,
) (*CountryRepositorySrc, error) {
	return &CountryRepositorySrc{
		pg:     pg,
		logger: logger,
	}, nil
}

func (r *CountryRepositorySrc) Create(ctx context.Context, Country *entity.Country) error {
	err := r.pg.Db.WithContext(ctx).Create(Country).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CountryRepositorySrc) Get(ctx context.Context, query *dto.GormQuery) (*[]entity.Country, error) {
	var items []entity.Country
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

func (r *CountryRepositorySrc) Find(ctx context.Context, query *dto.GormQuery) (*entity.Country, error) {
	var items entity.Country
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

func (r *CountryRepositorySrc) Save(ctx context.Context, Country *entity.Country) error {
	result := r.pg.Db.WithContext(ctx).Save(Country)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *CountryRepositorySrc) Delete(ctx context.Context, query *dto.GormQuery) error {
	var item entity.Country
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Delete(&item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
