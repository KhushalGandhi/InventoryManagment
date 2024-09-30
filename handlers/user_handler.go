package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"inventory-management/database"
	"inventory-management/models"
	"os"
	"time"
)

// RegisterUser registers a new user
// User Registration Handler
func RegisterUser(c *fiber.Ctx) error {
	var body models.User

	// Parse request body
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Default role to "user" if not provided
	if body.Role == "" {
		body.Role = "user"
	}

	// Validate if role is either "user" or "admin"
	if body.Role != "user" && body.Role != "admin" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid role. Allowed values are 'user' or 'admin'",
		})
	}

	// Encrypt password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not encrypt password",
		})
	}

	body.Password = string(hashedPassword)

	// Save the user to the database
	if err := database.DB.Create(&body).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create user",
		})
	}

	// Return created user (without the password)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":    body.ID,
		"name":  body.Name,
		"email": body.Email,
		"role":  body.Role,
	})
}

// LoginUser logs a user in
func LoginUser(c *fiber.Ctx) error {
	var loginData models.User
	if err := c.BodyParser(&loginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	var user models.User
	database.DB.Where("user_name = ?", loginData.UserName).First(&user)

	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["username"] = user.UserName
	claims["role"] = user.Role
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	return c.JSON(fiber.Map{"token": tokenString})
}
