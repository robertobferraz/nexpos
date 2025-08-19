package entity

import (
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/robertobff/nexpos/utils"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.TagMap["orderStatus"] = govalidator.Validator(func(str string) bool {
		res := str == ORDER_STATUS_PENDING.String()
		res = res || str == ORDER_STATUS_SHIPPED.String()
		res = res || str == ORDER_STATUS_DELIVERED.String()
		res = res || str == ORDER_STATUS_CANCELED.String()
		return res
	})

	govalidator.SetFieldsRequiredByDefault(true)
}

type OrderStatus int

const (
	ORDER_STATUS_PENDING OrderStatus = iota
	ORDER_STATUS_SHIPPED
	ORDER_STATUS_DELIVERED
	ORDER_STATUS_CANCELED
)

func newOrderStatus[T OrderStatus | int](orderStatus T) *OrderStatus {
	v := (OrderStatus)(orderStatus)
	return &v
}

func (t OrderStatus) String() string {
	switch t {
	case ORDER_STATUS_PENDING:
		return "pending"
	case ORDER_STATUS_SHIPPED:
		return "shipped"
	case ORDER_STATUS_DELIVERED:
		return "delivered"
	case ORDER_STATUS_CANCELED:
		return "canceled"
	}
	return ""
}

type Order struct {
	Base              `json:",inline" valid:"-"`
	UserID            *string      `json:"-" valid:"-"`
	User              *User        `json:"user" valid:"-"`
	BillingAddressID  *string      `json:"-" valid:"-"`
	BillingAddress    *Address     `json:"billingAddress" valid:"-"`
	ShippingAddressID *string      `json:"-" valid:"-"`
	ShippingAddress   *Address     `json:"shippingAddress" valid:"-"`
	Date              *time.Time   `json:"date" valid:"-"`
	Total             *float64     `json:"total" valid:"-"`
	SubTotal          *float64     `json:"subTotal" valid:"-"`
	DiscountAmount    *float64     `json:"discountAmount" valid:"-"`
	ShippingCost      *float64     `json:"shippingCost" valid:"-"`
	TaxAmount         *float64     `json:"taxAmount" valid:"-"`
	Status            *OrderStatus `json:"status" valid:"orderStatus"`
	PaymentMethod     *string      `json:"paymentMethod" valid:"-"`
	TrackingNumber    *string      `json:"trackingNumber" valid:"-"`
}

func NewOrders(
	user *User,
	billingAddress *Address,
	shippingAddress *Address,
	date *time.Time,
	total,
	subTotal,
	discountAmount,
	shippingCost,
	taxAmount *float64,
	paymentMethod,
	trackingNumber *string,
) (*Order, error) {
	userOrders := &Order{
		UserID:            user.ID,
		User:              user,
		BillingAddressID:  billingAddress.ID,
		BillingAddress:    billingAddress,
		ShippingAddressID: shippingAddress.ID,
		ShippingAddress:   shippingAddress,
		Date:              date,
		Total:             total,
		SubTotal:          subTotal,
		DiscountAmount:    discountAmount,
		ShippingCost:      shippingCost,
		TaxAmount:         taxAmount,
		Status:            newOrderStatus(ORDER_STATUS_PENDING),
		PaymentMethod:     paymentMethod,
		TrackingNumber:    trackingNumber,
	}

	userOrders.ID = utils.PString(uuid.NewV4().String())

	userOrders.CreatedAt = utils.PTime(time.Now())

	if err := userOrders.isValid(); err != nil {
		return nil, err
	}

	return userOrders, nil
}

func (p *Order) isValid() error {
	_, err := govalidator.ValidateStruct(p)
	return err
}
