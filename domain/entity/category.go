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
	Name        *string   `json:"name" valid:"-"`
	Description *string   `json:"description" valid:"-"`
	ImageID     *string   `json:"-" valid:"-"`
	Image       *Image    `json:"image" valid:"-"`
	ParentID    *string   `json:"-" valid:"-"`
	Parent      *Category `json:"parent" valid:"-"`
}

func NewCategory(name, description *string, image *Image, parent *Category) (*Category, error) {
	var imgID *string
	if image != nil {
		imgID = image.ID
	}

	category := &Category{
		Name:        name,
		Description: description,
		ImageID:     imgID,
		Image:       image,
		ParentID:    parent.ID,
		Parent:      parent,
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
