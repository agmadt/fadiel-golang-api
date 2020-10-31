package controllers

import (
	"database/sql"
	"fmt"
	"golang-api/app/helpers"
	"golang-api/app/models"
	"golang-api/app/structs"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

type ProductController struct{}

func (controller ProductController) Index(c *gin.Context) {

	var productResponses = []structs.ProductResponse{}
	var productImages = []structs.ProductImage{}
	var productVariants = []structs.ProductVariant{}
	var productVariantOptions = []structs.ProductVariantOption{}
	var productCategories = []structs.Category{}
	var products []structs.Product
	var limit int = 10
	var page int = 1
	var err error

	if c.Query("limit") != "" {
		limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			c.JSON(500, gin.H{"message": "Limit is not a number"})
			return
		}
	}

	if c.Query("page") != "" {
		page, err = strconv.Atoi(c.Query("page"))
		if err != nil {
			c.JSON(500, gin.H{"message": "Page is not a number"})
			return
		}
	}

	paginatedProducts, err := models.PaginateProduct(structs.PaginationParameters{
		Limit: limit,
		Page:  page,
	})
	if err != nil {
		c.JSON(500, gin.H{"message": "Server error"})
		return
	}

	for _, product := range paginatedProducts.Products {
		products = append(products, structs.Product{
			ID: product.ID,
		})
	}

	if len(products) > 0 {
		productImages, err = models.FindAllProductImageIn(products)
		if err != nil && err != sql.ErrNoRows {
			c.JSON(500, gin.H{"message": "Server error"})
			return
		}

		productVariants, err = models.FindAllProductVariantIn(products)
		if err != nil && err != sql.ErrNoRows {
			c.JSON(500, gin.H{"message": "Server error"})
			return
		}

		if len(productVariants) > 0 {
			productVariantOptions, err = models.FindAllProductVariantOptionIn(productVariants)
			if err != nil && err != sql.ErrNoRows {
				c.JSON(500, gin.H{"message": "Server error"})
				return
			}
		}

		productCategories, err = models.FindAllProductCategoryIn(products)
		if err != nil && err != sql.ErrNoRows {
			c.JSON(500, gin.H{"message": "Server error"})
			return
		}

		for _, product := range paginatedProducts.Products {
			for _, productImage := range productImages {
				if product.ID == productImage.ProductID {
					product.Images = append(product.Images, structs.ProductImageResponse{
						ID:    productImage.ID,
						Image: productImage.Image,
					})
				}
			}

			for _, productVariant := range productVariants {
				if product.ID == productVariant.ProductID {

					productVariantResponse := structs.ProductVariantResponse{
						ID:   productVariant.ID,
						Name: productVariant.Name,
					}

					for _, productVariantOption := range productVariantOptions {
						if productVariant.ID == productVariantOption.ProductVariantID {
							productVariantResponse.Options = append(productVariantResponse.Options, structs.ProductVariantOptionResponse{
								ID:   productVariantOption.ID,
								Name: productVariantOption.Name,
							})
						}
					}

					product.Variants = append(product.Variants, productVariantResponse)
				}
			}

			for _, productCategory := range productCategories {
				if product.ID == productCategory.ProductID {
					product.Categories = append(product.Categories, productCategory.Response())
				}
			}

			productResponses = append(productResponses, structs.ProductResponse{
				ID:          product.ID,
				Name:        product.Name,
				Price:       product.Price,
				Description: product.Description,
				Images:      product.Images,
				Variants:    product.Variants,
				Categories:  product.Categories,
			})
		}
	}

	paginatedProducts.Products = productResponses

	c.JSON(200, paginatedProducts)
}

func (controller ProductController) Show(c *gin.Context) {

	var product = structs.Product{ID: c.Param("id")}
	var productImageResponses = []structs.ProductImageResponse{}
	var productVariantResponses = []structs.ProductVariantResponse{}
	var productCategoryResponses = []structs.CategoryResponse{}
	var productVariantOptions []structs.ProductVariantOption

	product, err := models.FindProduct(product)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"message": "Product not found"})
		} else {
			c.JSON(500, gin.H{"message": "Server error"})
		}
		return
	}

	productImages, err := models.FindProductImages(product)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(500, gin.H{"message": "Server error"})
		return
	}

	productVariants, err := models.FindProductVariants(product)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(500, gin.H{"message": "Server error"})
		return
	}

	productCategories, err := models.FindProductCategories(product)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(500, gin.H{"message": "Server error"})
		return
	}

	if len(productVariants) > 0 {
		productVariantOptions, err = models.FindAllProductVariantOptionIn(productVariants)
		if err != nil && err != sql.ErrNoRows {
			c.JSON(500, gin.H{"message": "Server error"})
			return
		}
	}

	for _, productImage := range productImages {
		productImageResponses = append(productImageResponses, structs.ProductImageResponse{
			ID:    productImage.ID,
			Image: productImage.Image,
		})
	}

	for _, productVariant := range productVariants {
		productVariantResponse := structs.ProductVariantResponse{
			ID:   productVariant.ID,
			Name: productVariant.Name,
		}

		for _, productVariantOption := range productVariantOptions {
			if productVariant.ID == productVariantOption.ProductVariantID {
				productVariantResponse.Options = append(productVariantResponse.Options, structs.ProductVariantOptionResponse{
					ID:   productVariantOption.ID,
					Name: productVariantOption.Name,
				})
			}
		}

		productVariantResponses = append(productVariantResponses, productVariantResponse)
	}

	for _, productCategory := range productCategories {
		productCategoryResponses = append(productCategoryResponses, productCategory.Response())
	}

	var productResponse = structs.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
		Images:      productImageResponses,
		Variants:    productVariantResponses,
		Categories:  productCategoryResponses,
	}

	c.JSON(200, productResponse)
}

func (controller ProductController) Store(c *gin.Context) {

	var productRequest structs.ProductRequest
	var failedValidations = map[string]interface{}{}

	err := c.ShouldBindJSON(&productRequest)
	if err != nil {
		c.JSON(400, gin.H{"message": "Something wrong with the request"})
		fmt.Println("Store product bind json error", err)
		return
	}

	validate = validator.New()
	err = validate.Struct(productRequest)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errFailedField := helpers.ValidatorRemoveNamespace(strcase.ToSnake(err.Namespace()))
			failedValidations[errFailedField] = []string{helpers.ValidatorMessage(errFailedField, err.ActualTag(), err.Param())}
		}

		c.JSON(422, helpers.Validator{
			Message: "The given data was invalid",
			Errors:  failedValidations,
		})
		return
	}

	if len(productRequest.Categories) > 0 {
		for _, productCategoryRequest := range productRequest.Categories {

			_, err := models.FindCategory(structs.Category{ID: productCategoryRequest.ID})
			if err != nil {
				if err == sql.ErrNoRows {
					c.JSON(404, gin.H{"message": "Category not found"})
				} else {
					c.JSON(500, gin.H{"message": "Server error"})
				}
				return
			}
		}
	}

	product, err := models.StoreProduct(productRequest)
	if err != nil {
		c.JSON(500, gin.H{"message": "Server error"})
		return
	}

	if len(productRequest.Images) > 0 {
		for _, productImageRequest := range productRequest.Images {
			_, err := models.StoreProductImage(productImageRequest, product)
			if err != nil {
				c.JSON(500, gin.H{"message": "Server error"})
				return
			}
		}
	}

	if len(productRequest.Variants) > 0 {
		for _, productVariantRequest := range productRequest.Variants {
			productVariant, err := models.StoreProductVariant(productVariantRequest, product)
			if err != nil {
				c.JSON(500, gin.H{"message": "Server error"})
				return
			}

			for _, productVariantOptionRequest := range productVariantRequest.Options {
				_, err := models.StoreProductVariantOption(productVariantOptionRequest, productVariant)
				if err != nil {
					c.JSON(500, gin.H{"message": "Server error"})
					return
				}
			}
		}
	}

	if len(productRequest.Categories) > 0 {
		for _, productCategoryRequest := range productRequest.Categories {

			category := structs.Category{ID: productCategoryRequest.ID}
			_, err = models.StoreProductCategory(category, product)
			if err != nil {
				c.JSON(500, gin.H{"message": "Server error"})
				return
			}
		}
	}

	c.JSON(200, product)
}

func (controller ProductController) Update(c *gin.Context) {

	var productRequest structs.ProductUpdateRequest
	var failedValidations = map[string]interface{}{}

	product, err := models.FindProduct(structs.Product{ID: c.Param("id")})
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"message": "Product not found"})
		} else {
			c.JSON(500, gin.H{"message": "Server error"})
		}
		return
	}

	err = c.ShouldBindJSON(&productRequest)
	if err != nil {
		c.JSON(400, gin.H{"message": "Something wrong with the request"})
		fmt.Println("Update product bind json error", err)
		return
	}

	validate = validator.New()
	err = validate.Struct(productRequest)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errFailedField := helpers.ValidatorRemoveNamespace(strcase.ToSnake(err.Namespace()))
			failedValidations[errFailedField] = []string{helpers.ValidatorMessage(errFailedField, err.ActualTag(), err.Param())}
		}

		c.JSON(422, helpers.Validator{
			Message: "The given data was invalid",
			Errors:  failedValidations,
		})
		return
	}

	if len(productRequest.Categories) > 0 {
		for _, productCategoryRequest := range productRequest.Categories {
			_, err := models.FindCategory(structs.Category{ID: productCategoryRequest.ID})
			if err != nil {
				if err == sql.ErrNoRows {
					c.JSON(404, gin.H{"message": "Category not found"})
				} else {
					c.JSON(500, gin.H{"message": "Server error"})
				}
				return
			}
		}
	}

	product, err = models.UpdateProduct(productRequest, product)
	if err != nil {
		c.JSON(500, gin.H{"message": "Server error"})
		return
	}

	if len(productRequest.Images) > 0 {
		err = models.DeleteAllImageFromProduct(product)
		if err != nil {
			c.JSON(500, gin.H{"message": "Server error"})
			return
		}

		for _, productImageRequest := range productRequest.Images {
			_, err := models.StoreProductImage(productImageRequest, product)
			if err != nil {
				c.JSON(500, gin.H{"message": "Server error"})
				return
			}
		}
	}

	if len(productRequest.Variants) > 0 {
		err = models.DeleteAllVariantAndOptionFromProduct(product)
		if err != nil {
			c.JSON(500, gin.H{"message": "Server error"})
			return
		}

		for _, productVariantRequest := range productRequest.Variants {
			productVariant, err := models.StoreProductVariant(productVariantRequest, product)
			if err != nil {
				c.JSON(500, gin.H{"message": "Server error"})
				return
			}

			for _, productVariantOptionRequest := range productVariantRequest.Options {
				_, err := models.StoreProductVariantOption(productVariantOptionRequest, productVariant)
				if err != nil {
					c.JSON(500, gin.H{"message": "Server error"})
					return
				}
			}
		}
	}

	if len(productRequest.Categories) > 0 {
		err = models.DeleteAllCategoryFromProduct(product)
		if err != nil {
			c.JSON(500, gin.H{"message": "Server error"})
			return
		}

		for _, productCategoryRequest := range productRequest.Categories {
			category := structs.Category{ID: productCategoryRequest.ID}
			_, err = models.StoreProductCategory(category, product)
			if err != nil {
				c.JSON(500, gin.H{"message": "Server error"})
				return
			}
		}
	}

	c.JSON(200, product)
}

func (controller ProductController) Delete(c *gin.Context) {

	var product = structs.Product{ID: c.Param("id")}

	product, err := models.FindProduct(product)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"message": "Product not found"})
		} else {
			c.JSON(500, gin.H{"message": "Server error"})
		}
		return
	}

	err = models.DeleteProduct(product)
	if err != nil {
		c.JSON(500, gin.H{"message": "Server error"})
		return
	}

	c.JSON(200, gin.H{"message": "Product successfully deleted"})
}
