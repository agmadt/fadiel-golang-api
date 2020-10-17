package structs

import (
	helpers "golang-api/app/helpers"
	"time"
)

type Product struct {
	ID          string             `db:"id" json:"id"`
	Name        string             `db:"name" json:"name"`
	Price       float32            `db:"price" json:"price"`
	Description helpers.NullString `db:"description" json:"description"`
	CreatedAt   time.Time          `db:"created_at" json:"-"`
	UpdatedAt   time.Time          `db:"updated_at" json:"-"`
}

type ProductRequest struct {
	Name        string                   `json:"name"`
	Price       float32                  `json:"price"`
	Description helpers.NullString       `json:"description"`
	Images      []ProductImageRequest    `json:"images" validate:"dive"`
	Variants    []ProductVariantRequest  `json:"variants" validate:"dive"`
	Categories  []CategoryProductRequest `json:"categories" validate:"min=1,dive"`
}

type ProductResponse struct {
	ID          string                    `json:"id"`
	Name        string                    `json:"name"`
	Price       float32                   `json:"price"`
	Description helpers.NullString        `json:"description"`
	Images      []ProductImageResponse    `json:"images"`
	Variants    []ProductVariantResponse  `json:"variants"`
	Categories  []ProductCategoryResponse `json:"categories"`
}

type ProductImageResponse struct {
	ID    string `json:"id"`
	Image string `json:"image"`
}

type ProductCategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProductPagination struct {
	Products []ProductResponse `json:"products"`
	Limit    int               `json:"limit"`
	Page     int               `json:"page"`
	Total    int               `json:"total"`
}

func (product *Product) Response() ProductResponse {
	return ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Price:       product.Price,
		Description: product.Description,
	}
}
