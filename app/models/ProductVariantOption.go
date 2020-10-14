package models

import (
	"golang-api/app"
	"golang-api/app/structs"

	"github.com/gin-gonic/gin"
)

func FindProductVariantOption(c *gin.Context, pk string) (structs.ProductVariantOption, error) {

	var productVariantOption structs.ProductVariantOption

	db := app.GetDB()

	err := db.SelectOne(&productVariantOption, "SELECT id, product_variant_id, name  FROM product_variant_options WHERE id=?", pk)

	if err != nil {
		c.JSON(404, gin.H{
			"message": "Product variant not found",
		})

		return productVariantOption, err
	}

	return productVariantOption, err
}
