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

type District struct {
	Base   `json:",inline" valid:"-"`
	Name   *string `json:"name" valid:"required"`
	CityID *string `json:"-" valid:"-"`
	City   *City   `json:"city" valid:"-"`
}

func NewDistrict(name *string, city *City) (*District, error) {
	district := &District{
		Name: name,
		City: city,
	}

	district.ID = utils.PString(uuid.NewV4().String())
	district.CreatedAt = utils.PTime(time.Now())

	if err := district.isValid(); err != nil {
		return nil, err
	}

	return district, nil
}

func (p *District) isValid() error {
	_, err := govalidator.ValidateStruct(p)
	return err
}
