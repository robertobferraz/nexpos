package repository

import (
	"context"
	"github.com/robertobff/food-service/domain/dto"
	"github.com/robertobff/food-service/domain/entity"
)

type UserAddressRepository interface {
	Create(context.Context, *entity.UserAddress) error
	Get(context.Context, *dto.GormQuery) (*[]entity.UserAddress, error)
	Find(context.Context, *dto.GormQuery) (*entity.UserAddress, error)
	Save(context.Context, *entity.UserAddress) error
	Delete(context.Context, *dto.GormQuery) error
}
