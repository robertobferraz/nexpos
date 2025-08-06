package dto

import (
	"github.com/robertobff/food-service/domain/errors"
	"time"
)

type Base struct {
	Success *bool      `json:"success"`
	Error   *BaseError `json:"error,omitempty"`
	Message *string    `json:"message,omitempty"`
	Data    any        `json:"data,omitempty"`
}

type BaseError struct {
	Code    errors.ErrorCode `json:"code"`
	Message *string          `json:"message"`
}

type CreateUserInDto struct {
	Username    *string `json:"username" validate:"required"`
	Password    *string `json:"password" validate:"required"`
	Name        *string `json:"name" validate:"required"`
	Email       *string `json:"email" validate:"required,email"`
	Birthdate   *string `json:"birth_date" validate:"-"`
	PhoneNumber *string `json:"phone_number" validate:"required"`
	ExternalID  *string `json:"external_id" validate:"-"`
	Cpf         *string `json:"cpf" validate:"-"`
}

type CreateUserOutDto struct {
	ID    *string `json:"id"`
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

type DeleteUserInDto struct {
	ID *string `json:"id"`
}

type CreateCategoryInDto struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Image       *string `json:"image"`
}

type CreateCategoryOutDto struct {
	ID          *string `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Image       *string `json:"image"`
}

type DeleteCategoryInDto struct {
	ID *string `json:"id"`
}

type CreateItemInDto struct {
	Name        *string  `json:"name"`
	Image       *string  `json:"image"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	CategoryID  *string  `json:"category_id"`
}

type CreateItemOutDto struct {
	ID          *string               `json:"id"`
	Name        *string               `json:"name"`
	Description *string               `json:"description"`
	Image       *string               `json:"image"`
	Price       *float64              `json:"price"`
	Category    *CreateCategoryOutDto `json:"category"`
}

type DeleteItemInDto struct {
	ID *string `json:"id"`
}

type CreateDiscountInDto struct {
	ItemID     *string    `json:"item_id"`
	CategoryID *string    `json:"category_id"`
	Value      *float64   `json:"value"`
	Date       *time.Time `json:"date"`
}

type CreateDiscountOutDto struct {
	ID       *string               `json:"id"`
	Category *CreateCategoryOutDto `json:"category,omitempty"`
	Item     *CreateItemOutDto     `json:"item,omitempty"`
	Value    *float64              `json:"value"`
	Date     *time.Time            `json:"date"`
}

type DeleteDiscountInDto struct {
	ID *string `json:"id"`
}

type SignInInDto struct {
	Token *string `json:"token"`
}

type SignInOutDto struct {
	ID   *string `json:"id"`
	Name *string `json:"name"`
}
