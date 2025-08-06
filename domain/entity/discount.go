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

type Discount struct {
	Base       `json:",inline" valid:"-"`
	ItemID     *string    `json:"-" valid:"-"`
	Item       *Item      `json:"item" valid:"-"`
	CategoryID *string    `json:"-" valid:"-"`
	Category   *Category  `json:"category" valid:"-"`
	Date       *time.Time `json:"date" valid:"-"`
	Value      *float64   `json:"value" valid:"-"`
}

func NewDiscount(category *Category, item *Item, date *time.Time, value *float64) (*Discount, error) {
	discount := &Discount{
		Date:       date,
		Value:      value,
		CategoryID: nil,
		Category:   category,
		ItemID:     nil,
		Item:       item,
	}
	if category != nil {
		discount.CategoryID = category.ID
	}

	if item != nil {
		discount.ItemID = item.ID
	}

	discount.ID = utils.PString(uuid.NewV4().String())
	discount.CreatedAt = utils.PTime(time.Now())

	if err := discount.isValid(); err != nil {
		return nil, err
	}

	return discount, nil
}

func (p *Discount) isValid() error {
	_, err := govalidator.ValidateStruct(p)
	return err
}
