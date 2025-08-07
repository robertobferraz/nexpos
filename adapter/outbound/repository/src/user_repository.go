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

var UserModule = fx.Module(
	"user_repository",
	fx.Provide(NewUserRepositorySrc),
	fx.Provide(func(p *UserRepositorySrc) repository.UserRepository { return p }),
)

type UserRepositorySrc struct {
	pg     *postgres.Postgres
	logger *zap.SugaredLogger
}

func NewUserRepositorySrc(
	pg *postgres.Postgres,
	logger *zap.SugaredLogger,
) (*UserRepositorySrc, error) {
	return &UserRepositorySrc{
		pg:     pg,
		logger: logger,
	}, nil
}

func (r *UserRepositorySrc) Create(ctx context.Context, user *entity.User) error {
	err := r.pg.Db.WithContext(ctx).Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepositorySrc) Get(ctx context.Context, query *dto.GormQuery) (*[]entity.User, error) {
	var items []entity.User
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

func (r *UserRepositorySrc) Find(ctx context.Context, query *dto.GormQuery) (*entity.User, error) {
	var items entity.User
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

func (r *UserRepositorySrc) Save(ctx context.Context, user *entity.User) error {
	result := r.pg.Db.WithContext(ctx).Save(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *UserRepositorySrc) Delete(ctx context.Context, query *dto.GormQuery) error {
	var item entity.User
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Delete(&item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
