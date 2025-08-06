package entity

import (
	"github.com/asaskevich/govalidator"
	"github.com/robertobff/food-service/utils"
	uuid "github.com/satori/go.uuid"
	"time"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

type Country struct {
	Base       `json:",inline" valid:"-"`
	Name       *string `json:"name" valid:"required"`
	Identifier *string `json:"identifier" valid:"required"`
}

func NewCountry(name, identifier *string) (*Country, error) {
	country := &Country{
		Name:       name,
		Identifier: identifier,
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
