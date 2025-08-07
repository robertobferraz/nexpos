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

type City struct {
	Base     `json:",inline" valid:"-"`
	Name     *string `json:"name" valid:"required"`
	StateID  *string `json:"-" valid:"-"`
	State    *State  `json:"state" valid:"-"`
	StreetID *string `json:"-" valid:"-"`
	Street   *Street `json:"street" valid:"-"`
}

func NewCity(name *string, street *Street, state *State) (*City, error) {
	city := &City{
		Name:     name,
		StreetID: street.ID,
		Street:   street,
		StateID:  state.ID,
		State:    state,
	}

	city.ID = utils.PString(uuid.NewV4().String())
	city.CreatedAt = utils.PTime(time.Now())

	if err := city.isValid(); err != nil {
		return nil, err
	}

	return city, nil
}

func (p *City) isValid() error {
	_, err := govalidator.ValidateStruct(p)
	return err
}
