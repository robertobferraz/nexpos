package repository

import (
	"context"

	"github.com/robertobff/nexpos/domain/dto"
	"github.com/robertobff/nexpos/domain/entity"
)

type ItemRepository interface {
	Create(context.Context, *entity.Item) error
	Get(context.Context, *dto.GormQuery) (*[]entity.Item, error)
	Find(context.Context, *dto.GormQuery) (*entity.Item, error)
	Save(context.Context, *entity.Item) error
	Delete(context.Context, *dto.GormQuery) error
}
