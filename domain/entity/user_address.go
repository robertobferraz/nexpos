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

type UserAddress struct {
	Base   `json:",inline" valid:"-"`
	UserID *string `json:"-" valid:"-"`
	User   *User   `json:"user" valid:"-"`
	CityID *string `json:"-" valid:"-"`
	City   *City   `json:"city" valid:"-"`
}

func NewUserAddress(user *User, city *City) (*UserAddress, error) {
	userAddress := &UserAddress{
		UserID: user.ID,
		User:   user,
		CityID: city.ID,
		City:   city,
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
