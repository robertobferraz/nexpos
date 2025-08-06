package http

import (
	"github.com/robertobff/food-service/adapter/inbound/http/middleware"
	"go.uber.org/fx"
)

var MiddlewareModule = fx.Module(
	"middleware",
	middleware.UserMiddlewareModule,
)
