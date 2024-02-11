package v1

import (
	V1Domains "Shorty.Server.Go.Mangment/internal/business/domains/v1"
	"Shorty.Server.Go.Mangment/internal/constants"
	"github.com/gofiber/fiber/v2"
)

type ExampleHandler struct {
}

func NewExampleHandler() *ExampleHandler {
	return &ExampleHandler{}
}

func (h *ExampleHandler) Example(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"message": ctx.Locals(constants.UserDetails).(V1Domains.UserDetails).Email,
	})
}
