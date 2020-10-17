package models

import (
	"fmt"
	"golang-api/app"
	"golang-api/app/structs"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func FindAllProduct(c *gin.Context) (structs.ProductPagination, error) {

	var products []structs.Product
	var totalResponse int
	var productResponses []structs.ProductResponse
	var productPagination structs.ProductPagination
	var limit string = "10"
	var page string = "1"
	var limitInt int
	var pageInt int
	var offset string = "0"
	var sqlQuery string

	if c.Query("limit") != "" {
		limit = c.Query("limit")
	}

	limitInt, _ = strconv.Atoi(limit)
	pageInt, _ = strconv.Atoi(page)

	if c.Query("page") != "" {
		page = c.Query("page")
		pageInt, _ := strconv.ParseInt(page, 10, 64)

		offsetInt := ((int(pageInt)) - 1) * limitInt
		offset = strconv.Itoa(int(offsetInt))
	}

	db := app.GetDB()

	sqlQuery = "SELECT id, name, price, description FROM products ORDER BY created_at DESC"
	sqlQuery += " LIMIT " + offset + ", " + limit

	_, err := db.Select(&products, sqlQuery)

	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error while getting all records",
		})
		fmt.Println(err)

		return productPagination, err
	}

	err = db.SelectOne(&totalResponse, "SELECT count(id) as total FROM products")

	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error while getting all records",
		})

		return productPagination, err
	}

	for _, product := range products {
		productResponses = append(productResponses, product.Response())
	}

	productPagination = structs.ProductPagination{
		Products: productResponses,
		Total:    totalResponse,
		Limit:    limitInt,
		Page:     pageInt,
	}

	return productPagination, err
}

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

func FindProductImages(c *gin.Context, product structs.Product) ([]structs.ProductImage, error) {

	var images = []structs.ProductImage{}

	db := app.GetDB()

	_, err := db.Select(&images, "SELECT id, product_id, image FROM product_images WHERE product_id=?", product.ID)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error while getting product images",
		})
		fmt.Println(err)

		return images, err
	}

	return images, err
}

func FindProductVariants(c *gin.Context, product structs.Product) ([]structs.ProductVariant, error) {

	var variants = []structs.ProductVariant{}

	db := app.GetDB()

	_, err := db.Select(&variants, "SELECT id, product_id, name FROM product_variants WHERE product_id=?", product.ID)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error while getting product variants",
		})
		fmt.Println(err)

		return variants, err
	}

	return variants, err
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

func FindAllProductImageIn(products []structs.Product) ([]structs.ProductImage, error) {

	var images = []structs.ProductImage{}
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
		fmt.Println(err)
	}

	return images, err
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
		fmt.Println(err)
	}

	return variants, err
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
		fmt.Println(err)
	}

	return variantOptions, err
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
		fmt.Println(err)
	}

	return categories, err
}

func StoreProduct(c *gin.Context, request structs.ProductRequest) (structs.Product, error) {

	var err error

	db := app.GetDB()

	var product = structs.Product{
		ID:          uuid.New().String(),
		Name:        request.Name,
		Price:       request.Price,
		Description: request.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	_, err = db.Exec("INSERT INTO products(id, name, price, description, created_at, updated_at) VALUES (?,?,?,?,?,?)", product.ID, product.Name, product.Price, product.Description, product.CreatedAt, product.UpdatedAt)

	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error while storing data",
		})
		fmt.Println(err)

		return product, err
	}

	return product, err
}

func StoreProductCategory(c *gin.Context, category structs.Category, product structs.Product) (structs.Category, error) {

	db := app.GetDB()

	_, err := db.Exec("INSERT INTO product_categories(id, product_id, category_id) VALUES (?,?,?)", uuid.New().String(), product.ID, category.ID)

	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error while storing data",
		})
		fmt.Println(err)

		return category, err
	}

	return category, err
}
