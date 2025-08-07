package entity

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/robertobff/nexpos/utils"
	uuid "github.com/satori/go.uuid"
)

func init() {
	govalidator.SetFieldsRequiredByDefault(true)

	govalidator.TagMap["password_strength"] = govalidator.Validator(func(str string) bool {
		return len(str) >= 8 && regexp.MustCompile(`[A-Za-z]`).MatchString(str) && regexp.MustCompile(`[0-9]`).MatchString(str)
	})

	govalidator.TagMap["phone"] = govalidator.Validator(func(str string) bool {
		return regexp.MustCompile(`^\+?[1-9]\d{1,14}$`).MatchString(str)
	})
}

type User struct {
	Base        `json:",inline" valid:"-"`
	Email       *string    `json:"email" valid:"email,required~Email is invalid or missing"`
	Username    *string    `json:"username" valid:"required"`
	Name        *string    `json:"name" valid:"required~Name is missing,stringlength(2|50)~Name must be 2-50 characters"`
	PhoneNumber *string    `json:"phone_number" valid:"phone~Invalid phone number format,optional"`
	BirthDate   *time.Time `json:"birth_date" valid:"-"`
	Cpf         *string    `json:"cpf" valid:"-"`
	ExternalID  *string    `json:"external_id" valid:"-"`
}

func NewUser(name, username, email, cpf, phoneNumber *string, birthdate *string, externalID *string) (*User, error) {
	if cpf != nil {
		isValid := utils.CpfValidator(*cpf)
		if !isValid {
			return nil, fmt.Errorf("invalid cpf")
		}
	}

	bDate, err := time.Parse("2006-01-02", *birthdate)
	if err != nil {
		return nil, err
	}

	user := &User{
		Name:        name,
		Username:    username,
		Email:       email,
		BirthDate:   utils.PTime(bDate),
		Cpf:         cpf,
		PhoneNumber: phoneNumber,
		ExternalID:  externalID,
	}

	user.ID = utils.PString(uuid.NewV4().String())
	user.CreatedAt = utils.PTime(time.Now())

	if err := user.isValid(); err != nil {
		return nil, err
	}

	return user, nil
}

func (p *User) isValid() error {
	_, err := govalidator.ValidateStruct(p)
	if err != nil {
		var valErrors govalidator.Errors
		if errors.As(err, &valErrors) {
			var errorMessages []string
			for _, e := range valErrors {
				errorMessages = append(errorMessages, e.Error())
			}
			return errors.New(strings.Join(errorMessages, "; "))
		}
		return err
	}

	if p.BirthDate != nil && p.BirthDate.After(time.Now()) {
		return errors.New("birth date cannot be in the future")
	}

	return nil
}
