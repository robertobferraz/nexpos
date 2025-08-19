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

type Address struct {
	Base       `json:",inline" valid:"-"`
	UserID     *string `json:"-" valid:"-"`
	User       *User   `json:"user" valid:"-"`
	StreetID   *string `json:"-" valid:"-"`
	Street     *Street `json:"street" valid:"-"`
	Number     *uint   `json:"number" valid:"-"`
	Complement *string `json:"complement" valid:"-"`
	IsBilling  *bool   `json:"is_billing" valid:"-"`
	IsShipping *bool   `json:"is_shipping" valid:"-"`
	IsDefault  *bool   `json:"is_default" valid:"-"`
}

func NewAddress(
	user *User,
	street *Street,
	number *uint,
	complement *string,
	isBilling,
	isShipping,
	isDefault *bool,
) (*Address, error) {
	address := &Address{
		UserID:     user.ID,
		User:       user,
		StreetID:   street.ID,
		Street:     street,
		Number:     number,
		Complement: complement,
		IsBilling:  isBilling,
		IsShipping: isShipping,
		IsDefault:  isDefault,
	}

	address.ID = utils.PString(uuid.NewV4().String())
	address.CreatedAt = utils.PTime(time.Now())

	if err := address.isValid(); err != nil {
		return nil, err
	}

	return address, nil
}

func (p *Address) isValid() error {
	_, err := govalidator.ValidateStruct(p)
	return err
}
