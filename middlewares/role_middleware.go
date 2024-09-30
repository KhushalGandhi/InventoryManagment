package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4" // Import JWT package
)

func AdminOnly(c *fiber.Ctx) error {
	// Extract the user token from the context
	token := c.Locals("user")
	if token == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized, token not found",
		})
	}

	// Assert the token to jwt.Token type
	jwtToken, ok := token.(*jwt.Token)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized, invalid token",
		})
	}

	// Extract claims from the token
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized, invalid token claims",
		})
	}

	// Extract the user role
	role, exists := claims["role"].(string)
	if !exists || role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden, only admin users are allowed",
		})
	}

	// Call the next handler if the role is "admin"
	return c.Next()
}

func UserOnly(c *fiber.Ctx) error {
	// Extract the user token from the context
	token := c.Locals("user")
	if token == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized, token not found",
		})
	}

	// Assert the token to jwt.Token type
	jwtToken, ok := token.(*jwt.Token)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized, invalid token",
		})
	}

	// Extract claims from the token
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized, invalid token claims",
		})
	}

	// Extract the user role
	role, exists := claims["role"].(string)
	if !exists || role != "user" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden, only users with 'user' role are allowed",
		})
	}

	// Call the next handler if the role is "user"
	return c.Next()
}
