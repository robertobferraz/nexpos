package countryStateCity

import (
	"github.com/Netflix/go-env"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Config struct {
	CountryStateCityBaseURL *string `json:"country_state_city_base_url" env:"COUNTRY_STATE_CITY_BASE_URL"`
}

var ConfigModule = fx.Module(
	"country_state_api_config",
	fx.Provide(NewConfig),
)

func NewConfig(l *zap.SugaredLogger) (*Config, error) {
	var cfg Config
	err := cfg.loadConfig()
	if err != nil {
		l.Error(err)
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) loadConfig() error {
	_, err := env.UnmarshalFromEnviron(c)
	if err != nil {
		return err
	}

	return nil
}
