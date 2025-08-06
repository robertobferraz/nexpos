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

var CityModule = fx.Module(
	"city_repository",
	fx.Provide(NewCityRepositorySrc),
	fx.Provide(func(p *CityRepositorySrc) repository.CityRepository { return p }),
)

type CityRepositorySrc struct {
	pg     *postgres.Postgres
	logger *zap.SugaredLogger
}

func NewCityRepositorySrc(
	pg *postgres.Postgres,
	logger *zap.SugaredLogger,
) (*CityRepositorySrc, error) {
	return &CityRepositorySrc{
		pg:     pg,
		logger: logger,
	}, nil
}

func (r *CityRepositorySrc) Create(ctx context.Context, city *entity.City) error {
	err := r.pg.Db.WithContext(ctx).Create(city).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CityRepositorySrc) Get(ctx context.Context, query *dto.GormQuery) (*[]entity.City, error) {
	var items []entity.City
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

func (r *CityRepositorySrc) Find(ctx context.Context, query *dto.GormQuery) (*entity.City, error) {
	var items entity.City
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

func (r *CityRepositorySrc) Save(ctx context.Context, city *entity.City) error {
	result := r.pg.Db.WithContext(ctx).Save(city)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *CityRepositorySrc) Delete(ctx context.Context, query *dto.GormQuery) error {
	var item entity.City
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Delete(&item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
