package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/robertobff/nexpos/adapter/inbound/http/handler"
	"go.uber.org/fx"
)

var HandlerModule = fx.Module(
	"handler",
	fx.Invoke(HandleRoutes),
	handler.SwaggerHandlerModule,
	handler.AuthHandlerModule,
	handler.CategoryHandlerModule,
	handler.UserHandlerModule,
	handler.ImageHandlerModule,
)

func HandleRoutes(
	http *Http,
	swaggerHandler *handler.SwaggerHandler,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	imageHandler *handler.ImageHandler,
	categoryHandler *handler.CategoryHandler,
) {
	http.App.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/v1/swagger/index.html")
	})

	v1 := http.App.Group("/v1")
	swaggerHandler.RegisterRoutes(v1)
	authHandler.RegisterRoutes(v1)
	userHandler.RegisterRoutes(v1)
	imageHandler.RegisterRoutes(v1)
	categoryHandler.RegisterRoutes(v1)
}
