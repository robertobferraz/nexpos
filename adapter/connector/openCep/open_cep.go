package openCep

import (
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"open_cep",
	ConfigModule,
	fx.Provide(NewOpenCep),
)

type OpenCep struct {
	client  *http.Client
	logger  *zap.SugaredLogger
	BaseUrl *string
}

func NewOpenCep(c *Config, l *zap.SugaredLogger) *OpenCep {
	client := &http.Client{}
	return &OpenCep{
		client:  client,
		logger:  l,
		BaseUrl: c.OpenCepBaseURL,
	}
}

func (c *OpenCep) Client() *http.Client {
	return c.client
}
