package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yourname/fiber-jwt-app/controllers"
	"github.com/yourname/fiber-jwt-app/middlewares"
)

func SetupRoutes(app *fiber.App) {
	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)

	protected := app.Group("/api")
	protected.Use(middlewares.Protected())
	protected.Get("/profile", func(c *fiber.Ctx) error {
		return c.SendString("مرحبا بالمستخدم المحمي")
	})
}
