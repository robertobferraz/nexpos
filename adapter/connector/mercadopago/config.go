package mercadopago

import (
	"github.com/Netflix/go-env"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var ConfigModule = fx.Module(
	"mercadopago_conn_config",
	fx.Provide(NewConfig),
)

type Config struct {
	AccessToken *string `env:"MERCADOPAGO_ACCESS_TOKEN,required=true"`
}

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
