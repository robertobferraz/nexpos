package redis

import (
	"fmt"

	"github.com/Netflix/go-env"
	"github.com/robertobff/nexpos/utils"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var ConfigModule = fx.Module(
	"redis_config",
	fx.Provide(NewConfig),
)

type Config struct {
	Addr     *string `env:"REDIS_HOST" required:"true"`
	Password *string `env:"REDIS_PASSWORD" required:"true"`
	DB       *int    `env:"REDIS_DATABASE" required:"true"`
	Port     *string `env:"REDIS_ADDRESS_PORT" required:"true"`
}

func (c *Config) ConnectionString() *string {
	return utils.PString(fmt.Sprintf(
		"%s:%s",
		*c.Addr, *c.Port,
	))
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
