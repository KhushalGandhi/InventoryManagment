package handlers

import (
	"github.com/gofiber/fiber/v2"
	"inventory-management/database"
	"inventory-management/models"
	"inventory-management/s3"
)

// UploadProduct allows a user to upload a product
func UploadProduct(c *fiber.Ctx) error {
	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	userID := c.Locals("userID").(uint)
	product.UserID = userID
	product.Status = "pending"

	// Upload image to S3
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No file uploaded"})
	}

	// Upload to S3 and get the image URL
	imageURL, err := s3.UploadToS3(file, "products")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload image"})
	}

	// Set the uploaded image URL to the product
	product.ImageURL = []string{imageURL}

	// Save the product to the database
	result := database.DB.Create(&product)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not upload product"})
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

// ListProducts shows all the products uploaded by a user
func ListProducts(c *fiber.Ctx) error {
	var products []models.Product
	database.DB.Find(&products)

	return c.JSON(products)
}

// AdminApproveProduct allows an admin to approve a product
func AdminApproveProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product

	result := database.DB.First(&product, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	product.Status = "approved"
	database.DB.Save(&product)

	return c.JSON(product)
}

// AdminRejectProduct allows an admin to reject a product
func AdminRejectProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product

	result := database.DB.First(&product, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	product.Status = "rejected"
	database.DB.Save(&product)

	return c.JSON(product)
}
