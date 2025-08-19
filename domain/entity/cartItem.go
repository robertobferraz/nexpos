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

type CartItem struct {
	Base     `json:",inline" valid:"-"`
	CartID   *string `json:"-" valid:"-"`
	Cart     *Cart   `json:"cart" valid:"-"`
	ItemID   *string `json:"-" valid:"-"`
	Item     *Item   `json:"item" valid:"-"`
	Quantity *uint   `json:"quantity" valid:"-"`
}

func NewCartItem(cart *Cart, item *Item, quantity *uint) (*CartItem, error) {
	cartItem := &CartItem{
		CartID:   cart.ID,
		Cart:     cart,
		ItemID:   item.ID,
		Item:     item,
		Quantity: quantity,
	}

	cartItem.ID = utils.PString(uuid.NewV4().String())
	cartItem.CreatedAt = utils.PTime(time.Now())

	if err := cart.isValid(); err != nil {
		return nil, err
	}

	return cartItem, nil
}

func (p *CartItem) isValid() error {
	_, err := govalidator.ValidateStruct(p)
	return err
}
