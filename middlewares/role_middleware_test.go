package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a test token
func CreateTestToken(role string) string {
	claims := jwt.MapClaims{
		"role": role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte("secret")) // Use the same secret key as in your middleware
	return signedToken
}

// Test AdminOnly middleware
func TestAdminOnly(t *testing.T) {
	app := fiber.New()

	app.Use(AdminOnly) // Apply the AdminOnly middleware

	app.Get("/admin", func(c *fiber.Ctx) error {
		return c.SendString("Admin access granted")
	})

	// Test case: valid admin token
	req := httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+CreateTestToken("admin"))
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test case: invalid user token
	req = httptest.NewRequest(http.MethodGet, "/admin", nil)
	req.Header.Set("Authorization", "Bearer "+CreateTestToken("user"))
	resp, _ = app.Test(req)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

// Test UserOnly middleware
func TestUserOnly(t *testing.T) {
	app := fiber.New()

	app.Use(UserOnly) // Apply the UserOnly middleware

	app.Get("/user", func(c *fiber.Ctx) error {
		return c.SendString("User access granted")
	})

	// Test case: valid user token
	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	req.Header.Set("Authorization", "Bearer "+CreateTestToken("user"))
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test case: invalid admin token
	req = httptest.NewRequest(http.MethodGet, "/user", nil)
	req.Header.Set("Authorization", "Bearer "+CreateTestToken("admin"))
	resp, _ = app.Test(req)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}
