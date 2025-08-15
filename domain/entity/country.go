package entity

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/robertobff/nexpos/utils"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Country struct {
	Base         `json:",inline" valid:"-"`
	Name         *string `json:"name" valid:"required"`
	Iso2         *string `json:"iso2" valid:"required"`
	Iso3         *string `json:"iso3" valid:"required"`
	PhoneCode    *string `json:"phone_code" valid:"-"`
	Capital      *string `json:"capital" valid:"-"`
	CurrencyCode *string `json:"currency_code" valid:"-"`
	Emoji        *string `json:"emoji" valid:"-"`
	ExternalID   *int    `json:"-" valid:"-"`
}

func NewCountry(name, iso2, iso3, phoneCode, capital, currencyCode, emoji *string, externalID *int) (*Country, error) {
	country := &Country{
		Name:         name,
		Iso2:         iso2,
		Iso3:         iso3,
		PhoneCode:    phoneCode,
		Capital:      capital,
		CurrencyCode: currencyCode,
		Emoji:        emoji,
		ExternalID:   externalID,
	}

	country.ID = utils.PString(uuid.NewV4().String())
	country.CreatedAt = utils.PTime(time.Now())

	if err := country.isValid(); err != nil {
		return nil, err
	}

	return country, nil
}

func (p *Country) isValid() error {
	_, err := govalidator.ValidateStruct(p)
	return err
}
