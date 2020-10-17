package controllers

import (
	"fmt"
	"golang-api/app/helpers"
	"golang-api/app/models"
	"golang-api/app/structs"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

type ProductController struct{}

func (controller ProductController) Index(c *gin.Context) {

	var productResponses []structs.ProductResponse
	var products []structs.Product

	paginatedProducts, err := models.FindAllProduct(c)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Server error",
		})
		return
	}

	for _, product := range paginatedProducts.Products {
		products = append(products, structs.Product{
			ID: product.ID,
		})
	}

	productImages, err := models.FindAllProductImageIn(products)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Server error",
		})
		return
	}

	productVariants, err := models.FindAllProductVariantIn(products)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Server error",
		})
		return
	}

	productVariantOptions, err := models.FindAllProductVariantOptionIn(productVariants)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Server error",
		})
		return
	}

	productCategories, err := models.FindAllProductCategoryIn(products)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Server error",
		})
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
				product.Categories = append(product.Categories, structs.ProductCategoryResponse{
					ID:   productCategory.ID,
					Name: productCategory.Name,
				})
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

	paginatedProducts.Products = productResponses

	c.JSON(200, paginatedProducts)
}

func (controller ProductController) Show(c *gin.Context) {

	var product = structs.Product{ID: c.Param("id")}
	var productImageResponses []structs.ProductImageResponse
	var productVariantResponses []structs.ProductVariantResponse
	var productCategoryResponses []structs.ProductCategoryResponse
	var productVariantOptions []structs.ProductVariantOption

	product, err := models.FindProduct(c, product)
	if err != nil {
		return
	}

	productImages, err := models.FindProductImages(c, product)
	if err != nil {
		return
	}

	productVariants, err := models.FindProductVariants(c, product)
	if err != nil {
		return
	}

	productCategories, err := models.FindProductCategories(c, product)
	if err != nil {
		return
	}

	if len(productVariants) > 0 {
		productVariantOptions, err = models.FindAllProductVariantOptionIn(productVariants)
		if err != nil {
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
			productVariantResponse.Options = append(productVariantResponse.Options, structs.ProductVariantOptionResponse{
				ID:   productVariantOption.ID,
				Name: productVariantOption.Name,
			})
		}

		productVariantResponses = append(productVariantResponses, productVariantResponse)
	}

	for _, productCategory := range productCategories {
		productCategoryResponses = append(productCategoryResponses, structs.ProductCategoryResponse{
			ID:   productCategory.ID,
			Name: productCategory.Name,
		})
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
		c.JSON(400, gin.H{
			"message": "Something wrong with the request",
		})
		fmt.Println(err)
		return
	}

	validate = validator.New()
	err = validate.Struct(productRequest)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errFailedField := helpers.ValidatorRemoveNamespace(strcase.ToSnake(err.Namespace()))
			failedValidations[errFailedField] = []string{helpers.ValidatorMessage(errFailedField, err.ActualTag(), err.Param())}
		}

		c.JSON(200, helpers.Validator{
			Message: "The given data was invalid",
			Errors:  failedValidations,
		})
		return
	}

	product, err := models.StoreProduct(c, productRequest)
	if err != nil {
		return
	}

	if len(productRequest.Images) > 0 {
		for _, productImageRequest := range productRequest.Images {
			_, err := models.StoreProductImage(c, productImageRequest, product)
			if err != nil {
				return
			}
		}
	}

	if len(productRequest.Variants) > 0 {
		for _, productVariantRequest := range productRequest.Variants {
			productVariant, err := models.StoreProductVariant(c, productVariantRequest, product)
			if err != nil {
				return
			}

			for _, productVariantOptionRequest := range productVariantRequest.Options {
				_, err := models.StoreProductVariantOption(c, productVariantOptionRequest, productVariant)
				if err != nil {
					return
				}
			}
		}
	}

	if len(productRequest.Categories) > 0 {
		for _, productCategoryRequest := range productRequest.Categories {

			category, err := models.FindCategory(c, structs.Category{ID: productCategoryRequest.ID})
			if err != nil {
				return
			}

			_, err = models.StoreProductCategory(c, category, product)
			if err != nil {
				return
			}
		}
	}

	c.JSON(200, product)
}

func (controller ProductController) Update(c *gin.Context) {

	products, err := models.FindAllProduct(c)
	if err != nil {
		return
	}

	c.JSON(200, products)
}

func (controller ProductController) Delete(c *gin.Context) {

	products, err := models.FindAllProduct(c)
	if err != nil {
		return
	}

	c.JSON(200, products)
}
