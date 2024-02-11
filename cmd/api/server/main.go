package server

import (
	"Shorty.Server.Go.Mangment/internal/config"
	"Shorty.Server.Go.Mangment/internal/constants"
	V1Routes "Shorty.Server.Go.Mangment/internal/http/routes/v1"
	"Shorty.Server.Go.Mangment/internal/utils"
	"Shorty.Server.Go.Mangment/pkg/logger"
	"fmt"
	"github.com/gofiber/fiber/v2"
	CorsMiddleware "github.com/gofiber/fiber/v2/middleware/cors"
	HelmetMiddleware "github.com/gofiber/fiber/v2/middleware/helmet"
	LoggerMiddleware "github.com/gofiber/fiber/v2/middleware/logger"
	MonitorMiddleware "github.com/gofiber/fiber/v2/middleware/monitor"
	RecoverMiddleware "github.com/gofiber/fiber/v2/middleware/recover"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	App *fiber.App
}

func NewApp() (*Server, error) {
	redisClient := utils.SetupRedisConn()

	app := setupRouter()
	V1Routes.NewExampleRoute(app, redisClient).RegisterRoutes()

	return &Server{app}, nil
}

func setupRouter() *fiber.App {

	router := fiber.New(fiber.Config{
		AppName:     "Shorty Manage Server",
		ReadTimeout: 60 * time.Second,
	})

	router.Use(LoggerMiddleware.New(LoggerMiddleware.Config{
		Format:        "[${ip}] [${time}] ${status} - ${method} ${path}\n",
		DisableColors: false,
		TimeFormat:    "02-01-2006/15:04:05",
		TimeZone:      "Europe/Moscow",
	}))
	router.Use(HelmetMiddleware.New())

	router.Use(CorsMiddleware.New(CorsMiddleware.Config{
		AllowOrigins: constants.AllowOrigin,
		AllowHeaders: constants.AllowHeader,
		AllowMethods: constants.AllowMethods,
	}))

	router.Use(RecoverMiddleware.New())

	router.Get("/monitor", MonitorMiddleware.New(MonitorMiddleware.Config{
		Title: "Shorty Manage Server",
	}))

	return router

}

func (a *Server) Start() error {
	go func() {
		addr := fmt.Sprintf(":%d", config.AppConfig.Port)
		if err := a.App.Listen(addr); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	logger.Info("Shutting down server...", nil)

	if err := a.App.Shutdown(); err != nil {
		logger.Error(err.Error(), nil)
	}

	logger.Info("Server exiting", nil)

	return nil

}
