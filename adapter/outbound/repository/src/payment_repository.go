package src

import (
	"context"
	"errors"

	"github.com/robertobff/nexpos/adapter/outbound/database/postgres"
	"github.com/robertobff/nexpos/domain/dto"
	"github.com/robertobff/nexpos/domain/entity"
	"github.com/robertobff/nexpos/domain/repository"
	"gorm.io/gorm"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var PaymentModule = fx.Module(
	"payment_repository",
	fx.Provide(NewPaymentRepositorySrc),
	fx.Provide(func(p *PaymentRepositorySrc) repository.PaymentRepository { return p }),
)

type PaymentRepositorySrc struct {
	pg     *postgres.Postgres
	logger *zap.SugaredLogger
}

func NewPaymentRepositorySrc(
	pg *postgres.Postgres,
	logger *zap.SugaredLogger,
) (*PaymentRepositorySrc, error) {
	return &PaymentRepositorySrc{
		pg:     pg,
		logger: logger,
	}, nil
}

func (r *PaymentRepositorySrc) Create(ctx context.Context, payment *entity.Payment) error {
	err := r.pg.Db.WithContext(ctx).Create(payment).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PaymentRepositorySrc) Get(ctx context.Context, query *dto.GormQuery) (*[]entity.Payment, error) {
	var items []entity.Payment
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Find(&items)
	if result.Error != nil {
		if errors.Is(gorm.ErrRecordNotFound, result.Error) {
			return &items, nil
		} else {
			return nil, result.Error
		}
	}
	return &items, nil
}

func (r *PaymentRepositorySrc) Find(ctx context.Context, query *dto.GormQuery) (*entity.Payment, error) {
	var items entity.Payment
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Find(&items)
	if result.Error != nil {
		if errors.Is(gorm.ErrRecordNotFound, result.Error) {
			return nil, nil
		} else {
			return nil, result.Error
		}
	}

	if items.ID == nil {
		return nil, nil
	}

	return &items, nil
}

func (r *PaymentRepositorySrc) Save(ctx context.Context, payment *entity.Payment) error {
	result := r.pg.Db.WithContext(ctx).Save(payment)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *PaymentRepositorySrc) Delete(ctx context.Context, query *dto.GormQuery) error {
	var item entity.Payment
	gormDB := QueryConstructor(r.pg.Db, query)
	result := gormDB.WithContext(ctx).Delete(&item)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
