package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/robertobff/food-service/adapter/inbound/http/handler"
	"go.uber.org/fx"
)

var HandlerModule = fx.Module(
	"handler",
	fx.Invoke(HandleRoutes),
	handler.SwaggerHandlerModule,
	handler.AuthHandlerModule,
)

func HandleRoutes(
	http *Http,
	swaggerHandler *handler.SwaggerHandler,
	authHandler *handler.AuthHandler,

) {
	http.App.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/v1/swagger/index.html")
	})

	v1 := http.App.Group("/v1")
	swaggerHandler.RegisterRoutes(v1)
	authHandler.RegisterRoutes(v1)
}
