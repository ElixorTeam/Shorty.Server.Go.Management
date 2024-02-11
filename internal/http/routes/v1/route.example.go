package v1

import (
	"Shorty.Server.Go.Mangment/internal/http/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type ExampleRoute struct {
	router         fiber.Router
	AuthMiddleware func(roles ...string) fiber.Handler
}

func NewExampleRoute(api fiber.Router, rdc *redis.Client) *ExampleRoute {
	KeycloaAuthMiddleware := middlewares.KeycloakAuthMiddleware(rdc)

	return &ExampleRoute{
		router:         api,
		AuthMiddleware: KeycloaAuthMiddleware,
	}
}

func (r *ExampleRoute) RegisterRoutes() {
	exampleRoutes := r.router.Group("/example")
	{
		exampleRoutes.Get("/test", r.AuthMiddleware(), func(ctx *fiber.Ctx) error {
			return ctx.JSON(fiber.Map{
				"message": "example",
			})
		})
	}
}
