package repository

import (
	"context"

	"github.com/robertobff/nexpos/domain/dto"
	"github.com/robertobff/nexpos/domain/entity"
)

type ReviewRepository interface {
	Create(context.Context, *entity.Review) error
	Get(context.Context, *dto.GormQuery) (*[]entity.Review, error)
	Find(context.Context, *dto.GormQuery) (*entity.Review, error)
	Save(context.Context, *entity.Review) error
	Delete(context.Context, *dto.GormQuery) error
}
