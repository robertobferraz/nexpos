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

var AddressModule = fx.Module(
	"address_repository",
	fx.Provide(NewAddressRepositorySrc),
	fx.Provide(func(p *AddressRepositorySrc) repository.AddressRepository { return p }),
)

type AddressRepositorySrc struct {
	pg     *postgres.Postgres
	logger *zap.SugaredLogger
}

func NewAddressRepositorySrc(
	pg *postgres.Postgres,
	logger *zap.SugaredLogger,
) (*AddressRepositorySrc, error) {
	return &AddressRepositorySrc{
		pg:     pg,
		logger: logger,
	}, nil
}

func (r *AddressRepositorySrc) Create(ctx context.Context, address *entity.Address) error {
	err := r.pg.Db.WithContext(ctx).Create(address).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *AddressRepositorySrc) Get(ctx context.Context, query *dto.GormQuery) (*[]entity.Address, error) {
	var items []entity.Address
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

func (r *AddressRepositorySrc) Find(ctx context.Context, query *dto.GormQuery) (*entity.Address, error) {
	var items entity.Address
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

func (r *AddressRepositorySrc) Save(ctx context.Context, address *entity.Address) error {
	result := r.pg.Db.WithContext(ctx).Save(address)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *AddressRepositorySrc) Delete(ctx context.Context, query *dto.GormQuery) error {
	var item entity.Address
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Delete(&item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
