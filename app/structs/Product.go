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
	Name        string                   `json:"name" validate:"required"`
	Price       float32                  `json:"price" validate:"required"`
	Description helpers.NullString       `json:"description" validate:"required"`
	Images      []ProductImageRequest    `json:"images" validate:"dive"`
	Variants    []ProductVariantRequest  `json:"variants" validate:"dive"`
	Categories  []CategoryProductRequest `json:"categories" validate:"min=1,dive"`
}

type ProductUpdateRequest struct {
	Name        string                   `json:"name" validate:"required"`
	Price       float32                  `json:"price" validate:"required"`
	Description helpers.NullString       `json:"description"`
	Images      []ProductImageRequest    `json:"images" validate:"dive"`
	Variants    []ProductVariantRequest  `json:"variants" validate:"dive"`
	Categories  []CategoryProductRequest `json:"categories" validate:"dive"`
}

type ProductResponse struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Price       float32                  `json:"price"`
	Description helpers.NullString       `json:"description"`
	Images      []ProductImageResponse   `json:"images"`
	Variants    []ProductVariantResponse `json:"variants"`
	Categories  []CategoryResponse       `json:"categories"`
}

type ProductImageResponse struct {
	ID    string `json:"id"`
	Image string `json:"image"`
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
		Images:      []ProductImageResponse{},
		Variants:    []ProductVariantResponse{},
		Categories:  []CategoryResponse{},
	}
}
