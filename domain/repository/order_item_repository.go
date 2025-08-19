package repository

import (
	"context"

	"github.com/robertobff/nexpos/domain/dto"
	"github.com/robertobff/nexpos/domain/entity"
)

type OrderItemRepository interface {
	Create(context.Context, *entity.OrderItem) error
	Get(context.Context, *dto.GormQuery) (*[]entity.OrderItem, error)
	Find(context.Context, *dto.GormQuery) (*entity.OrderItem, error)
	Save(context.Context, *entity.OrderItem) error
	Delete(context.Context, *dto.GormQuery) error
}
