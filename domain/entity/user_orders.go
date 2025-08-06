package entity

import (
	"github.com/asaskevich/govalidator"
	"github.com/robertobff/food-service/utils"
	uuid "github.com/satori/go.uuid"
	"time"
)

func init() {
	govalidator.TagMap["paidStatus"] = govalidator.Validator(func(str string) bool {
		res := str == PAID_STATUS_PENDING.String()
		res = res || str == PAID_STATUS_ACCEPTED.String()
		res = res || str == PAID_STATUS_REFUSED.String()
		res = res || str == PAID_STATUS_CANCELED.String()
		res = res || str == PAID_STATUS_DELETED.String()
		res = res || str == PAID_STATUS_REFUND.String()
		return res
	})

	govalidator.SetFieldsRequiredByDefault(true)
}

type PaidStatus int

const (
	PAID_STATUS_PENDING PaidStatus = iota
	PAID_STATUS_ACCEPTED
	PAID_STATUS_REFUSED
	PAID_STATUS_CANCELED
	PAID_STATUS_DELETED
	PAID_STATUS_REFUND
)

func newPaidStatus[T PaidStatus | int](paidStatus T) *PaidStatus {
	v := (PaidStatus)(paidStatus)
	return &v
}

func (t PaidStatus) String() string {
	switch t {
	case PAID_STATUS_PENDING:
		return "pending"
	case PAID_STATUS_ACCEPTED:
		return "accepted"
	case PAID_STATUS_REFUSED:
		return "refused"
	case PAID_STATUS_CANCELED:
		return "canceled"
	case PAID_STATUS_DELETED:
		return "deleted"
	case PAID_STATUS_REFUND:
		return "refund"
	}
	return ""
}

type UserOrders struct {
	Base   `json:",inline" valid:"-"`
	UserID *string     `json:"-" valid:"-"`
	User   *User       `json:"user" valid:"-"`
	Item   []*Item     `json:"item" valid:"-"`
	Date   *time.Time  `json:"date" valid:"-"`
	Total  *float64    `json:"total" valid:"-"`
	Status *PaidStatus `json:"status" valid:"paidStatus,optional"`
}

func NewUserOrders(user *User, date *time.Time, total *float64) (*UserOrders, error) {
	userOrders := &UserOrders{
		UserID: user.ID,
		User:   user,
		Date:   date,
		Total:  total,
		Status: newPaidStatus(PAID_STATUS_PENDING),
	}

	userOrders.ID = utils.PString(uuid.NewV4().String())

	userOrders.CreatedAt = utils.PTime(time.Now())

	if err := userOrders.isValid(); err != nil {
		return nil, err
	}

	return userOrders, nil
}

func (p *UserOrders) isValid() error {
	_, err := govalidator.ValidateStruct(p)
	return err
}

func (p *UserOrders) SetItems(items ...*Item) error {
	p.Item = items
	p.UpdatedAt = utils.PTime(time.Now())
	return p.isValid()
}
