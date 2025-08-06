package repository

import (
	"context"
	"github.com/robertobff/food-service/domain/dto"
	"github.com/robertobff/food-service/domain/entity"
)

type UserOrdersRepository interface {
	Create(context.Context, *entity.UserOrders) error
	Get(context.Context, *dto.GormQuery) (*[]entity.UserOrders, error)
	Find(context.Context, *dto.GormQuery) (*entity.UserOrders, error)
	Save(context.Context, *entity.UserOrders) error
	Delete(context.Context, *dto.GormQuery) error
}
