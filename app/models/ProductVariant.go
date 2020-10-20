package models

import (
	"database/sql"
	"fmt"
	"golang-api/app"
	"golang-api/app/structs"

	"github.com/google/uuid"
)

func FindProductVariant(productVariant structs.ProductVariant) (structs.ProductVariant, error) {

	db := app.GetDB()

	err := db.SelectOne(&productVariant, "SELECT id, product_id, name  FROM product_variants WHERE id=?", productVariant)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Find product error", err)
		}
		return productVariant, err
	}

	return productVariant, err
}

func StoreProductVariant(request structs.ProductVariantRequest, product structs.Product) (structs.ProductVariant, error) {

	db := app.GetDB()

	productVariant := structs.ProductVariant{
		ID:   uuid.New().String(),
		Name: request.Name,
	}

	_, err := db.Exec("INSERT INTO product_variants(id, product_id, name) VALUES (?,?,?)", productVariant.ID, product.ID, productVariant.Name)
	if err != nil {
		fmt.Println("Store product variant error", err)
		return productVariant, err
	}

	return productVariant, err
}
