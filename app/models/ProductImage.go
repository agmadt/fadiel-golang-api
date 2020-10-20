package models

import (
	"fmt"
	"golang-api/app"
	"golang-api/app/structs"

	"github.com/google/uuid"
)

func StoreProductImage(request structs.ProductImageRequest, product structs.Product) (structs.ProductImage, error) {

	db := app.GetDB()

	productImage := structs.ProductImage{
		ID:    uuid.New().String(),
		Image: request.Image,
	}

	_, err := db.Exec("INSERT INTO product_images(id, product_id, image) VALUES (?,?,?)", productImage.ID, product.ID, productImage.Image)
	if err != nil {
		fmt.Println("Store product image error", err)
		return productImage, err
	}

	return productImage, err
}
