package repository

import (
	"context"

	"github.com/robertobff/nexpos/domain/dto"
	"github.com/robertobff/nexpos/domain/entity"
)

type UserRepository interface {
	Create(context.Context, *entity.User) error
	Get(context.Context, *dto.GormQuery) (*[]entity.User, error)
	Find(context.Context, *dto.GormQuery) (*entity.User, error)
	Save(context.Context, *entity.User) error
	Delete(context.Context, *dto.GormQuery) error
}
