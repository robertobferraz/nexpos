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
	Type                    *string `json:"type" env:"FIREBASE_TYPE"`
	ProjectId               *string `json:"project_id" env:"FIREBASE_PROJECT_ID"`
	PrivateKeyId            *string `json:"private_key_id" env:"FIREBASE_PRIVATE_KEY_ID"`
	PrivateKey              *string `json:"private_key" env:"FIREBASE_PRIVATE_KEY"`
	ClientEmail             *string `json:"client_email" env:"FIREBASE_CLIENT_EMAIL"`
	ClientId                *string `json:"client_id" env:"FIREBASE_CLIENT_ID"`
	AuthUri                 *string `json:"auth_uri" env:"FIREBASE_AUTH_URI"`
	TokenUri                *string `json:"token_uri" env:"FIREBASE_TOKEN_URI"`
	AuthProviderX509CertUrl *string `json:"auth_provider_x509_cert_url" env:"FIREBASE_AUTH_PROVIDER_X509_CERT_URL"`
	ClientX509CertUrl       *string `json:"client_x509_cert_url" env:"FIREBASE_CLIENT_X509_CERT_URL"`
	UniverseDomain          *string `json:"universe_domain" env:"FIREBASE_UNIVERSE_DOMAIN"`
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
