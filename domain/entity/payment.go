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

type Payment struct {
	Base          `json:",inline" valid:"-"`
	OrderID       *string    `json:"-" valid:"-"`
	Order         *Order     `json:"order" valid:"-"`
	Amount        *float64   `json:"amount" valid:"-"`
	Date          *time.Time `json:"date" valid:"-"`
	Status        *string    `json:"status" valid:"-"`
	TransactionID *string    `json:"transaction_id" valid:"-"`
}

func NewPayment(
	order *Order,
	amount *float64,
	date *time.Time,
	status *string,
	transactionID *string,
) (*Payment, error) {
	payment := &Payment{
		OrderID:       order.ID,
		Order:         order,
		Amount:        amount,
		Date:          date,
		Status:        status,
		TransactionID: transactionID,
	}

	payment.ID = utils.PString(uuid.NewV4().String())
	payment.CreatedAt = utils.PTime(time.Now())

	if err := payment.isValid(); err != nil {
		return nil, err
	}

	return payment, nil
}

func (p *Payment) isValid() error {
	_, err := govalidator.ValidateStruct(p)
	return err
}
