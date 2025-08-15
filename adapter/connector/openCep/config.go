package openCep

import (
	"github.com/Netflix/go-env"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Config struct {
	OpenCepBaseURL *string `json:"open_cep_base_url" env:"OPEN_CEP_BASE_URL"`
}

var ConfigModule = fx.Module(
	"open_cep_config",
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
