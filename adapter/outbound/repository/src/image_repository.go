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

var ImageModule = fx.Module(
	"image_repository",
	fx.Provide(NewImageRepositorySrc),
	fx.Provide(func(p *ImageRepositorySrc) repository.ImageRepository { return p }),
)

type ImageRepositorySrc struct {
	pg     *postgres.Postgres
	logger *zap.SugaredLogger
}

func NewImageRepositorySrc(
	pg *postgres.Postgres,
	logger *zap.SugaredLogger,
) (*ImageRepositorySrc, error) {
	return &ImageRepositorySrc{
		pg:     pg,
		logger: logger,
	}, nil
}

func (r *ImageRepositorySrc) Create(ctx context.Context, image *entity.Image) error {
	err := r.pg.Db.WithContext(ctx).Create(image).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ImageRepositorySrc) Get(ctx context.Context, query *dto.GormQuery) (*[]entity.Image, error) {
	var items []entity.Image
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

func (r *ImageRepositorySrc) Find(ctx context.Context, query *dto.GormQuery) (*entity.Image, error) {
	var items entity.Image
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

func (r *ImageRepositorySrc) Save(ctx context.Context, image *entity.Image) error {
	result := r.pg.Db.WithContext(ctx).Save(image)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ImageRepositorySrc) Delete(ctx context.Context, query *dto.GormQuery) error {
	var item entity.Image
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Delete(&item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
