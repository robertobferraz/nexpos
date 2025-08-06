package stripe

import (
	food "github.com/robertobff/food-service/adapter/connector/stripe"
	"github.com/robertobff/food-service/domain/entity"
	"github.com/stripe/stripe-go/v76"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"stripe_api",
	ConfigModule,
	fx.Provide(NewStripeAPI),
)

type StripeAPI struct {
	stripe         *food.Stripe
	logger         *zap.SugaredLogger
	EndpointSecret *string
	SecretKey      *string
	RedirectURL    *string
}

func NewStripeAPI(
	stripe *food.Stripe,
	c *Config,
	logger *zap.SugaredLogger,
) *StripeAPI {
	return &StripeAPI{
		stripe:      stripe,
		logger:      logger,
		RedirectURL: c.RedirectURL,
	}
}

func (s *StripeAPI) CreateCustomer(u *entity.User) (*stripe.Customer, error) {
	customer, err := s.stripe.Client().Customers.New(&stripe.CustomerParams{
		Email: u.Email,
		Name:  u.Name,
		Metadata: map[string]string{
			"user_id": *u.ID,
		},
	})
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *StripeAPI) DeletePlanCheckout(externalID *string) (*stripe.CheckoutSession, error) {
	return s.stripe.Client().CheckoutSessions.Expire(*externalID, &stripe.CheckoutSessionExpireParams{})
}
