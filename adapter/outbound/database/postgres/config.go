package postgres

import (
	"fmt"

	"github.com/Netflix/go-env"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var ConfigModule = fx.Module(
	"postgres_config",
	fx.Provide(NewConfig),
)

type Config struct {
	User        *string `env:"POSTGRES_USER,required=true"`
	Password    *string `env:"POSTGRES_PASSWORD,default=''"`
	Host        *string `env:"POSTGRES_HOST,required=true"`
	Port        *int    `env:"POSTGRES_PORT,default=5432"`
	Database    *string `env:"POSTGRES_DB,required=true"`
	ZeroThrust  *bool   `env:"ZERO_THRUST,default=false"`
	AutoMigrate *bool   `env:"AUTO_MIGRATE,default=true"`
}

func (c Config) GetDsn(zeroThrust bool) string {
	if zeroThrust {
		return fmt.Sprintf(
			"user=%s dbname=%s sslmode=disable",
			*c.User, *c.Database,
		)
	}
	return fmt.Sprintf(
		"host=%s user=%s password=%s database=%s port=%d sslmode=disable",
		*c.Host, *c.User, *c.Password, *c.Database, *c.Port,
	)
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
