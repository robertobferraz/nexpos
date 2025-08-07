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

var CategoryModule = fx.Module(
	"category_repository",
	fx.Provide(NewCategoryRepositorySrc),
	fx.Provide(func(p *CategoryRepositorySrc) repository.CategoryRepository { return p }),
)

type CategoryRepositorySrc struct {
	pg     *postgres.Postgres
	logger *zap.SugaredLogger
}

func NewCategoryRepositorySrc(
	pg *postgres.Postgres,
	logger *zap.SugaredLogger,
) (*CategoryRepositorySrc, error) {
	return &CategoryRepositorySrc{
		pg:     pg,
		logger: logger,
	}, nil
}

func (r *CategoryRepositorySrc) Create(ctx context.Context, category *entity.Category) error {
	err := r.pg.Db.WithContext(ctx).Create(category).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CategoryRepositorySrc) Get(ctx context.Context, query *dto.GormQuery) (*[]entity.Category, error) {
	var category []entity.Category
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Find(&category)
	if result.Error != nil {
		if errors.Is(gorm.ErrRecordNotFound, result.Error) {
			return &category, nil
		} else {
			return nil, result.Error
		}
	}
	return &category, nil
}

func (r *CategoryRepositorySrc) Find(ctx context.Context, query *dto.GormQuery) (*entity.Category, error) {
	var category entity.Category
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Find(&category)
	if result.Error != nil {
		if errors.Is(gorm.ErrRecordNotFound, result.Error) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}
	return &category, nil
}

func (r *CategoryRepositorySrc) Save(ctx context.Context, category *entity.Category) error {
	result := r.pg.Db.WithContext(ctx).Save(category)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *CategoryRepositorySrc) Delete(ctx context.Context, query *dto.GormQuery) error {
	var category entity.Category
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Delete(&category)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
