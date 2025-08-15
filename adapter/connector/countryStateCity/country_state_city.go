package countryStateCity

import (
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"country_state_city",
	ConfigModule,
	fx.Provide(NewCountryStateCity),
)

type CountryStateCity struct {
	client *http.Client
	logger *zap.SugaredLogger
	ApiKey *string
}

func NewCountryStateCity(c *Config, l *zap.SugaredLogger) *CountryStateCity {
	client := &http.Client{}
	return &CountryStateCity{
		client: client,
		logger: l,
		ApiKey: c.CountryStateCityAPIKey,
	}
}

func (c *CountryStateCity) Client() *http.Client {
	return c.client
}
