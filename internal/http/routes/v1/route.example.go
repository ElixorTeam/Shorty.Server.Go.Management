package v1

import (
	V1Handlers "Shorty.Server.Go.Management/internal/http/handlers/v1"
	"Shorty.Server.Go.Management/internal/http/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type ExampleRoute struct {
	router         fiber.Router
	AuthMiddleware func(roles ...string) fiber.Handler
	ExampleHandler *V1Handlers.ExampleHandler
}

func NewExampleRoute(api fiber.Router, rdc *redis.Client) *ExampleRoute {
	AuthMiddleware := middlewares.KeycloakAuthMiddleware(rdc)
	ExampleHandler := V1Handlers.NewExampleHandler()

	return &ExampleRoute{
		router:         api,
		AuthMiddleware: AuthMiddleware,
		ExampleHandler: ExampleHandler,
	}
}

func (r *ExampleRoute) RegisterRoutes() {
	exampleRoutes := r.router.Group("/example")
	{
		exampleRoutes.Get("/test", r.AuthMiddleware(), r.ExampleHandler.Example)
	}
}
