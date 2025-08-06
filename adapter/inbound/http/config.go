package http

import (
	"github.com/Netflix/go-env"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var ConfigModule = fx.Module(
	"http_config",
	fx.Provide(NewConfig),
)

type Config struct {
	Port                  *int  `env:"HTTP_PORT,required=true"`
	DisableStartupMessage *bool `env:"HTTP_DISABLE_STARTUP_MESSAGE,default=true" `
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
