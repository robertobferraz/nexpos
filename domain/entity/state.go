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
	Iso2       *string  `json:"iso2" valid:"required"`
	CountryID  *string  `json:"-" valid:"-"`
	Country    *Country `json:"country" valid:"-"`
	ExternalID *int     `json:"-" valid:"-"`
}

func NewState(name, iso2 *string, externalId *int, country *Country) (*State, error) {
	state := &State{
		Name:       name,
		Iso2:       iso2,
		CountryID:  country.ID,
		Country:    country,
		ExternalID: externalId,
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
