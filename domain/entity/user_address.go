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

type UserAddress struct {
	Base     `json:",inline" valid:"-"`
	UserID   *string `json:"-" valid:"-"`
	User     *User   `json:"user" valid:"-"`
	StreetID *string `json:"-" valid:"-"`
	Street   *Street `json:"street" valid:"-"`
}

func NewUserAddress(user *User, street *Street) (*UserAddress, error) {
	userAddress := &UserAddress{
		UserID:   user.ID,
		User:     user,
		StreetID: street.ID,
		Street:   street,
	}

	userAddress.ID = utils.PString(uuid.NewV4().String())
	userAddress.CreatedAt = utils.PTime(time.Now())

	if err := userAddress.isValid(); err != nil {
		return nil, err
	}

	return userAddress, nil
}

func (p *UserAddress) isValid() error {
	_, err := govalidator.ValidateStruct(p)
	return err
}
