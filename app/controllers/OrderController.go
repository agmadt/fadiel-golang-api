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

type OrderController struct{}

var validate *validator.Validate

func (controller OrderController) Index(c *gin.Context) {

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

	paginatedOrders, err := models.PaginateOrder(structs.PaginationParameters{
		Limit: limit,
		Page:  page,
	})
	if err != nil {
		c.JSON(500, gin.H{"message": "Server error"})
		return
	}

	c.JSON(200, paginatedOrders)
}

func (controller OrderController) Show(c *gin.Context) {

	var order = structs.Order{
		ID: c.Param("id"),
	}

	order, err := models.FindOrder(order)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"message": "Order not found"})
		} else {
			c.JSON(500, gin.H{"message": "Server error"})
		}
		return
	}

	orderProducts, err := models.FindOrderProducts(order)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(500, gin.H{"message": "Server error"})
		return
	}

	var orderProduct = structs.OrderProductResponse{
		OrderResponse: order.Response(),
		Products:      orderProducts,
	}

	c.JSON(200, orderProduct)
}

func (controller OrderController) Store(c *gin.Context) {

	var orderRequest structs.OrderRequest
	var failedValidations = map[string]interface{}{}
	var orderProducts = []structs.OrderProductProduct{}
	var orderProductImages = []structs.OrderProductProductImage{}
	var orderProductVariants = []structs.OrderProductProductVariant{}
	var orderProductCategories = []structs.OrderProductProductCategory{}

	err := c.ShouldBindJSON(&orderRequest)
	if err != nil {
		c.JSON(400, gin.H{"message": "Something wrong with the request"})
		fmt.Println("Store order bind json error", err)
		return
	}

	validate = validator.New()
	err = validate.Struct(orderRequest)

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

	var totalOrder float32 = 0

	for _, itemProduct := range orderRequest.Products {

		var product = structs.Product{
			ID: itemProduct.ID,
		}

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

		if len(productImages) > 0 {
			for _, productImage := range productImages {
				orderProductImages = append(orderProductImages, structs.OrderProductProductImage{
					Image: productImage.Image,
				})
			}
		}

		productCategories, err := models.FindProductCategories(product)
		if err != nil && err != sql.ErrNoRows {
			c.JSON(500, gin.H{"message": "Server error"})
			return
		}

		if len(productCategories) > 0 {
			for _, productCategory := range productCategories {
				orderProductCategories = append(orderProductCategories, structs.OrderProductProductCategory{
					Name: productCategory.Name,
				})
			}
		}

		for _, itemVariant := range itemProduct.Variants {
			var productVariant = structs.ProductVariant{
				ID: itemVariant.VariantID,
			}

			productVariant, err := models.FindProductVariant(productVariant)
			if err != nil {
				if err == sql.ErrNoRows {
					c.JSON(404, gin.H{"message": "Product variant not found"})
				} else {
					c.JSON(500, gin.H{"message": "Server error"})
				}
				return
			}

			var productVariantOption = structs.ProductVariantOption{
				ID: itemVariant.OptionID,
			}

			productVariantOption, err = models.FindProductVariantOption(productVariantOption)
			if err != nil {
				if err == sql.ErrNoRows {
					c.JSON(404, gin.H{"message": "Product variant option not found"})
				} else {
					c.JSON(500, gin.H{"message": "Server error"})
				}
				return
			}

			orderProductVariants = append(orderProductVariants, structs.OrderProductProductVariant{
				VariantID:         productVariant.ID,
				VariantName:       productVariant.Name,
				VariantOptionID:   productVariantOption.ID,
				VariantOptionName: productVariantOption.ID,
			})
		}

		subtotal := (product.Price * float32(itemProduct.Quantity))
		totalOrder += (product.Price * float32(itemProduct.Quantity))

		orderProducts = append(orderProducts, structs.OrderProductProduct{
			ID:          product.ID,
			Name:        product.Name,
			Price:       product.Price,
			Description: product.Description,
			Quantity:    itemProduct.Quantity,
			Subtotal:    subtotal,
			Images:      orderProductImages,
			Variants:    orderProductVariants,
			Categories:  orderProductCategories,
		})
	}

	orderRequest.Total = totalOrder

	order, err := models.StoreOrder(orderRequest)
	if err != nil {
		c.JSON(500, gin.H{"message": "Server error"})
		return
	}

	for _, orderProduct := range orderProducts {
		_, err := models.StoreOrderProduct(order, orderProduct)
		if err != nil {
			c.JSON(500, gin.H{"message": "Server error"})
			return
		}
	}

	c.JSON(200, order.Response())
}
