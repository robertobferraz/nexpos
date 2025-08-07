package src

import (
	"context"
	"errors"
	"fmt"

	"github.com/robertobff/nexpos/domain/dto"

	"gorm.io/gorm"
)

func QueryConstructor(db *gorm.DB, query *dto.GormQuery) *gorm.DB {
	if query.Where != nil {
		for _, v := range *query.Where {
			db = db.Where(fmt.Sprint(v.Column, " ", v.Condition, " ?"), v.Value)
		}
	}

	if query.Preload != nil {
		for _, p := range *query.Preload {
			db = db.Preload(p.Field)
		}
	}

	if query.Order != nil {
		for _, o := range *query.Order {
			db = db.Order(o.Field)
		}
	}

	if query.InnerJoins != nil {
		for _, ij := range *query.InnerJoins {
			if ij.Where != nil {
				dbs := []*gorm.DB{}
				for _, ijw := range *ij.Where {
					dbs = append(dbs, Where(db, ijw))
				}
				db = db.InnerJoins(ij.Field, dbs)
			} else {
				db = db.InnerJoins(ij.Field)
			}
		}
	}

	if query.Debug {
		db = db.Debug()
	}

	if query.Unscoped {
		db = db.Unscoped()
	}

	return db
}

func Where(db *gorm.DB, query dto.GormWhere) *gorm.DB {
	return db.Where(fmt.Sprint(query.Column, " ", query.Condition, " ?"), query.Value)
}

func Find(item any, db *gorm.DB, ctx context.Context, query *dto.GormQuery) (bool, error) {
	gormDB := QueryConstructor(db, query)
	result := gormDB.WithContext(ctx).First(item)
	if result.Error != nil {
		if errors.Is(gorm.ErrRecordNotFound, result.Error) {
			return true, nil
		} else {
			return false, result.Error
		}
	}
	return false, nil
}

func Get(item any, db *gorm.DB, ctx context.Context, query *dto.GormQuery) (bool, error) {
	gormDB := QueryConstructor(db, query)
	result := gormDB.WithContext(ctx).Find(item)
	if result.Error != nil {
		if errors.Is(gorm.ErrRecordNotFound, result.Error) {
			return true, nil
		} else {
			return false, result.Error
		}
	}
	return false, nil
}
