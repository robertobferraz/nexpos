package dto

import (
	"time"

	"github.com/robertobff/nexpos/domain/entity"
	"github.com/robertobff/nexpos/domain/errors"
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

type Image struct {
	Name        *string `json:"name"`
	Data        *[]byte `json:"-"`
	Url         *string `json:"url,omitempty"`
	ContentType *string `json:"content_type"`
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

type DeleteUserInDto struct {
	ID *string `json:"id"`
}

type GetUserByUIDInDto struct {
	UID *string `json:"uid"`
}

type CreateCategoryInDto struct {
	Name        *string `json:"name"`
	Image       *Image  `json:"image"`
	Description *string `json:"description"`
}

type DeleteCategoryInDto struct {
	ID *string `json:"id"`
}

type CreateItemInDto struct {
	Name        *string  `json:"name"`
	Image       *Image   `json:"image"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	CategoryID  *string  `json:"category_id"`
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

type SaveUserInDto struct {
	ID          *string `json:"id"`
	Username    *string `json:"username"`
	Name        *string `json:"name"`
	Birthdate   *string `json:"birth_date"`
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
	Image       *Image  `json:"image"`
}

type GetUsersInDto struct {
	Protocol *string `json:"-"`
	HostName *string `json:"-"`
}
type GetUsersOutDto struct {
	ID          *string    `json:"id"`
	Username    *string    `json:"username"`
	Name        *string    `json:"name"`
	Email       *string    `json:"email"`
	Birthdate   *time.Time `json:"birth_date"`
	PhoneNumber *string    `json:"phone_number"`
	Image       *Image     `json:"image"`
}

type CreateCountryInDto struct {
	Name         *string `json:"name"`
	Iso2         *string `json:"iso2"`
	Iso3         *string `json:"iso3"`
	CurrencyCode *string `json:"currency_code"`
	PhoneCode    *string `json:"phone_code"`
	Capital      *string `json:"capital"`
	Emoji        *string `json:"emoji"`
	ExternalID   *int    `json:"external_id"`
}

type CreateStateInDto struct {
	Name       *string `json:"name"`
	Iso2       *string `json:"iso2"`
	CountryID  *string `json:"country_id"`
	ExternalID *int    `json:"external_id"`
}

type CreateCityInDto struct {
	Name       *string `json:"name"`
	StateID    *string `json:"state_id"`
	ExternalID *int    `json:"external_id"`
}

type CreateDistrictInDto struct {
	Name   *string `json:"name"`
	CityID *string `json:"city_id"`
}

type CreateStreetInDto struct {
	Name       *string `json:"name"`
	DistrictID *string `json:"district_id"`
	ZipCode    *string `json:"zip_code"`
	Number     *string `json:"number"`
}

type FindCountryByIdentifierInDto struct {
	Identifier *string `json:"identifier"`
}

type GetUserAddressInDto struct {
	UserID   *string `json:"-"`
	Protocol *string `json:"-"`
	HostName *string `json:"-"`
}

type GetUserAddressOutDto struct {
	User   *GetUsersOutDto `json:"user"`
	Street *entity.Street  `json:"street"`
}

type CreateUserAddressInDto struct {
	UserID       *string            `json:"-"`
	CountryID    *string            `json:"-"`
	ActionBrazil *ActionBrazilInDto `json:"action_brazil"`
	ActionOthers *ActionOthersInDto `json:"action_others"`
}

type ActionBrazilInDto struct {
	CountryID *string `json:"-"`
	UserID    *string `json:"-"`
	Cep       *string `json:"cep"`
	Number    *string `json:"number"`
}

type ActionOthersInDto struct {
	CountryID  *string   `json:"-"`
	UserID     *string   `json:"-"`
	StateID    *string   `json:"state_id"`
	DistrictID *string   `json:"district_id"`
	District   *District `json:"district"`
	StreetID   *string   `json:"street_id"`
	Street     *Street   `json:"street"`
}

type District struct {
	Name   *string `json:"name"`
	CityID *string `json:"city_id"`
}

type Street struct {
	Name       *string `json:"name"`
	ZipCode    *string `json:"zip_code"`
	Number     *string `json:"number"`
	DistrictID *string `json:"district_id"`
}
