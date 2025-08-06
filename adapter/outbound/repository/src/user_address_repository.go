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

var UserAddressModule = fx.Module(
	"user_address_repository",
	fx.Provide(NewUserAddressRepositorySrc),
	fx.Provide(func(p *UserAddressRepositorySrc) repository.UserAddressRepository { return p }),
)

type UserAddressRepositorySrc struct {
	pg     *postgres.Postgres
	logger *zap.SugaredLogger
}

func NewUserAddressRepositorySrc(
	pg *postgres.Postgres,
	logger *zap.SugaredLogger,
) (*UserAddressRepositorySrc, error) {
	return &UserAddressRepositorySrc{
		pg:     pg,
		logger: logger,
	}, nil
}

func (r *UserAddressRepositorySrc) Create(ctx context.Context, userAddress *entity.UserAddress) error {
	err := r.pg.Db.WithContext(ctx).Create(userAddress).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *UserAddressRepositorySrc) Get(ctx context.Context, query *dto.GormQuery) (*[]entity.UserAddress, error) {
	var items []entity.UserAddress
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

func (r *UserAddressRepositorySrc) Find(ctx context.Context, query *dto.GormQuery) (*entity.UserAddress, error) {
	var items entity.UserAddress
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

func (r *UserAddressRepositorySrc) Save(ctx context.Context, userAddress *entity.UserAddress) error {
	result := r.pg.Db.WithContext(ctx).Save(userAddress)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *UserAddressRepositorySrc) Delete(ctx context.Context, query *dto.GormQuery) error {
	var item entity.UserAddress
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Delete(&item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
