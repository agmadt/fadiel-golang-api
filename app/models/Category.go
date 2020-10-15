package models

import (
	"fmt"
	"golang-api/app"
	"golang-api/app/structs"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func FindAllCategory(c *gin.Context) (structs.CategoryPagination, error) {
	var categories []structs.Category
	var totalResponse int
	var categoryResponses []structs.CategoryResponse
	var categoryPagination structs.CategoryPagination
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

	sqlQuery = "SELECT id, name FROM categories ORDER BY created_at DESC"
	sqlQuery += " LIMIT " + offset + ", " + limit

	_, err := db.Select(&categories, sqlQuery)

	fmt.Println(sqlQuery)

	if err != nil {
		c.JSON(404, gin.H{
			"message": "Error while getting all records",
		})
		fmt.Println(err)

		return categoryPagination, err
	}

	err = db.SelectOne(&totalResponse, "SELECT count(id) as total FROM categories")

	if err != nil {
		c.JSON(404, gin.H{
			"message": "Error while getting all records",
		})
		fmt.Println(err)

		return categoryPagination, err
	}

	for _, category := range categories {
		categoryResponses = append(categoryResponses, category.Response())
	}

	categoryPagination = structs.CategoryPagination{
		Categories: categoryResponses,
		Total:      totalResponse,
		Limit:      limitInt,
		Page:       pageInt,
	}

	return categoryPagination, err
}

func FindCategory(c *gin.Context, pk string) (structs.Category, error) {

	var category structs.Category

	db := app.GetDB()

	err := db.SelectOne(&category, "SELECT id, name FROM categories WHERE id=?", pk)

	if err != nil {
		c.JSON(404, gin.H{
			"message": "Category not found",
		})
		fmt.Println(err)

		return category, err
	}

	return category, err
}

func FindCategoryProducts(c *gin.Context, category structs.Category) ([]structs.Category, error) {

	var categories []structs.Category

	db := app.GetDB()

	_, err := db.Select(&categories, "SELECT id FROM product_categories WHERE category_id=?", category.ID)
	if err != nil {
		c.JSON(404, gin.H{
			"message": "Category products not found",
		})
		fmt.Println(err)

		return categories, err
	}

	return categories, err
}

func StoreCategory(c *gin.Context, request structs.CategoryRequest) (structs.Category, error) {

	var category = structs.Category{
		ID:        uuid.New().String(),
		Name:      request.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db := app.GetDB()

	_, err := db.Exec("INSERT INTO categories(id, name, created_at, updated_at) VALUES (?,?,?,?)", category.ID, category.Name, category.CreatedAt, category.UpdatedAt)

	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error while storing data",
		})
		fmt.Println(err)

		return category, err
	}

	return category, err
}

func UpdateCategory(c *gin.Context, request structs.CategoryRequest, category structs.Category) (structs.Category, error) {

	category.Name = request.Name

	db := app.GetDB()

	_, err := db.Exec("UPDATE categories SET name = ? WHERE id = ?", category.Name, category.ID)

	if err != nil {
		c.JSON(500, gin.H{
			"message": "Error while updating data",
		})
		fmt.Println(err)

		return category, err
	}

	return category, err
}

func DeleteCategory(c *gin.Context, category structs.Category) (structs.Category, error) {

	db := app.GetDB()

	_, err := db.Exec("DELETE FROM categories WHERE id = ?", category.ID)

	if err != nil {
		c.JSON(500, gin.H{
			"message": "Server error",
		})
		fmt.Println(err)

		return category, err
	}

	return category, err
}
