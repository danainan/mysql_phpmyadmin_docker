package main

import (
	"auth/connection"
	"auth/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	connection.InitMySQL()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))
	router.SetUpRoutes(app)

	app.Listen(":8000")
}
