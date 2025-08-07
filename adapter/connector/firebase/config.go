package firebase

import (
	"github.com/Netflix/go-env"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var ConfigModule = fx.Module(
	"firebase_conn_config",
	fx.Provide(NewConfig),
)

type Config struct {
	Type                    *string `json:"type" env:"TYPE"`
	ProjectId               *string `json:"project_id" env:"PROJECT_ID"`
	PrivateKeyId            *string `json:"private_key_id" env:"PRIVATE_KEY_ID"`
	PrivateKey              *string `json:"private_key" env:"PRIVATE_KEY"`
	ClientEmail             *string `json:"client_email" env:"CLIENT_EMAIL"`
	ClientId                *string `json:"client_id" env:"CLIENT_ID"`
	AuthUri                 *string `json:"auth_uri" env:"AUTH_URI"`
	TokenUri                *string `json:"token_uri" env:"TOKEN_URI"`
	AuthProviderX509CertUrl *string `json:"auth_provider_x509_cert_url" env:"AUTH_PROVIDER_X509_CERT_URL"`
	ClientX509CertUrl       *string `json:"client_x509_cert_url" env:"CLIENT_X509_CERT_URL"`
	UniverseDomain          *string `json:"universe_domain" env:"UNIVERSE_DOMAIN"`
	Env                     *string `env:"ENV" required:"true"`
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
