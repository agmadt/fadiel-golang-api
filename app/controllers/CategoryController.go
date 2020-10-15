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

type CategoryController struct{}

func (controller CategoryController) Index(c *gin.Context) {

	categories, err := models.FindAllCategory(c)
	if err != nil {
		return
	}

	c.JSON(200, categories)
}

func (controller CategoryController) Show(c *gin.Context) {

	id := c.Param("id")

	category, err := models.FindCategory(c, id)
	if err != nil {
		return
	}

	c.JSON(200, category.Response())
}

func (controller CategoryController) Store(c *gin.Context) {

	var categoryRequest structs.CategoryRequest
	var failedValidations map[string]interface{}

	err := c.ShouldBindJSON(&categoryRequest)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Something wrong with the request",
		})
		fmt.Println(err)
		return
	}

	validate = validator.New()
	err = validate.Struct(categoryRequest)

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

	category, err := models.StoreCategory(c, categoryRequest)
	if err != nil {
		return
	}

	c.JSON(200, category.Response())
}

func (controller CategoryController) Update(c *gin.Context) {

	var categoryRequest structs.CategoryRequest
	var failedValidations map[string]interface{}

	id := c.Param("id")

	category, err := models.FindCategory(c, id)
	if err != nil {
		return
	}

	err = c.ShouldBindJSON(&categoryRequest)
	if err != nil {
		c.JSON(400, gin.H{
			"message": "Something wrong with the request",
		})
		fmt.Println(err)
		return
	}

	validate = validator.New()
	err = validate.Struct(categoryRequest)

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

	category, err = models.UpdateCategory(c, categoryRequest, category)
	if err != nil {
		return
	}

	c.JSON(200, category.Response())
}

func (controller CategoryController) Delete(c *gin.Context) {

	var categories []structs.Category

	id := c.Param("id")

	category, err := models.FindCategory(c, id)
	if err != nil {
		return
	}

	categories, err = models.FindCategoryProducts(c, category)
	if err != nil {
		return
	}

	if len(categories) > 0 {
		c.JSON(403, gin.H{
			"message": "Category is still being used by products",
		})
		return
	}

	_, err = models.DeleteCategory(c, category)
	if err != nil {
		return
	}

	c.JSON(400, gin.H{
		"message": "Category successfully deleted",
	})
}
