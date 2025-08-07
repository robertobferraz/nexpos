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

var StateModule = fx.Module(
	"state_repository",
	fx.Provide(NewStateRepositorySrc),
	fx.Provide(func(p *StateRepositorySrc) repository.StateRepository { return p }),
)

type StateRepositorySrc struct {
	pg     *postgres.Postgres
	logger *zap.SugaredLogger
}

func NewStateRepositorySrc(
	pg *postgres.Postgres,
	logger *zap.SugaredLogger,
) (*StateRepositorySrc, error) {
	return &StateRepositorySrc{
		pg:     pg,
		logger: logger,
	}, nil
}

func (r *StateRepositorySrc) Create(ctx context.Context, state *entity.State) error {
	err := r.pg.Db.WithContext(ctx).Create(state).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *StateRepositorySrc) Get(ctx context.Context, query *dto.GormQuery) (*[]entity.State, error) {
	var items []entity.State
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

func (r *StateRepositorySrc) Find(ctx context.Context, query *dto.GormQuery) (*entity.State, error) {
	var items entity.State
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

func (r *StateRepositorySrc) Save(ctx context.Context, state *entity.State) error {
	result := r.pg.Db.WithContext(ctx).Save(state)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *StateRepositorySrc) Delete(ctx context.Context, query *dto.GormQuery) error {
	var item entity.State
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Delete(&item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
