package models

import (
	"fmt"
	"golang-api/app"
	"golang-api/app/structs"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func FindProductVariantOptions(c *gin.Context, productVariant structs.ProductVariant) ([]structs.ProductVariantOption, error) {

	var variantOptions = []structs.ProductVariantOption{}

	db := app.GetDB()

	_, err := db.Select(&variantOptions, "SELECT id, product_variant_id, name FROM product_variant_options WHERE product_variant_id=?", productVariant.ID)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error while getting product variants",
		})
		fmt.Println(err)

		return variantOptions, err
	}

	return variantOptions, err
}

func StoreProductVariant(c *gin.Context, request structs.ProductVariantRequest, product structs.Product) (structs.ProductVariant, error) {

	db := app.GetDB()

	productVariant := structs.ProductVariant{
		ID:   uuid.New().String(),
		Name: request.Name,
	}

	_, err := db.Exec("INSERT INTO product_variants(id, product_id, name) VALUES (?,?,?)", productVariant.ID, product.ID, productVariant.Name)

	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error while storing data",
		})
		fmt.Println(err)

		return productVariant, err
	}

	return productVariant, err
}
