package utils

import (
	"Shorty.Server.Go.Management/internal/config"
	"Shorty.Server.Go.Management/pkg/logger"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func CustomErrorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		fmt.Println(e.Message + " Ya error by fiber")
		code = e.Code
	}

	if config.AppConfig.Debug {
		logger.Error("Error", logrus.Fields{"err": err})
		return ctx.Status(code).JSON(fiber.Map{"message": err.Error()})
	}

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).Send(nil)
	}

	// Return from handler
	return ctx.Next()
}
