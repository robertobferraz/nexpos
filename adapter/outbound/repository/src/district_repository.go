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

var DistrictModule = fx.Module(
	"district_repository",
	fx.Provide(NewDistrictRepositorySrc),
	fx.Provide(func(p *DistrictRepositorySrc) repository.DistrictRepository { return p }),
)

type DistrictRepositorySrc struct {
	pg     *postgres.Postgres
	logger *zap.SugaredLogger
}

func NewDistrictRepositorySrc(
	pg *postgres.Postgres,
	logger *zap.SugaredLogger,
) (*DistrictRepositorySrc, error) {
	return &DistrictRepositorySrc{
		pg:     pg,
		logger: logger,
	}, nil
}

func (r *DistrictRepositorySrc) Create(ctx context.Context, district *entity.District) error {
	err := r.pg.Db.WithContext(ctx).Create(district).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *DistrictRepositorySrc) Get(ctx context.Context, query *dto.GormQuery) (*[]entity.District, error) {
	var items []entity.District
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

func (r *DistrictRepositorySrc) Find(ctx context.Context, query *dto.GormQuery) (*entity.District, error) {
	var items entity.District
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

func (r *DistrictRepositorySrc) Save(ctx context.Context, district *entity.District) error {
	result := r.pg.Db.WithContext(ctx).Save(district)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *DistrictRepositorySrc) Delete(ctx context.Context, query *dto.GormQuery) error {
	var item entity.District
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Delete(&item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
