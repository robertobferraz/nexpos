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

type Cart struct {
	Base   `json:",inline" valid:"-"`
	UserID *string `json:"-" valid:"-"`
	User   *User   `json:"user" valid:"-"`
}

func NewCart(user *User) (*Cart, error) {
	cart := &Cart{
		UserID: user.ID,
		User:   user,
	}

	cart.ID = utils.PString(uuid.NewV4().String())
	cart.CreatedAt = utils.PTime(time.Now())

	if err := cart.isValid(); err != nil {
		return nil, err
	}

	return cart, nil
}

func (p *Cart) isValid() error {
	_, err := govalidator.ValidateStruct(p)
	return err
}
