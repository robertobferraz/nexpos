package repository

import (
	"context"
	"github.com/robertobff/food-service/domain/dto"
	"github.com/robertobff/food-service/domain/entity"
)

type CountryRepository interface {
	Create(context.Context, *entity.Country) error
	Get(context.Context, *dto.GormQuery) (*[]entity.Country, error)
	Find(context.Context, *dto.GormQuery) (*entity.Country, error)
	Save(context.Context, *entity.Country) error
	Delete(context.Context, *dto.GormQuery) error
}
