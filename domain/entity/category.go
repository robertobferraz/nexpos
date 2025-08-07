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

type Category struct {
	Base        `json:",inline" valid:"-"`
	Name        *string `json:"name" valid:"-"`
	Description *string `json:"description" valid:"-"`
	Image       *string `json:"image" valid:"-"`
}

func NewCategory(name, description, image *string) (*Category, error) {
	category := &Category{
		Name:        name,
		Description: description,
		Image:       image,
	}

	category.ID = utils.PString(uuid.NewV4().String())
	category.CreatedAt = utils.PTime(time.Now())

	if err := category.isValid(); err != nil {
		return nil, err
	}

	return category, nil
}

func (p *Category) isValid() error {
	_, err := govalidator.ValidateStruct(p)
	return err
}
