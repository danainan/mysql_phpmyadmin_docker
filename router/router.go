package router

import (
	"auth/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetUpRoutes(app *fiber.App) {
	api := app.Group("api", logger.New())
	api.Get("/", handler.Hello)
	api.Post("/register", handler.Register)
	api.Post("/login", handler.Login)
	api.Get("/users", handler.GetAllUser)
}
