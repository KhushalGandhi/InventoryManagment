package handlers

import (
	"bytes"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a test token
func createTestToken(role string) string {
	claims := jwt.MapClaims{
		"role": role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte("secret")) // Use the same secret key as in your middleware
	return signedToken
}

// Test UploadProduct function
func TestUploadProduct(t *testing.T) {
	app := fiber.New()

	app.Post("/api/admin/product/upload", UploadProduct) // Assuming you have this handler

	// Prepare a mock product payload
	payload := `{"name": "Test Product", "price": 10.99}`

	// Test case: successful upload
	req := httptest.NewRequest(http.MethodPost, "/api/admin/product/upload", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+createTestToken("admin")) // Mock admin token
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test case: unauthorized upload
	req = httptest.NewRequest(http.MethodPost, "/api/admin/product/upload", bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+createTestToken("user")) // Mock user token
	resp, _ = app.Test(req)

	assert.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Test case: bad request (invalid JSON)
	req = httptest.NewRequest(http.MethodPost, "/api/admin/product/upload", bytes.NewBuffer([]byte(`{"invalid"`)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+createTestToken("admin")) // Mock admin token
	resp, _ = app.Test(req)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
