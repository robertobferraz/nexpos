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

type Item struct {
	Base        `json:",inline" valid:"-"`
	Name        *string   `json:"name" valid:"-"`
	Description *string   `json:"description" valid:"-"`
	Image       *string   `json:"image" valid:"-"`
	Price       *float64  `json:"price" valid:"-"`
	CategoryID  *string   `json:"-" valid:"-"`
	Category    *Category `json:"category" valid:"-"`
}

func NewItem(name, description, image *string, price *float64, category *Category) (*Item, error) {
	item := &Item{
		Name:        name,
		Description: description,
		Image:       image,
		Price:       price,
		CategoryID:  category.ID,
		Category:    category,
	}

	item.ID = utils.PString(uuid.NewV4().String())
	item.CreatedAt = utils.PTime(time.Now())

	if err := item.isValid(); err != nil {
		return nil, err
	}

	return item, nil
}

func (p *Item) isValid() error {
	_, err := govalidator.ValidateStruct(p)
	return err
}
