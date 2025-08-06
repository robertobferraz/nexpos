package mercadopago

import (
	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
	"github.com/mercadopago/sdk-go/pkg/preference"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"mercadopago",
	ConfigModule,
	fx.Provide(NewMercadopago),
)

type Client struct {
	logger           *zap.SugaredLogger
	config           *config.Config
	preferenceClient *preference.Client
	paymentClient    *payment.Client
}

func NewMercadopago(c *Config, logger *zap.SugaredLogger) (*Client, error) {
	cfg, err := config.New(*c.AccessToken)
	if err != nil {
		logger.Errorf("failed to create Mercado Pago config: %v", err)
		return nil, err
	}

	preferenceClient := preference.NewClient(cfg)
	paymentClient := payment.NewClient(cfg)

	return &Client{
		logger:           logger,
		config:           cfg,
		preferenceClient: &preferenceClient,
		paymentClient:    &paymentClient,
	}, nil
}

func (m *Client) Client() *preference.Client {
	return m.preferenceClient
}
