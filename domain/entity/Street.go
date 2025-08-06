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

type Street struct {
	Base       `json:",inline" valid:"-"`
	Name       *string   `json:"name" valid:"required"`
	DistrictID *string   `json:"-" valid:"-"`
	District   *District `json:"district" valid:"-"`
	ZipCode    *string   `json:"zip_code" valid:"-"`
	Number     *string   `json:"number" valid:"-"`
}

func NewStreet(name, zipCode, number *string, district *District) (*Street, error) {
	street := &Street{
		Name:       name,
		ZipCode:    zipCode,
		Number:     number,
		DistrictID: district.ID,
		District:   district,
	}

	street.ID = utils.PString(uuid.NewV4().String())
	street.CreatedAt = utils.PTime(time.Now())

	if err := street.isValid(); err != nil {
		return nil, err
	}

	return street, nil
}

func (p *Street) isValid() error {
	_, err := govalidator.ValidateStruct(p)
	return err
}
