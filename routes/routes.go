package routes

import (
	"github.com/gofiber/fiber/v2"
	"inventory-management/handlers"
	"inventory-management/middlewares"
)

func SetupRoutes(app *fiber.App) {
	// User routes
	app.Post("/api/register", handlers.RegisterUser)
	app.Post("/api/login", handlers.LoginUser)

	// JWT-protected routes
	app.Use("/api", middlewares.JWTProtected())

	// Product routes
	app.Post("/api/product", handlers.UploadProduct)
	app.Get("/api/products", handlers.ListProducts)
	app.Post("/upload", handlers.UploadImageHandler)

	// Admin routes (Admin only)
	app.Patch("/api/admin/product/:id/approve", middlewares.AdminOnly, handlers.AdminApproveProduct)
	app.Patch("/api/admin/product/:id/reject", middlewares.AdminOnly, handlers.AdminRejectProduct)
}
