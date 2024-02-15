package v1

import (
	"Shorty.Server.Go.Management/internal/constants"
	"github.com/gofiber/fiber/v2"
)

type ExampleHandler struct {
}

func NewExampleHandler() *ExampleHandler {
	return &ExampleHandler{}
}

func (h *ExampleHandler) Example(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"user_details": ctx.Locals(constants.UserDetails),
		"message":      "Hello, World!",
	})
}
