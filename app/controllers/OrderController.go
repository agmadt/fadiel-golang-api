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

type OrderController struct{}

var validate *validator.Validate

func (controller OrderController) Index(c *gin.Context) {

	orders, err := models.FindAllOrder(c)

	if err != nil {
		return
	}

	c.JSON(200, orders.Response())
}

func (controller OrderController) Show(c *gin.Context) {

	id := c.Param("id")

	order, err := models.FindOrder(c, id)
	if err != nil {
		return
	}

	orderProducts, err := models.FindOrderProducts(c, order)
	if err != nil {
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
		c.JSON(400, gin.H{
			"message": "Something wrong with the request",
		})
		fmt.Println(err)
		return
	}

	validate = validator.New()
	err = validate.Struct(orderRequest)

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

	var totalOrder float32 = 0

	for _, itemProduct := range orderRequest.Products {
		product, err := models.FindProduct(c, itemProduct.ID)
		if err != nil {
			return
		}

		productImages, err := models.FindProductImages(c, product)
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(productImages) > 0 {
			for _, productImage := range productImages {
				orderProductImages = append(orderProductImages, structs.OrderProductProductImage{
					Image: productImage.Image,
				})
			}
		}

		productCategories, err := models.FindProductCategories(c, product)
		if err != nil {
			fmt.Println(err)
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
			productVariant, err := models.FindProductVariant(c, itemVariant.VariantID)
			if err != nil {
				return
			}

			productVariantOption, err := models.FindProductVariantOption(c, itemVariant.OptionID)
			if err != nil {
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

	order, err := models.StoreOrder(c, orderRequest)
	if err != nil {
		return
	}

	for _, orderProduct := range orderProducts {
		_, err := models.StoreOrderProduct(c, order, orderProduct)
		if err != nil {
			return
		}
	}

	c.JSON(200, order.Response())
}
