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

type State struct {
	Base       `json:",inline" valid:"-"`
	Name       *string  `json:"name" valid:"required"`
	Identifier *string  `json:"identifier" valid:"required"`
	CountryID  *string  `json:"-" valid:"-"`
	Country    *Country `json:"country" valid:"-"`
}

func NewState(name, identifier *string, country *Country) (*State, error) {
	state := &State{
		Name:       name,
		Identifier: identifier,
		CountryID:  country.ID,
		Country:    country,
	}

	state.ID = utils.PString(uuid.NewV4().String())
	state.CreatedAt = utils.PTime(time.Now())

	if err := state.isValid(); err != nil {
		return nil, err
	}

	return state, nil
}

func (p *State) isValid() error {
	_, err := govalidator.ValidateStruct(p)
	return err
}
