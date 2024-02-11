package v1

import "github.com/gofiber/fiber/v2"

type ExampleHandler struct {
}

func NewExampleHandler() *ExampleHandler {
	return &ExampleHandler{}
}

func (h *ExampleHandler) Example(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"message": "example",
	})
}
