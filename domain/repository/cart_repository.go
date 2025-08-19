package repository

import (
	"context"

	"github.com/robertobff/nexpos/domain/dto"
	"github.com/robertobff/nexpos/domain/entity"
)

type CartRepository interface {
	Create(context.Context, *entity.Cart) error
	Get(context.Context, *dto.GormQuery) (*[]entity.Cart, error)
	Find(context.Context, *dto.GormQuery) (*entity.Cart, error)
	Save(context.Context, *entity.Cart) error
	Delete(context.Context, *dto.GormQuery) error
}
