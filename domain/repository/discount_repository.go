package repository

import (
	"context"
	"github.com/robertobff/food-service/domain/dto"
	"github.com/robertobff/food-service/domain/entity"
)

type DiscountRepository interface {
	Create(context.Context, *entity.Discount) error
	Get(context.Context, *dto.GormQuery) (*[]entity.Discount, error)
	Find(context.Context, *dto.GormQuery) (*entity.Discount, error)
	Save(context.Context, *entity.Discount) error
	Delete(context.Context, *dto.GormQuery) error
}
