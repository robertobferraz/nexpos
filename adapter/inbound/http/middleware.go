package http

import (
	"github.com/robertobff/nexpos/adapter/inbound/http/middleware"
	"go.uber.org/fx"
)

var MiddlewareModule = fx.Module(
	"middleware",
	middleware.AuthMiddlewareModule,
	middleware.UserMiddlewareModule,
)
