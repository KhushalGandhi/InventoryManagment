package main

import (
	"github.com/gofiber/fiber/v2"
	"inventory-management/database"
	"inventory-management/routes"
)

func main() {
	// Initialize Fiber app
	app := fiber.New()

	// Connect to the database
	db := database.ConnectDB()

	// Run migrations
	err := db.RunMigrations()
	if err != nil {
		panic("Migrations failed")
	}

	// Set up routes
	routes.SetupRoutes(app)

	// Start the server
	app.Listen(":3000")
}
