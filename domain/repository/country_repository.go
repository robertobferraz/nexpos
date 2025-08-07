package repository

import (
	"context"

	"github.com/robertobff/nexpos/domain/dto"
	"github.com/robertobff/nexpos/domain/entity"
)

type CountryRepository interface {
	Create(context.Context, *entity.Country) error
	Get(context.Context, *dto.GormQuery) (*[]entity.Country, error)
	Find(context.Context, *dto.GormQuery) (*entity.Country, error)
	Save(context.Context, *entity.Country) error
	Delete(context.Context, *dto.GormQuery) error
}
