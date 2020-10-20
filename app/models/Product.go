package models

import (
	"database/sql"
	"fmt"
	"golang-api/app"
	"golang-api/app/structs"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func PaginateProduct(paginationParam structs.PaginationParameters) (structs.ProductPagination, error) {

	var products []structs.Product
	var totalResponse int
	var productResponses []structs.ProductResponse
	var productPagination structs.ProductPagination
	var offset string = "0"
	var sqlQuery string
	var limit string = strconv.Itoa(paginationParam.Limit)

	offset = strconv.Itoa(((paginationParam.Page) - 1) * paginationParam.Limit)

	db := app.GetDB()

	sqlQuery = "SELECT id, name, price, description FROM products ORDER BY created_at DESC"
	sqlQuery += " LIMIT " + offset + ", " + limit

	_, err := db.Select(&products, sqlQuery)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("Query to get all product error", err)
		return productPagination, err
	}

	err = db.SelectOne(&totalResponse, "SELECT count(id) as total FROM products")
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("Query to count all product error", err)
		return productPagination, err
	}

	for _, product := range products {
		productResponses = append(productResponses, product.Response())
	}

	productPagination = structs.ProductPagination{
		Products: productResponses,
		Total:    totalResponse,
		Limit:    paginationParam.Limit,
		Page:     paginationParam.Page,
	}

	return productPagination, nil
}

func FindProduct(product structs.Product) (structs.Product, error) {

	db := app.GetDB()

	err := db.SelectOne(&product, "SELECT id, name, price, description, created_at, updated_at FROM products WHERE id=?", product.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Find product error", err)
		}
		return product, err
	}

	return product, nil
}

func DeleteProduct(product structs.Product) error {

	db := app.GetDB()

	_, err := db.Exec("DELETE FROM products WHERE id=?", product.ID)
	if err != nil {
		fmt.Println("Delete product error", err)
		return err
	}

	return nil
}

func FindProductImages(product structs.Product) ([]structs.ProductImage, error) {

	var images = []structs.ProductImage{}

	db := app.GetDB()

	_, err := db.Select(&images, "SELECT id, product_id, image FROM product_images WHERE product_id=?", product.ID)
	if err != nil {
		fmt.Println("Find product images error", err)
		return images, err
	}

	return images, nil
}

func FindProductVariants(product structs.Product) ([]structs.ProductVariant, error) {

	var variants = []structs.ProductVariant{}

	db := app.GetDB()

	_, err := db.Select(&variants, "SELECT id, product_id, name FROM product_variants WHERE product_id=?", product.ID)
	if err != nil {
		fmt.Println("Find product variants error", err)
		return variants, err
	}

	return variants, nil
}

func FindProductCategories(product structs.Product) ([]structs.Category, error) {

	var categories = []structs.Category{}

	db := app.GetDB()

	_, err := db.Select(&categories, "SELECT categories.id, categories.name, categories.created_at, categories.updated_at FROM product_categories JOIN categories ON product_categories.category_id = categories.id WHERE product_categories.product_id=?", product.ID)
	if err != nil {
		fmt.Println("Find product categories error", err)
		return categories, err
	}

	return categories, nil
}

func FindAllProductImageIn(products []structs.Product) ([]structs.ProductImage, error) {

	var images []structs.ProductImage
	var productIDs string

	db := app.GetDB()

	productIDs = "("
	for _, product := range products {
		productIDs += "'" + product.ID + "',"
	}
	productIDs += ")"
	productIDs = strings.Replace(productIDs, ",)", ")", -1)

	_, err := db.Select(&images, "SELECT id, product_id, image FROM product_images WHERE product_id IN "+productIDs)
	if err != nil {
		fmt.Println("Find all product image in products error", err)
		return images, err
	}

	return images, nil
}

func FindAllProductVariantIn(products []structs.Product) ([]structs.ProductVariant, error) {

	var variants = []structs.ProductVariant{}
	var productIDs string

	db := app.GetDB()

	productIDs = "("
	for _, product := range products {
		productIDs += "'" + product.ID + "',"
	}
	productIDs += ")"
	productIDs = strings.Replace(productIDs, ",)", ")", -1)

	_, err := db.Select(&variants, "SELECT id, product_id, name FROM product_variants WHERE product_id IN "+productIDs)
	if err != nil {
		fmt.Println("Find all product variant in products error", err)
		return variants, err
	}

	return variants, nil
}

func FindAllProductVariantOptionIn(productVariants []structs.ProductVariant) ([]structs.ProductVariantOption, error) {

	var variantOptions = []structs.ProductVariantOption{}
	var productVariantIDs string

	db := app.GetDB()

	productVariantIDs = "("
	for _, productVariant := range productVariants {
		productVariantIDs += "'" + productVariant.ID + "',"
	}
	productVariantIDs += ")"
	productVariantIDs = strings.Replace(productVariantIDs, ",)", ")", -1)

	_, err := db.Select(&variantOptions, "SELECT id, product_variant_id, name FROM product_variant_options WHERE product_variant_id IN "+productVariantIDs)
	if err != nil {
		fmt.Println("Find all variant option in variants error", err)
		return variantOptions, err
	}

	return variantOptions, nil
}

func FindAllProductCategoryIn(products []structs.Product) ([]structs.Category, error) {

	var categories = []structs.Category{}
	var productIDs string

	db := app.GetDB()

	productIDs = "("
	for _, product := range products {
		productIDs += "'" + product.ID + "',"
	}
	productIDs += ")"
	productIDs = strings.Replace(productIDs, ",)", ")", -1)

	_, err := db.Select(&categories, "SELECT categories.id, categories.name, product_categories.product_id FROM product_categories JOIN categories ON product_categories.category_id = categories.id WHERE product_categories.product_id IN "+productIDs)
	if err != nil {
		fmt.Println("Find all product category in products error", err)
		return categories, err
	}

	return categories, nil
}

func StoreProduct(request structs.ProductRequest) (structs.Product, error) {

	db := app.GetDB()

	var product = structs.Product{
		ID:          uuid.New().String(),
		Name:        request.Name,
		Price:       request.Price,
		Description: request.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err := db.Exec("INSERT INTO products(id, name, price, description, created_at, updated_at) VALUES (?,?,?,?,?,?)", product.ID, product.Name, product.Price, product.Description, product.CreatedAt, product.UpdatedAt)
	if err != nil {
		fmt.Println("Store product error", err)
		return product, err
	}

	return product, nil
}

func UpdateProduct(request structs.ProductUpdateRequest, product structs.Product) (structs.Product, error) {

	var err error
	var newProduct = structs.Product{
		ID:          product.ID,
		Name:        request.Name,
		Price:       request.Price,
		Description: request.Description,
	}

	db := app.GetDB()

	_, err = db.Exec("UPDATE products SET name=?, price=?, description=? WHERE id=?", request.Name, request.Price, request.Description, product.ID)
	if err != nil {
		fmt.Println("Update product error", err)
		return newProduct, err
	}

	return newProduct, err
}

func StoreProductCategory(category structs.Category, product structs.Product) (structs.Category, error) {

	db := app.GetDB()

	_, err := db.Exec("INSERT INTO product_categories(id, product_id, category_id) VALUES (?,?,?)", uuid.New().String(), product.ID, category.ID)
	if err != nil {
		fmt.Println("Store product category error", err)
		return category, err
	}

	return category, err
}

func DeleteAllImageFromProduct(product structs.Product) error {

	db := app.GetDB()

	_, err := db.Exec("DELETE FROM product_images WHERE product_id=?", product.ID)
	if err != nil {
		fmt.Println("Delete all image from product error", err)
		return err
	}

	return nil
}

func DeleteAllVariantAndOptionFromProduct(product structs.Product) error {

	db := app.GetDB()

	_, err := db.Exec("DELETE pv, pvo FROM product_variants AS pv, product_variant_options AS pvo WHERE pv.product_id=?", product.ID)
	if err != nil {
		fmt.Println("Delete all variant error", err)
		return err
	}

	return nil
}

func DeleteAllCategoryFromProduct(product structs.Product) error {

	db := app.GetDB()

	_, err := db.Exec("DELETE FROM product_categories WHERE product_id=?", product.ID)
	if err != nil {
		fmt.Println("Delete all category error", err)
		return err
	}

	return nil
}
