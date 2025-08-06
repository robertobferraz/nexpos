package repository

import (
	"context"
	"github.com/robertobff/food-service/domain/dto"
	"github.com/robertobff/food-service/domain/entity"
)

type StateRepository interface {
	Create(context.Context, *entity.State) error
	Get(context.Context, *dto.GormQuery) (*[]entity.State, error)
	Find(context.Context, *dto.GormQuery) (*entity.State, error)
	Save(context.Context, *entity.State) error
	Delete(context.Context, *dto.GormQuery) error
}
