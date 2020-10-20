package models

import (
	"database/sql"
	"fmt"
	"golang-api/app"
	"golang-api/app/structs"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func PaginateCategory(paginationParam structs.PaginationParameters) (structs.CategoryPagination, error) {
	var categories []structs.Category
	var totalResponse int
	var categoryResponses []structs.CategoryResponse
	var categoryPagination structs.CategoryPagination
	var offset string = "0"
	var sqlQuery string
	var limit string = strconv.Itoa(paginationParam.Limit)

	offset = strconv.Itoa(((paginationParam.Page) - 1) * paginationParam.Limit)

	db := app.GetDB()

	sqlQuery = "SELECT id, name FROM categories ORDER BY created_at DESC"
	sqlQuery += " LIMIT " + offset + ", " + limit

	_, err := db.Select(&categories, sqlQuery)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("Query to get category error", err)
		return categoryPagination, err
	}

	err = db.SelectOne(&totalResponse, "SELECT count(id) as total FROM categories")
	if err != nil && err != sql.ErrNoRows {
		fmt.Println("Query to count all category error", err)
		return categoryPagination, err
	}

	for _, category := range categories {
		categoryResponses = append(categoryResponses, category.Response())
	}

	categoryPagination = structs.CategoryPagination{
		Categories: categoryResponses,
		Total:      totalResponse,
		Limit:      paginationParam.Limit,
		Page:       paginationParam.Page,
	}

	return categoryPagination, nil
}

func FindCategory(category structs.Category) (structs.Category, error) {

	db := app.GetDB()

	err := db.SelectOne(&category, "SELECT id, name FROM categories WHERE id=?", category.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			fmt.Println("Find category error", err)
		}
		return category, err
	}

	return category, nil
}

func FindCategoryProducts(category structs.Category) ([]structs.Category, error) {

	var categories = []structs.Category{}

	db := app.GetDB()

	_, err := db.Select(&categories, "SELECT id FROM product_categories WHERE category_id=?", category.ID)
	if err != nil {
		fmt.Println("Find category products error", err)
		return categories, err
	}

	return categories, nil
}

func StoreCategory(request structs.CategoryRequest) (structs.Category, error) {

	var category = structs.Category{
		ID:        uuid.New().String(),
		Name:      request.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	db := app.GetDB()

	_, err := db.Exec("INSERT INTO categories(id, name, created_at, updated_at) VALUES (?,?,?,?)", category.ID, category.Name, category.CreatedAt, category.UpdatedAt)
	if err != nil {
		fmt.Println("Store category error", err)
		return category, err
	}

	return category, nil
}

func UpdateCategory(request structs.CategoryRequest, category structs.Category) (structs.Category, error) {

	category.Name = request.Name

	db := app.GetDB()

	_, err := db.Exec("UPDATE categories SET name = ? WHERE id = ?", category.Name, category.ID)
	if err != nil {
		fmt.Println("Update category error", err)
		return category, err
	}

	return category, nil
}

func DeleteCategory(category structs.Category) (structs.Category, error) {

	db := app.GetDB()

	_, err := db.Exec("DELETE FROM categories WHERE id = ?", category.ID)
	if err != nil {
		fmt.Println("Delete category error", err)
		return category, err
	}

	return category, nil
}
