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

type CategoryController struct{}

func (controller CategoryController) Index(c *gin.Context) {

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

	categories, err := models.PaginateCategory(structs.PaginationParameters{
		Limit: limit,
		Page:  page,
	})
	if err != nil {
		c.JSON(500, gin.H{"message": "Server error"})
		return
	}

	c.JSON(200, categories)
}

func (controller CategoryController) Show(c *gin.Context) {

	var category = structs.Category{ID: c.Param("id")}

	category, err := models.FindCategory(category)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"message": "Category not found"})
		} else {
			c.JSON(500, gin.H{"message": "Server error"})
		}
		return
	}

	c.JSON(200, category.Response())
}

func (controller CategoryController) Store(c *gin.Context) {

	var categoryRequest structs.CategoryRequest
	var failedValidations map[string]interface{}

	err := c.ShouldBindJSON(&categoryRequest)
	if err != nil {
		c.JSON(400, gin.H{"message": "Something wrong with the request"})
		fmt.Println("Store category bind json error", err)
		return
	}

	validate = validator.New()
	err = validate.Struct(categoryRequest)

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

	category, err := models.StoreCategory(categoryRequest)
	if err != nil {
		c.JSON(500, gin.H{"message": "Server error"})
		return
	}

	c.JSON(200, category.Response())
}

func (controller CategoryController) Update(c *gin.Context) {

	var category structs.Category
	var categoryRequest structs.CategoryRequest
	var failedValidations map[string]interface{}

	category.ID = c.Param("id")

	category, err := models.FindCategory(category)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"message": "Category not found"})
		} else {
			c.JSON(500, gin.H{"message": "Server error"})
		}
		return
	}

	err = c.ShouldBindJSON(&categoryRequest)
	if err != nil {
		c.JSON(400, gin.H{"message": "Something wrong with the request"})
		fmt.Println("Store category bind json error", err)
		return
	}

	validate = validator.New()
	err = validate.Struct(categoryRequest)

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

	category, err = models.UpdateCategory(categoryRequest, category)
	if err != nil {
		c.JSON(500, gin.H{"message": "Server error"})
		return
	}

	c.JSON(200, category.Response())
}

func (controller CategoryController) Delete(c *gin.Context) {

	var category structs.Category
	var categories []structs.Category

	category.ID = c.Param("id")

	category, err := models.FindCategory(category)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"message": "Category not found"})
		} else {
			c.JSON(500, gin.H{"message": "Server error"})
		}
		return
	}

	categories, err = models.FindCategoryProducts(category)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(500, gin.H{"message": "Server error"})
		return
	}

	if len(categories) > 0 {
		c.JSON(403, gin.H{"message": "Category is still being used by products"})
		return
	}

	_, err = models.DeleteCategory(category)
	if err != nil {
		c.JSON(500, gin.H{"message": "Server error"})
		return
	}

	c.JSON(400, gin.H{"message": "Category successfully deleted"})
}
