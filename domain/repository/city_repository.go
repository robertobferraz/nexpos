package repository

import (
	"context"
	"github.com/robertobff/food-service/domain/dto"
	"github.com/robertobff/food-service/domain/entity"
)

type CityRepository interface {
	Create(context.Context, *entity.City) error
	Get(context.Context, *dto.GormQuery) (*[]entity.City, error)
	Find(context.Context, *dto.GormQuery) (*entity.City, error)
	Save(context.Context, *entity.City) error
	Delete(context.Context, *dto.GormQuery) error
}
