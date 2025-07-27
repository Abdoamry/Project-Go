package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"time"
	"os"
	"github.com/abdoamry/Project-go/database"
	"github.com/abdoamry/Project-go/models"
	"golang.org/x/crypto/bcrypt"

)

func Register(c *fiber.Ctx) error {
	u := new(models.User)
	if err := c.BodyParser(u); err != nil {
		return err
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	u.Password = string(hashedPassword)
	database.DB.Create(&u)
	return c.JSON(u)
}

func Login(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var input LoginInput
	c.BodyParser(&input)

	var user models.User
	database.DB.Where("email = ?", input.Email).First(&user)
	if user.ID == 0 {
		return c.Status(401).SendString("Invalid email or password")
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return c.Status(401).SendString("Invalid email or password")
	}

	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	signedToken, _ := token.SignedString([]byte(secret))

	// تخزين في Redis كمثال
	database.Redis.Set(database.Ctx, signedToken, user.Email, time.Hour*72)

	return c.JSON(fiber.Map{"token": signedToken})
}