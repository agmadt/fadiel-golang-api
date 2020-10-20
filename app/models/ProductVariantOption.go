package models

import (
	"database/sql"
	"fmt"
	"golang-api/app"
	"golang-api/app/structs"

	"github.com/google/uuid"
)

func FindProductVariantOption(productVariantOption structs.ProductVariantOption) (structs.ProductVariantOption, error) {

	db := app.GetDB()

	err := db.SelectOne(&productVariantOption, "SELECT id, product_variant_id, name  FROM product_variant_options WHERE id=?", productVariantOption)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Find product variant option error", err)
		}
		return productVariantOption, err
	}

	return productVariantOption, err
}

func StoreProductVariantOption(request structs.ProductVariantOptionRequest, productVariant structs.ProductVariant) (structs.ProductVariantOption, error) {

	db := app.GetDB()

	productVariantOption := structs.ProductVariantOption{
		ID:   uuid.New().String(),
		Name: request.Name,
	}

	_, err := db.Exec("INSERT INTO product_variant_options(id, product_variant_id, name) VALUES (?,?,?)", productVariantOption.ID, productVariant.ID, productVariantOption.Name)
	if err != nil {
		fmt.Println("Store product variant option error", err)
		return productVariantOption, err
	}

	return productVariantOption, err
}
