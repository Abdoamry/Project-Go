package routes

import (
	"github.com/abdoamry/Project-go/controllers"
	middlewares "github.com/abdoamry/Project-go/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)

	protected := app.Group("/api")
	protected.Use(middlewares.Protected())
	protected.Get("/profile", func(c *fiber.Ctx) error {
		return c.SendString("Welcome Bro ")
	})
}
