package middlewares

import (
	"github.com/gofiber/jwt/v3"
	"os"
)

func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	})
}