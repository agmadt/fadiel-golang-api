package models

import (
	"fmt"
	"golang-api/app"
	"golang-api/app/structs"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func FindProductVariantOption(c *gin.Context, productVariantOption structs.ProductVariantOption) (structs.ProductVariantOption, error) {

	db := app.GetDB()

	err := db.SelectOne(&productVariantOption, "SELECT id, product_variant_id, name  FROM product_variant_options WHERE id=?", productVariantOption)

	if err != nil {
		c.JSON(404, gin.H{
			"message": "Product variant not found",
		})

		return productVariantOption, err
	}

	return productVariantOption, err
}

func StoreProductVariantOption(c *gin.Context, request structs.ProductVariantOptionRequest, productVariant structs.ProductVariant) (structs.ProductVariantOption, error) {

	db := app.GetDB()

	productVariantOption := structs.ProductVariantOption{
		ID:   uuid.New().String(),
		Name: request.Name,
	}

	_, err := db.Exec("INSERT INTO product_variant_options(id, product_variant_id, name) VALUES (?,?,?)", productVariantOption.ID, productVariant.ID, productVariantOption.Name)

	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error while storing data",
		})
		fmt.Println(err)

		return productVariantOption, err
	}

	return productVariantOption, err
}
