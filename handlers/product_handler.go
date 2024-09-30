package handlers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"inventory-management/database"
	"inventory-management/models"
	"mime/multipart"
	"path/filepath"
)

// UploadProduct allows a user to upload a product
func UploadProduct(c *fiber.Ctx) error {
	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	token := c.Locals("user").(*jwt.Token)     // Assert the user to *jwt.Token
	claims, ok := token.Claims.(jwt.MapClaims) // Extract claims as jwt.MapClaims
	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized, invalid token",
		})
	}

	// Get the user ID from claims
	userIDFloat, exists := claims["id"].(float64) // assuming your claims have an "id" key
	if !exists {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized, user ID not found in token",
		})
	}

	// Convert userID to uint
	product.UserID = uint(userIDFloat)
	product.Status = "pending"

	// Upload image to S3
	//file, err := c.FormFile("image")
	//if err != nil {
	//	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No file uploaded"})
	//}
	//
	//// Upload to S3 and get the image URL
	//imageURL, err := s3.UploadToS3(file, "products")
	//if err != nil {
	//	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to upload image"})
	//}
	//
	//// Set the uploaded image URL to the product
	//product.ImageURL = []string{imageURL}

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

func UploadFileToS3(file multipart.File, fileHeader *multipart.FileHeader, bucketName, keyPrefix string) (string, error) {
	// Load the AWS SDK config
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-west-2")) // Set your region
	if err != nil {
		return "", fmt.Errorf("unable to load SDK config, %v", err)
	}

	// Create an S3 service client
	s3Client := s3.NewFromConfig(cfg)

	// Read the file content into a buffer
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(file)
	if err != nil {
		return "", fmt.Errorf("could not read file: %v", err)
	}

	// Create the S3 key (filename)
	fileName := keyPrefix + filepath.Base(fileHeader.Filename)

	// Upload input parameters
	putObjectInput := &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
		ACL:         types.ObjectCannedACLPublicRead, // Set the ACL as public read
	}

	// Upload the file
	_, err = s3Client.PutObject(context.TODO(), putObjectInput)
	if err != nil {
		return "", fmt.Errorf("failed to upload file, %v", err)
	}

	// Return the file URL
	fileURL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, fileName)
	return fileURL, nil
}

// Handler for uploading an image to S3
func UploadImageHandler(c *fiber.Ctx) error {
	// Get the file from the form
	fileHeader, err := c.FormFile("image")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Could not get uploaded file")
	}

	// Open the file to get a multipart.File
	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Could not open uploaded file")
	}
	defer file.Close() // Ensure the file is closed after reading

	// Specify the S3 bucket and key prefix
	bucketName := "your-s3-bucket-name" // Replace with your S3 bucket name
	keyPrefix := "uploads/images/"

	// Upload the file to S3
	fileURL, err := UploadFileToS3(file, fileHeader, bucketName, keyPrefix)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to upload file to S3")
	}

	// Return the file URL as a response
	return c.JSON(fiber.Map{
		"message": "File uploaded successfully",
		"url":     fileURL,
	})
}
