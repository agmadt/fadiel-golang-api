package models

import (
	"golang-api/app"
	"golang-api/app/structs"

	"github.com/gin-gonic/gin"
)

func FindProductVariant(c *gin.Context, productVariant structs.ProductVariant) (structs.ProductVariant, error) {

	db := app.GetDB()

	err := db.SelectOne(&productVariant, "SELECT id, product_id, name  FROM product_variants WHERE id=?", productVariant)

	if err != nil {
		c.JSON(404, gin.H{
			"message": "Product variant not found",
		})

		return productVariant, err
	}

	return productVariant, err
}
