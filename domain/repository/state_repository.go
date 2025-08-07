package repository

import (
	"context"

	"github.com/robertobff/nexpos/domain/dto"
	"github.com/robertobff/nexpos/domain/entity"
)

type StateRepository interface {
	Create(context.Context, *entity.State) error
	Get(context.Context, *dto.GormQuery) (*[]entity.State, error)
	Find(context.Context, *dto.GormQuery) (*entity.State, error)
	Save(context.Context, *entity.State) error
	Delete(context.Context, *dto.GormQuery) error
}
