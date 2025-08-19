package repository

import (
	"context"

	"github.com/robertobff/nexpos/domain/dto"
	"github.com/robertobff/nexpos/domain/entity"
)

type CartItemRepository interface {
	Create(context.Context, *entity.CartItem) error
	Get(context.Context, *dto.GormQuery) (*[]entity.CartItem, error)
	Find(context.Context, *dto.GormQuery) (*entity.CartItem, error)
	Save(context.Context, *entity.CartItem) error
	Delete(context.Context, *dto.GormQuery) error
}
