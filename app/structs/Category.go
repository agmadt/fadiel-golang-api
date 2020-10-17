package structs

import "time"

type Category struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	ProductID string    `db:"product_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CategoryPagination struct {
	Categories []CategoryResponse `json:"categories"`
	Limit      int                `json:"limit"`
	Page       int                `json:"page"`
	Total      int                `json:"total"`
}

type CategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

type CategoryProductRequest struct {
	ID string `json:"id" validate:"required"`
}

func (category *Category) Response() CategoryResponse {
	return CategoryResponse{
		ID:   category.ID,
		Name: category.Name,
	}
}
