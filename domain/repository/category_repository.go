package repository

import (
	"context"

	"github.com/robertobff/nexpos/domain/dto"
	"github.com/robertobff/nexpos/domain/entity"
)

type CategoryRepository interface {
	Create(context.Context, *entity.Category) error
	Get(context.Context, *dto.GormQuery) (*[]entity.Category, error)
	Find(context.Context, *dto.GormQuery) (*entity.Category, error)
	Save(context.Context, *entity.Category) error
	Delete(context.Context, *dto.GormQuery) error
}
