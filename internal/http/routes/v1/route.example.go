package v1

import (
	V1Handlers "Shorty.Server.Go.Mangment/internal/http/handlers/v1"
	"Shorty.Server.Go.Mangment/internal/http/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type ExampleRoute struct {
	router         fiber.Router
	AuthMiddleware func(roles ...string) fiber.Handler
	V1Handlers.ExampleHandler
}

func NewExampleRoute(api fiber.Router, rdc *redis.Client) *ExampleRoute {
	AuthMiddleware := middlewares.KeycloakAuthMiddleware(rdc)

	return &ExampleRoute{
		router:         api,
		AuthMiddleware: AuthMiddleware,
	}
}

func (r *ExampleRoute) RegisterRoutes() {
	exampleRoutes := r.router.Group("/example")
	{
		exampleRoutes.Get("/test", r.AuthMiddleware(), r.ExampleHandler.Example)
	}
}
