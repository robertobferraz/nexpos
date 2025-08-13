package repository

import (
	"context"

	"github.com/robertobff/nexpos/domain/dto"
	"github.com/robertobff/nexpos/domain/entity"
)

type ImageRepository interface {
	Create(context.Context, *entity.Image) error
	Get(context.Context, *dto.GormQuery) (*[]entity.Image, error)
	Find(context.Context, *dto.GormQuery) (*entity.Image, error)
	Save(context.Context, *entity.Image) error
	Delete(context.Context, *dto.GormQuery) error
}
