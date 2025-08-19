package repository

import (
	"context"

	"github.com/robertobff/nexpos/domain/dto"
	"github.com/robertobff/nexpos/domain/entity"
)

type OrderRepository interface {
	Create(context.Context, *entity.Order) error
	Get(context.Context, *dto.GormQuery) (*[]entity.Order, error)
	Find(context.Context, *dto.GormQuery) (*entity.Order, error)
	Save(context.Context, *entity.Order) error
	Delete(context.Context, *dto.GormQuery) error
}
