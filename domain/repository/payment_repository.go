package repository

import (
	"context"

	"github.com/robertobff/nexpos/domain/dto"
	"github.com/robertobff/nexpos/domain/entity"
)

type PaymentRepository interface {
	Create(context.Context, *entity.Payment) error
	Get(context.Context, *dto.GormQuery) (*[]entity.Payment, error)
	Find(context.Context, *dto.GormQuery) (*entity.Payment, error)
	Save(context.Context, *entity.Payment) error
	Delete(context.Context, *dto.GormQuery) error
}
