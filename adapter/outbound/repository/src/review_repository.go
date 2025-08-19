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

var ReviewModule = fx.Module(
	"review_repository",
	fx.Provide(NewReviewRepositorySrc),
	fx.Provide(func(p *ReviewRepositorySrc) repository.ReviewRepository { return p }),
)

type ReviewRepositorySrc struct {
	pg     *postgres.Postgres
	logger *zap.SugaredLogger
}

func NewReviewRepositorySrc(
	pg *postgres.Postgres,
	logger *zap.SugaredLogger,
) (*ReviewRepositorySrc, error) {
	return &ReviewRepositorySrc{
		pg:     pg,
		logger: logger,
	}, nil
}

func (r *ReviewRepositorySrc) Create(ctx context.Context, review *entity.Review) error {
	err := r.pg.Db.WithContext(ctx).Create(review).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *ReviewRepositorySrc) Get(ctx context.Context, query *dto.GormQuery) (*[]entity.Review, error) {
	var items []entity.Review
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

func (r *ReviewRepositorySrc) Find(ctx context.Context, query *dto.GormQuery) (*entity.Review, error) {
	var items entity.Review
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

func (r *ReviewRepositorySrc) Save(ctx context.Context, review *entity.Review) error {
	result := r.pg.Db.WithContext(ctx).Save(review)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *ReviewRepositorySrc) Delete(ctx context.Context, query *dto.GormQuery) error {
	var item entity.Review
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Delete(&item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
