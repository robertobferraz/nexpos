package handler

import (
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var SwaggerHandlerModule = fx.Module(
	"swagger_handler",
	fx.Provide(NewSwaggerHandler),
)

type SwaggerHandler struct {
	logger *zap.SugaredLogger
}

func NewSwaggerHandler(
	logger *zap.SugaredLogger,
) (*SwaggerHandler, error) {
	return &SwaggerHandler{logger: logger}, nil
}

func (h *SwaggerHandler) RegisterRoutes(r fiber.Router) {
	r.Get("/swagger/*", fiberSwagger.WrapHandler)
}
