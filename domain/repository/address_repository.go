package repository

import (
	"context"

	"github.com/robertobff/nexpos/domain/dto"
	"github.com/robertobff/nexpos/domain/entity"
)

type AddressRepository interface {
	Create(context.Context, *entity.Address) error
	Get(context.Context, *dto.GormQuery) (*[]entity.Address, error)
	Find(context.Context, *dto.GormQuery) (*entity.Address, error)
	Save(context.Context, *entity.Address) error
	Delete(context.Context, *dto.GormQuery) error
}
