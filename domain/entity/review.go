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

type Review struct {
	Base    `json:",inline" valid:"-"`
	UserID  *string    `json:"-" valid:"-"`
	User    *User      `json:"user" valid:"-"`
	ItemID  *string    `json:"-" valid:"-"`
	Item    *Item      `json:"item" valid:"-"`
	Rating  *uint      `json:"rating" valid:"-"`
	Comment *string    `json:"comment" valid:"-"`
	Date    *time.Time `json:"date" valid:"-"`
}

func NewReview(
	user *User,
	item *Item,
	rating *uint,
	comment *string,
	date *time.Time,
) (*Review, error) {
	review := &Review{
		UserID:  user.ID,
		User:    user,
		ItemID:  item.ID,
		Item:    item,
		Rating:  rating,
		Comment: comment,
		Date:    date,
	}

	review.ID = utils.PString(uuid.NewV4().String())
	review.CreatedAt = utils.PTime(time.Now())

	if err := review.isValid(); err != nil {
		return nil, err
	}

	return review, nil
}

func (p *Review) isValid() error {
	_, err := govalidator.ValidateStruct(p)
	return err
}
