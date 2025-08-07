package repository

import (
	"context"

	"github.com/robertobff/nexpos/domain/dto"
	"github.com/robertobff/nexpos/domain/entity"
)

type StreetRepository interface {
	Create(context.Context, *entity.Street) error
	Get(context.Context, *dto.GormQuery) (*[]entity.Street, error)
	Find(context.Context, *dto.GormQuery) (*entity.Street, error)
	Save(context.Context, *entity.Street) error
	Delete(context.Context, *dto.GormQuery) error
}
