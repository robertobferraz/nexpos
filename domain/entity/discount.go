package entity

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/robertobff/nexpos/utils"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.TagMap["discountType"] = govalidator.Validator(func(str string) bool {
		res := str == DISCOUNT_TYPE_FIXED.String()
		res = res || str == DISCOUNT_TYPE_PERCENT.String()

		return res
	})

	govalidator.SetFieldsRequiredByDefault(true)

}

func newDiscountType[T DiscountType | int](discountType T) *DiscountType {
	v := (DiscountType)(discountType)
	return &v
}

const (
	DISCOUNT_TYPE_FIXED DiscountType = iota
	DISCOUNT_TYPE_PERCENT
)

func (t DiscountType) String() string {
	switch t {
	case DISCOUNT_TYPE_FIXED:
		return "fixed"
	case DISCOUNT_TYPE_PERCENT:
		return "percentage"
	}
	return ""
}

type DiscountType int

type Discount struct {
	Base       `json:",inline" valid:"-"`
	Code       *string       `json:"code" valid:"-"`
	Type       *DiscountType `json:"type" valid:"-"`
	Value      *float64      `json:"value" valid:"-"`
	StartDate  *time.Time    `json:"start_date" valid:"-"`
	EndDate    *time.Time    `json:"end_date" valid:"-"`
	CategoryID *string       `json:"-" valid:"-"`
	Category   *Category     `json:"category" valid:"-"`
	ItemID     *string       `json:"-" valid:"-"`
	Item       *Item         `json:"item" valid:"-"`
	MinAmount  *float64      `json:"min_amount" valid:"-"`
	MaxUses    *float64      `json:"max_uses" valid:"-"`
}

func NewDiscount(
	code *string,
	value,
	minAmount,
	maxUses *float64,
	dType *DiscountType,
	startDate,
	endDate *time.Time,
	category *Category,
	item *Item,
) (*Discount, error) {
	discount := &Discount{
		Code:       code,
		Type:       dType,
		Value:      value,
		StartDate:  startDate,
		EndDate:    endDate,
		CategoryID: nil,
		Category:   category,
		ItemID:     nil,
		Item:       item,
		MinAmount:  minAmount,
		MaxUses:    maxUses,
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
