package stripe

import (
	"github.com/Netflix/go-env"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var ConfigModule = fx.Module(
	"stripe_conn_config",
	fx.Provide(NewConfig),
)

type Config struct {
	SecretKey *string `env:"STRIPE_SECRET_KEY,required=true"`
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
