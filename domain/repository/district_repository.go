package repository

import (
	"context"

	"github.com/robertobff/nexpos/domain/dto"
	"github.com/robertobff/nexpos/domain/entity"
)

type DistrictRepository interface {
	Create(context.Context, *entity.District) error
	Get(context.Context, *dto.GormQuery) (*[]entity.District, error)
	Find(context.Context, *dto.GormQuery) (*entity.District, error)
	Save(context.Context, *entity.District) error
	Delete(context.Context, *dto.GormQuery) error
}
