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

type Image struct {
	Base        `json:",inline" valid:"-"`
	Name        *string `json:"name" valid:"-"`
	Data        *[]byte `json:"-" valid:"-"`
	ContentType *string `json:"content_type" valid:"-"`
}

func NewImage(name, contentType *string, data *[]byte) (*Image, error) {
	image := &Image{
		Name:        name,
		Data:        data,
		ContentType: contentType,
	}

	image.ID = utils.PString(uuid.NewV4().String())
	image.CreatedAt = utils.PTime(time.Now())

	if err := image.isValid(); err != nil {
		return nil, err
	}

	return image, nil
}

func (p *Image) isValid() error {
	_, err := govalidator.ValidateStruct(p)
	return err
}
