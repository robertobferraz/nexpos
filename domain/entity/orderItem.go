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

type OrderItem struct {
	Base      `json:",inline" valid:"-"`
	OrderID   *string  `json:"-" valid:"-"`
	Order     *Order   `json:"order" valid:"-"`
	ItemID    *string  `json:"-" valid:"-"`
	Item      *Item    `json:"item" valid:"-"`
	Quantity  *int     `json:"-" valid:"-"`
	UnitPrice *float64 `json:"unit_price" valid:"-"`
	Discount  *float64 `json:"discount" valid:"-"`
}

func NewOrderItem(
	order *Order,
	item *Item,
	quantity *int,
	unitPrice,
	discount *float64,
) (*OrderItem, error) {
	orderItem := &OrderItem{
		OrderID:   order.ID,
		Order:     order,
		ItemID:    item.ID,
		Item:      item,
		Quantity:  quantity,
		UnitPrice: unitPrice,
		Discount:  discount,
	}

	orderItem.ID = utils.PString(uuid.NewV4().String())
	orderItem.CreatedAt = utils.PTime(time.Now())

	if err := orderItem.isValid(); err != nil {
		return nil, err
	}

	return orderItem, nil
}

func (p *OrderItem) isValid() error {
	_, err := govalidator.ValidateStruct(p)
	return err
}
