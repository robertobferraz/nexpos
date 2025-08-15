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
	ImageID     *string    `json:"-" valid:"-"`
	Image       *Image     `json:"image" valid:"-"`
	BirthDate   *time.Time `json:"birth_date" valid:"-"`
	Cpf         *string    `json:"cpf" valid:"-"`
	ExternalID  *string    `json:"external_id" valid:"-"`
}

func NewUser(name, username, email, cpf, phoneNumber, birthdate, externalID *string, image *Image) (*User, error) {
	if cpf != nil {
		isValid := utils.CpfValidator(*cpf)
		if !isValid {
			return nil, fmt.Errorf("invalid cpf")
		}
	}

	var bDate time.Time
	var err error
	if birthdate != nil {
		bDate, err = time.Parse("2006-01-02", *birthdate)
		if err != nil {
			return nil, err
		}
	}

	var imageID *string
	if image != nil {
		imageID = image.ID
	}
	var date *time.Time
	if bDate.Before(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)) {
		date = nil
	} else {
		date = utils.PTime(bDate)
	}

	user := &User{
		Name:        name,
		Username:    username,
		Email:       email,
		BirthDate:   date,
		Cpf:         cpf,
		ImageID:     imageID,
		Image:       image,
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

func (p *User) SetName(name *string) error {
	p.Name = name
	return p.isValid()
}

func (p *User) SetUsername(username *string) error {
	p.Username = username
	return p.isValid()
}

func (p *User) SetEmail(email *string) error {
	p.Email = email
	return p.isValid()
}

func (p *User) SetPhoneNumber(phoneNumber *string) error {
	p.PhoneNumber = phoneNumber
	return p.isValid()
}

func (p *User) SetImage(image *Image) error {
	p.Image = image
	return p.isValid()
}

func (p *User) SetBirthDate(birthDate *string) error {
	bDate, err := time.Parse("2006-01-02", *birthDate)
	if err != nil {
		return err
	}

	p.BirthDate = utils.PTime(bDate)
	return p.isValid()
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
