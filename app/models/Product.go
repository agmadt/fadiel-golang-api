package models

import (
	"golang-api/app"
	"golang-api/app/structs"

	"github.com/gin-gonic/gin"
)

func FindProduct(c *gin.Context, product structs.Product) (structs.Product, error) {

	db := app.GetDB()

	err := db.SelectOne(&product, "SELECT id, name, price, description, created_at, updated_at FROM products WHERE id=?", product.ID)
	if err != nil {
		c.JSON(404, gin.H{
			"message": "Product not found",
		})

		return product, err
	}

	return product, err
}

func FindProductImages(c *gin.Context, product structs.Product) ([]structs.Image, error) {

	var images = []structs.Image{}

	db := app.GetDB()

	_, err := db.Select(&images, "SELECT id, product_id, image FROM product_images WHERE product_id=?", product.ID)
	if err != nil {
		c.JSON(404, gin.H{
			"message": "Images not found",
		})

		return images, err
	}

	return images, err
}

func FindProductCategories(c *gin.Context, product structs.Product) ([]structs.Category, error) {

	var categories = []structs.Category{}

	db := app.GetDB()

	_, err := db.Select(&categories, "SELECT categories.id, categories.name, categories.created_at, categories.updated_at FROM product_categories JOIN categories ON product_categories.category_id = categories.id WHERE product_categories.product_id=?", product.ID)
	if err != nil {
		c.JSON(404, gin.H{
			"message": "Categories not found",
		})

		return categories, err
	}

	return categories, err
}
