package entity

import (
	"time"

	"github.com/asaskevich/govalidator"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type BaseTimestamps struct {
	CreatedAt *time.Time `json:"created_at" valid:"required"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" valid:"-"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" valid:"-"`
}

type BaseID struct {
	ID *string `json:"id" valid:"uuid"`
}

type Base struct {
	BaseID         `json:",inline"`
	BaseTimestamps `json:",inline"`
}
