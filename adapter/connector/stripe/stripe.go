package stripe

import (
	"github.com/stripe/stripe-go/v76/client"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"stripe",
	ConfigModule,
	fx.Provide(NewStripe),
)

type Stripe struct {
	client    *client.API
	logger    *zap.SugaredLogger
	SecretKey *string
}

func NewStripe(c *Config, logger *zap.SugaredLogger) *Stripe {
	sc := client.API{}
	sc.Init(*c.SecretKey, nil)
	return &Stripe{
		client:    &sc,
		logger:    logger,
		SecretKey: c.SecretKey,
	}
}

func (s *Stripe) Client() *client.API {
	return s.client
}
